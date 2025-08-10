package db

import cdb "github.com/ttrnecka/agent_poc/common/db"

type User struct {
	cdb.BaseModel `bson:",inline"`
	Username      string `bson:"username" json:"username"`
	Email         string `bson:"email" json:"email"`
	Password      string `bson:"password" json:"-"` // Excluded from JSON output
}
