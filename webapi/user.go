package main

type User struct {
	ID       int
	Email    string
	Username string
	Password string // hashed
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u User) GetResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
}
