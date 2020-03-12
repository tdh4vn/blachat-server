package entities

import (
	"github.com/gocql/gocql"
	"time"
)

type User struct {
	ID gocql.UUID `json:"id"`
	Name string	`json:"name"`
	Avatar string `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
