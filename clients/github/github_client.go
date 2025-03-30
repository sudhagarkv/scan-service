package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"

	"scan-service/constants"
	"scan-service/models"
	"scan-service/utils"
)

type Github struct {
	client    *http.Client
	token     string
	keyClient *azkeys.Client
}

func (s Github) HasPrivateAccess(ctx context.Context, url, encryptedToken string) (bool, error) {
	owner, repoName, err := utils.SplitGitHubURL(url)
	if err != nil {
		log.Printf("Unable to parse github owner and repo name %v", err)
		return false, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repoName), nil)
	if err != nil {
		log.Printf("Unable to create request to get repo details %v", err)
		return false, err
	}

	version := ""
	algorithmRSAOAEP := azkeys.EncryptionAlgorithmRSAOAEP
	decrypt, err := s.keyClient.Decrypt(ctx, constants.TokenKey, version, azkeys.KeyOperationParameters{
		Algorithm: &algorithmRSAOAEP,
		Value:     []byte(encryptedToken),
	}, nil)
	if err != nil {
		log.Printf("Unable to decrypt token using token key %v", err)
		return false, err
	}

	request.Header.Add("Authorization", "Bearer "+string(decrypt.Result))

	response, err := s.client.Do(request)
	if err != nil {
		log.Printf("Unable to make http request to get repo details %v", err)
		return false, err
	}

	defer response.Body.Close()
	var repo models.GithubRepo
	err = json.NewDecoder(response.Body).Decode(&repo)
	if err != nil {
		log.Printf("Unable to decode repo response %v", err)
		return false, err
	}

	return true, nil
}

func (s Github) HasPublicAccess(ctx context.Context, url string) (bool, error) {
	owner, repoName, err := utils.SplitGitHubURL(url)
	if err != nil {
		log.Printf("Unable to parse github owner and repo name %v", err)
		return false, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repoName), nil)
	if err != nil {
		log.Printf("Unable to create request to get repo details %v", err)
		return false, err
	}

	request.Header.Add("Authorization", "Bearer "+s.token)

	response, err := s.client.Do(request)
	if err != nil {
		log.Printf("Unable to make http request to get repo details %v", err)
		return false, err
	}

	defer response.Body.Close()
	var repo models.GithubRepo
	err = json.NewDecoder(response.Body).Decode(&repo)
	if err != nil {
		log.Printf("Unable to decode repo response %v", err)
		return false, err
	}

	return true, nil
}

func NewGithub(client *http.Client, keysClient *azkeys.Client, apiToken string) *Github {
	return &Github{
		client:    client,
		token:     apiToken,
		keyClient: keysClient,
	}
}
