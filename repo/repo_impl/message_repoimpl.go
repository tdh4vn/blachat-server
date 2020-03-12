package repo_impl

import (
	"blachat-server/entities"
	"github.com/gocql/gocql"
	"time"
)


const insertMessageQuery = `INSERT INTO messages(id, channel_id, author_id, content, created_at, system, type) VALUES(?,?,?,?,?,?,?)`

const findPagingMessageOfChannel = `SELECT id, channel_id, author_id, content, created_at, system, type FROM messages WHERE channel_id = ?` +
	` AND id < ? LIMIT ?`

const findByIDQuery = `SELECT id, channel_id, author_id, content, created_at, system, type FROM messages WHERE id = ?`

const findFirstPageMessageOfChannel = `SELECT id, channel_id, author_id, content, created_at, system, type FROM messages WHERE channel_id = ?` +
	` LIMIT ?`

const findLatestMessageOfChannel = `SELECT id, channel_id, author_id, content, created_at, system, type FROM messages WHERE channel_id = ?` +
	` AND id > ?`


type MessageRepoImpl struct {
	DbSession *gocql.Session
}

func NewMessageRepo(db *gocql.Session) *MessageRepoImpl {
	return &MessageRepoImpl{
		DbSession: db,
	}
}

func (repo *MessageRepoImpl) GetNewMessage(channelID gocql.UUID, latestItem *gocql.UUID) ([]*entities.Message, error) {
	var iter *gocql.Iter

	iter = repo.DbSession.Query(findLatestMessageOfChannel, channelID, latestItem).Iter()

	var messages []*entities.Message

	var id gocql.UUID
	var channelId gocql.UUID
	var authorID gocql.UUID
	var content string
	var createdAt time.Time
	var isSystem bool
	var messageType int

	for iter.Scan(&id, &channelId, &authorID, &content, &createdAt, &isSystem, &messageType) {

		messages = append(messages, &entities.Message{
			ID:id,
			ChannelID:channelId,
			AuthorID:authorID,
			Content:content,
			CreatedAt:createdAt,
			IsSystemMessage: isSystem,
			Type:messageType,
		})
	}

	return messages, nil
}

func (repo *MessageRepoImpl) SaveMessage(message *entities.Message) (*entities.Message, error) {
	if err := repo.DbSession.Query(insertMessageQuery,
		message.ID,
		message.ChannelID,
		message.AuthorID,
		message.Content,
		message.CreatedAt,
		message.IsSystemMessage,
		message.Type).Exec(); err != nil {
			return nil, err
	} else {
		return message, nil
	}
}


func (repo *MessageRepoImpl) FindById(messageID gocql.UUID) (*entities.Message, error) {

	message := entities.Message{}

	if err := repo.DbSession.Query(findByIDQuery, messageID).Consistency(gocql.One).Scan(
		&message.ID,
		&message.ChannelID,
		&message.AuthorID,
		&message.Content,
		&message.CreatedAt,
		&message.IsSystemMessage,
		&message.Type); err != nil {
		return nil, err
	} else {
		return &message, nil
	}
}

func (repo *MessageRepoImpl) GetMessages(_channelID gocql.UUID, lastItem *gocql.UUID, pageSize int) ([]*entities.Message, error) {

	var iter *gocql.Iter;

	if lastItem == nil {
		iter = repo.DbSession.Query(findFirstPageMessageOfChannel, _channelID, pageSize).Iter()
	} else {
		iter = repo.DbSession.Query(findPagingMessageOfChannel, _channelID, lastItem, pageSize).Iter()
	}

	var messages []*entities.Message

	var id gocql.UUID
	var channelId gocql.UUID
	var authorID gocql.UUID
	var content string
	var createdAt time.Time
	var isSystem bool
	var messageType int

	for iter.Scan(&id, &channelId, &authorID, &content, &createdAt, &isSystem, &messageType) {
		messages = append(messages, &entities.Message{
			ID: id,
			ChannelID: channelId,
			AuthorID: authorID,
			Content: content,
			CreatedAt: createdAt,
			IsSystemMessage: isSystem,
			Type: messageType,
		})
	}

	return messages, nil

}

