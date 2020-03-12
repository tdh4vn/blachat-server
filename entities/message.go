package entities

import (
	"github.com/gocql/gocql"
	"time"
)

type Message struct {
	ID gocql.UUID `json:"id"`
	ChannelID gocql.UUID `json:"channel_id"`
	AuthorID gocql.UUID `json:"author_id"`
	Content string	`json:"content"`
	Type int `json:"type"`
	IsSystemMessage bool `json:"is_system_message"`
	CreatedAt time.Time `json:"created_at"`
}