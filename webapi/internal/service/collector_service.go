package service

import (
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
)

type CollectorService interface {
	GenericService[entity.Collector]
}

type collectorService struct {
	GenericService[entity.Collector]
}

func NewCollectorService(r repository.CollectorRepository) CollectorService {
	return &collectorService{
		NewGenericService(r),
	}
}
