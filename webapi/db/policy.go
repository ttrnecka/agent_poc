package db

import cdb "github.com/ttrnecka/agent_poc/common/db"

type Policy struct {
	cdb.BaseModel `bson:",inline"`
	Name          string   `bson:"name" json:"name"`
	Description   string   `bson:"description" json:"description"`
	FileName      string   `bson:"file_name" json:"file_name"`
	Versions      []string `bson:"versions" json:"versions"`
}

func Policies() *cdb.CRUD[Policy] {
	return cdb.NewCRUD[Policy](dB.database, "policies")
}
