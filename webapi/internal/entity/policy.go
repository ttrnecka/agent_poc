package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/db"
)

type Policy struct {
	cdb.BaseModel `bson:",inline"`
	Name          string   `bson:"name"`
	Description   string   `bson:"description"`
	FileName      string   `bson:"file_name"`
	Versions      []string `bson:"versions"`
}

func Policies() *cdb.CRUD[Policy] {
	return cdb.NewCRUD[Policy](db.Database(), "policies")
}
