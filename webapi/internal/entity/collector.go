package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/db"
)

type Collector struct {
	cdb.BaseModel `bson:",inline"`
	Name          string `bson:"name"`
	Status        string `bson:"status"`
	Password      string `bson:"password"`
}

func Collectors() *cdb.CRUD[Collector] {
	return cdb.NewCRUD[Collector](db.Database(), "collectors")
}
