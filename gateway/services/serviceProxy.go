package gateway

import (
	domain "egosystem.org/micros/gateway/domain"
)

type ProxyService interface {
	SaveSecretKEY(string, string) error
}

type DefaultProxyService struct {
	repo domain.ProxyRepositoryStorage
}

func NewProxyService(repository domain.ProxyRepositoryStorage) DefaultProxyService {
	return DefaultProxyService{repo: repository}
}

func (s DefaultProxyService) SaveSecretKEY(engine, apikey string) (string, error) {
	err := s.repo.SaveSecretKEY(engine, apikey)
	if err != nil {
		return "failed", err
	}
	return "ok", nil
}
