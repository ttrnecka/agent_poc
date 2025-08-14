package mapper

import (
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToProbeEntity(p dto.ProbeDTO) entity.Probe {
	probe := entity.Probe{
		Policy:   p.Policy,
		Version:  p.Version,
		Address:  p.Address,
		Port:     p.Port,
		User:     p.User,
		Password: p.Password,
	}
	probe.ID, _ = primitive.ObjectIDFromHex(p.ID)
	probe.CollectorID, _ = primitive.ObjectIDFromHex(p.CollectorID)
	return probe
}

func ToProbeDTO(p entity.Probe) dto.ProbeDTO {
	return dto.ProbeDTO{
		ID:          p.ID.String(),
		Policy:      p.Policy,
		Version:     p.Version,
		Address:     p.Address,
		Port:        p.Port,
		User:        p.User,
		Password:    p.Password,
		CollectorID: p.CollectorID.String(),
	}
}
