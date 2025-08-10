package db

import cdb "github.com/ttrnecka/agent_poc/common/db"

type Collector struct {
	cdb.BaseModel `bson:",inline"`
	Name          string `bson:"name" json:"name"`
	Password      string `bson:"password" json:"-"` // Excluded from JSON output
	Status        string `bson:"status" json:"status"`
}

func Collectors() *cdb.CRUD[Collector] {
	return cdb.NewCRUD[Collector](dB.database, "collectors")
}
