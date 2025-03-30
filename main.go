package main

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"net/url"
	"os"
	"scan-service/middleware"
	"scan-service/models"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"scan-service/clients/github"
	"scan-service/constants"
	"scan-service/controller"
	"scan-service/repository"
	"scan-service/scm"
	"scan-service/service"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("unable to load .env")
	}
}

func init() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = validate.RegisterValidation("isPrivate", func(fl validator.FieldLevel) bool {
			request := fl.Parent().Interface().(models.ScanRequest)
			if request.IsPrivate {
				return len(request.EncryptedToken) > 0
			}
			return true
		}, false)
	}
}

func main() {
	router := gin.New()

	client := &http.Client{}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	keyVaultURL := os.Getenv("KEY_VAULT_URL")

	secretClient, err := azsecrets.NewClient(keyVaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	azkeysClient, err := azkeys.NewClient(keyVaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	dbPassword, err := secretClient.GetSecret(context.TODO(), constants.DBPasswordKey, "", nil)
	if err != nil {
		panic(err)
	}

	// PostgreSQL connection details
	dbInfo := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("DB_USER"), *dbPassword.Value),
		Host:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: "sslmode=disable",
	}

	// Connect to the database
	db, err := sqlx.Open("postgres", dbInfo.String())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	githubAPIToken, err := secretClient.GetSecret(context.Background(), constants.GithubAPITokenKey, "", nil)
	if err != nil {
		panic(err)
	}

	// Kafka broker URLs, replace with your IBM Event Streams broker URLs
	brokerList := []string{os.Getenv("KAFKA_HOST")}

	// Create Kafka producer configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing Kafka producer: %v", err)
		}
	}()

	newGithub := github.NewGithub(client, azkeysClient, *githubAPIToken.Value)
	scmFactoryService := scm.NewFactoryService(newGithub)
	scanRepository := repository.NewScanRepository(db)
	scanService := service.NewScanService(scmFactoryService, producer, scanRepository, os.Getenv("TOPIC_NAME"))
	scanController := controller.NewScanController(scanService)

	router.Use(middleware.CORS())
	router.Handle(http.MethodPost, "/api/v1/scan", scanController.ProcessRequest)

	port := os.Getenv("PORT")

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Unable to start the server %v", err)
	}
}
