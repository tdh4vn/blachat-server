package repo_impl

import (
	"blachat-server/entities"
	"github.com/gocql/gocql"
	"time"
)

const createChannelQuery = `INSERT INTO channel(id, name, avatar, type, deleted, created_at, updated_at) VALUES (?,?,?,?,?,?,?)`

const findByIdChannelQuery = `SELECT id, name, avatar, type, deleted, last_message_id, created_at, updated_at FROM channel WHERE id = ?`

const findChannelsIdUserJoin = `SELECT channel_id FROM participants WHERE user_id = ?`

const findChannelsPagingUserJoin = `SELECT id, name, avatar, type, deleted, last_message_id, created_at, updated_at FROM channel WHERE id IN` +
	` ? AND token(id) < token(?) LIMIT ?`

const findFirstChannelsUserJoin = `SELECT id, name, avatar, type, deleted, last_message_id, created_at, updated_at FROM channel WHERE id IN` +
	` ? LIMIT ?`

const deleteChannelQuery = `UPDATE channel SET deleted=? WHERE id = ?`

const updateLastMessageOfChannel = `UPDATE channel SET last_message_id = ?, updated_at = ? WHERE id = ?`

type ChannelRepoImpl struct {
	DbSession *gocql.Session
}

func NewChanelRepo(db *gocql.Session) *ChannelRepoImpl {
	return &ChannelRepoImpl{
		DbSession: db,
	}
}

func (c *ChannelRepoImpl) UpdateLastMessage(id gocql.UUID, mID gocql.UUID) (*entities.Channel, error) {
	channel := entities.Channel{}

	if err := c.DbSession.Query(updateLastMessageOfChannel, mID, time.Now(), id).Exec(); err != nil {
		return nil, err
	} else {
		return &channel, nil
	}
}

func (c *ChannelRepoImpl) Create(name string, avatar string, channelType int) (*entities.Channel, error) {
	if cID, err := gocql.RandomUUID(); err != nil {
		return nil, err
	} else {
		createdAt := time.Now()
		updatedAt := time.Now()

		if err := c.DbSession.Query(createChannelQuery, cID, name, avatar, channelType, false, time.Now(), time.Now()).Exec(); err != nil {
			return nil, err
		} else {
			return &entities.Channel{
				ID: cID,
				Name: name,
				Avatar: avatar,
				Type: channelType,
				UpdatedAt:updatedAt,
				CreatedAt:createdAt,
				IsDeleted: false,
				LastMessageID: nil,
			}, nil
		}
	}
}

func (c *ChannelRepoImpl) FindById(id gocql.UUID) (*entities.Channel, error) {

	channel := entities.Channel{}

	if err := c.DbSession.Query(findByIdChannelQuery, id).Consistency(gocql.One).Scan(
		&channel.ID,
		&channel.Name,
		&channel.Avatar,
		&channel.Type,
		&channel.IsDeleted,
		&channel.CreatedAt,
		&channel.UpdatedAt); err != nil {
		return nil, err
	} else {
		return &channel, nil
	}
}

func (c *ChannelRepoImpl) ChannelsUserJoin(userId gocql.UUID, lastId *gocql.UUID, pageSize int) ([]*entities.Channel, error) {
	var channels []*entities.Channel

	var iter *gocql.Iter

	var cId gocql.UUID

	var channelsIDUserJoined []gocql.UUID

	var iterForChannelsIDs *gocql.Iter

	iterForChannelsIDs = c.DbSession.Query(findChannelsIdUserJoin, userId).Iter()

	for iterForChannelsIDs.Scan(&cId) {
		channelsIDUserJoined = append(channelsIDUserJoined, cId)
	}

	if lastId != nil {
		iter = c.DbSession.Query(findChannelsPagingUserJoin, channelsIDUserJoined, lastId, pageSize).Iter()
	} else {

		iter = c.DbSession.Query(findFirstChannelsUserJoin, channelsIDUserJoined, pageSize).Iter()
	}

	var id gocql.UUID
	var name string
	var avatar string
	var cType int
	var deleted bool
	var updatedAt time.Time
	var createdAt time.Time
	var lastMessageId gocql.UUID

	for iter.Scan(&id, &name, &avatar, &cType, &deleted, &lastMessageId, &createdAt, &updatedAt) {
		var _lastMessageId *gocql.UUID = nil
		if lastMessageId.Time().Unix() > 0{
			__lastMessageId, _ := gocql.ParseUUID(lastMessageId.String())
			_lastMessageId = &__lastMessageId
		}

		channels = append(channels, &entities.Channel{
			ID: id,
			Name: name,
			Avatar: avatar,
			Type: cType,
			IsDeleted: deleted,
			LastMessageID: _lastMessageId,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return channels, nil

}

func (c *ChannelRepoImpl) Delete(id gocql.UUID) (bool, error) {
	if err := c.DbSession.Query(deleteChannelQuery, id).Exec(); err != nil {
		return false, err
	}

	return true, nil
}






