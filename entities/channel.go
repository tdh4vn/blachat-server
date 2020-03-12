package entities

import (
	"github.com/gocql/gocql"
	"time"
)

type Channel struct {
	ID gocql.UUID `json:"id"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
	Type int `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool `json:"is_deleted"`
	LastMessageID *gocql.UUID `json:"last_message_id"`

}
