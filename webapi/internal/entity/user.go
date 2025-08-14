package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/db"
)

type User struct {
	cdb.BaseModel `bson:",inline"`
	Username      string `bson:"username"`
	Email         string `bson:"email"`
	Password      string `bson:"password"`
}

func Users() *cdb.CRUD[User] {
	return cdb.NewCRUD[User](db.Database(), "users")
}
