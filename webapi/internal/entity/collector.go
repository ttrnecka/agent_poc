package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/db"
)

// omitempty - the field will not be updated if not provided - keeping the value in DB
type Collector struct {
	cdb.BaseModel `bson:",inline"`
	Name          string `bson:"name"`
	Status        string `bson:"status,omitempty"`
	Password      string `bson:"password,omitempty"`
}

func Collectors() *cdb.CRUD[Collector] {
	return cdb.NewCRUD[Collector](db.Database(), "collectors")
}
