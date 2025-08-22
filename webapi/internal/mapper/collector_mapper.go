package mapper

import (
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToCollectorEntity(c dto.CollectorDTO) entity.Collector {
	coll := entity.Collector{
		Name:     c.Name,
		Status:   c.Status,
		Password: c.Password,
	}
	coll.ID, _ = primitive.ObjectIDFromHex(c.ID)
	return coll
}

func ToCollectorDTO(c entity.Collector) dto.CollectorDTO {
	return dto.CollectorDTO{
		ID:     c.ID.Hex(),
		Name:   c.Name,
		Status: c.Status,
	}
}
