package scm

import (
	"context"
	"errors"
)

type Type string

const (
	Github Type = "Github"
	GitLab Type = "GitLab"
)

type FactoryService interface {
	HasPrivateAccess(ctx context.Context, url, encryptedToken string) (bool, error)
	HasPublicAccess(ctx context.Context, url string) (bool, error)
}

type Clients map[Type]FactoryService

func (c Clients) GetSCMService(scmType Type) (FactoryService, error) {
	if service, ok := c[scmType]; ok {
		return service, nil
	}
	return nil, errors.New("scm type not supported")
}

func NewFactoryService(githubClient FactoryService) Clients {
	return Clients{
		Github: githubClient,
	}
}
