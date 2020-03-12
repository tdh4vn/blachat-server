package entities

import (
	"github.com/gocql/gocql"
	"time"
)

type Participants struct {
	UserID gocql.UUID
	ChannelID gocql.UUID
	LastReceive time.Time
	LastSeen time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}