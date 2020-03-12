package dto

import (
	"blachat-server/entities"
	"github.com/gocql/gocql"
	"time"
)

type ChannelDTO struct {
	ID gocql.UUID `json:"id"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
	Type int `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool `json:"is_deleted"`
	LastMessageID *gocql.UUID `json:"last_message_id"`
	MemberIds []gocql.UUID `json:"member_ids"`
}

func MapToChannelDTO(channel entities.Channel, memberIds []gocql.UUID) ChannelDTO {
	return ChannelDTO{
		ID: channel.ID,
		Name: channel.Name,
		Avatar: channel.Avatar,
		Type: channel.Type,
		CreatedAt: channel.CreatedAt,
		UpdatedAt: channel.UpdatedAt,
		IsDeleted: channel.IsDeleted,
		LastMessageID: channel.LastMessageID,
		MemberIds: memberIds,
	}
}