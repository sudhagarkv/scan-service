package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"

	"scan-service/models"
	"scan-service/repository"
	"scan-service/scm"
	"scan-service/utils"
)

type ScanService interface {
	ProcessScanRequest(ctx context.Context, request models.ScanRequest) error
}

type scanService struct {
	scmFactory     scm.Clients
	producerClient sarama.SyncProducer
	scanRepository repository.ScanRepository
	topic          string
}

func (s scanService) ProcessScanRequest(ctx context.Context, request models.ScanRequest) error {
	service, err := s.scmFactory.GetSCMService(request.Type)
	if err != nil {
		log.Printf("Unable to get scm service %v", err)
		return err
	}

	var access bool
	if request.IsPrivate {
		access, err = service.HasPrivateAccess(ctx, request.URL, request.EncryptedToken)
		if err != nil {
			log.Printf("Unable to access repository %v", err)
			return err
		}
	} else {
		access, err = service.HasPublicAccess(ctx, request.URL)
		if err != nil {
			log.Printf("Unable to access repository %v", err)
			return err
		}
	}

	if !access {
		return errors.New("cannot access repository using provided token")
	}

	repoOwner, repoName, err := utils.SplitGitHubURL(request.URL)
	if err != nil {
		log.Printf("Unable to prase github url %v", err)
		return err
	}

	scanID, err := s.scanRepository.InsertScan(ctx, request)
	if err != nil {
		log.Printf("Unable to insert scan request %v", err)
		return err
	}

	scan := models.QueuedScan{
		RepoName:       repoName,
		IsPrivate:      request.IsPrivate,
		EncryptedToken: request.EncryptedToken,
		Namespace:      repoOwner,
		URL:            request.URL,
		ID:             *scanID,
	}

	marshal, err := json.Marshal(scan)
	if err != nil {
		log.Printf("Unable to convert queue scan to bytes %v", err)
		return err
	}
	partition, offset, err := s.producerClient.SendMessage(&sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.StringEncoder(marshal),
		Partition: 0,
	})
	if err != nil {
		log.Printf("Unable to send message into kafka queue %v", err)
		return err
	}

	log.Printf("Item pushed into the queue at offset %d and partition %d", offset, partition)

	err = s.scanRepository.UpdateQueueStatus(ctx, *scanID, models.Queued)
	if err != nil {
		log.Printf("Unable to update queue status %v", err)
		return err
	}

	return nil
}

func NewScanService(service scm.Clients, producer sarama.SyncProducer, scanRepository repository.ScanRepository, topic string) ScanService {
	return &scanService{
		scmFactory:     service,
		producerClient: producer,
		scanRepository: scanRepository,
		topic:          topic,
	}
}
