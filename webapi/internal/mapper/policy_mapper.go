package mapper

import (
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToPolicyEntity(p dto.PolicyDTO) entity.Policy {
	pol := entity.Policy{
		Name:        p.Name,
		Description: p.Description,
		FileName:    p.FileName,
		Versions:    p.Versions,
	}
	pol.ID, _ = primitive.ObjectIDFromHex(p.ID)
	return pol
}

func ToPolicyDTO(p entity.Policy) dto.PolicyDTO {
	return dto.PolicyDTO{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		FileName:    p.FileName,
		Versions:    p.Versions,
	}
}
