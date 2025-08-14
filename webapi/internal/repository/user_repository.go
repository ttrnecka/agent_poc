package repository

import (
	"context"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
)

// import (
//     "database/sql"
//     "myapp/internal/user/entity"
// )

type UserRepository interface {
	// Create(user entity.User) (entity.User, error)
	GetByField(context.Context, string, interface{}) (*entity.User, error)
}

type userRepository struct {
	*cdb.CRUD[entity.User]
}

func NewUserRepository(db *cdb.CRUD[entity.User]) UserRepository {
	return &userRepository{db}
}

// func (r *userRepository) Create(user entity.User) (entity.User, error) {
//     query := `INSERT INTO users (name, email) VALUES (?, ?) RETURNING id`
//     err := r.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
//     return user, err
// }
