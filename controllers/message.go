package controllers

import (
	"blachat-server/entities"
	"blachat-server/repo/repo_impl"
	"blachat-server/services"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
	"time"
)

type Message struct {
	UserRepo         *repo_impl.UserRepoImpl
	ParticipantsRepo *repo_impl.ParticipantsRepoImpl
	ChannelRepo      *repo_impl.ChannelRepoImpl
	MessageRepo 	 *repo_impl.MessageRepoImpl
}


func (controller *Message) GetMessages(_channelID gocql.UUID, lastItem *gocql.UUID, pageSize int) (int, gin.H) {
	if messages, err := controller.MessageRepo.GetMessages(_channelID, lastItem, pageSize); err != nil {
		return http.StatusInternalServerError, gin.H{"message": err.Error()}
	} else {
		return http.StatusOK, gin.H{"message": "success", "data": messages}
	}
}

func (controller *Message) GetNewMessages(_channelID gocql.UUID, latestItem *gocql.UUID) (int, gin.H) {
	if messages, err := controller.MessageRepo.GetNewMessage(_channelID, latestItem); err != nil {
		return http.StatusInternalServerError, gin.H{"message": err.Error()}
	} else {
		return http.StatusOK, gin.H{"message": "success", "data": messages}
	}
}

func (controller *Message) SendMarkReceive(channelID gocql.UUID, messageID gocql.UUID, receiveID gocql.UUID, actorActionID gocql.UUID) (int, gin.H) {
	go services.SendReceiveMessageEvent(messageID.String(), channelID.String(), receiveID.String(), actorActionID.String())
	return http.StatusOK, gin.H{"message": "success"}
}

func (controller *Message) SendMarkSeen(channelID gocql.UUID, messageID gocql.UUID, receiveID gocql.UUID, actorActionID gocql.UUID) (int, gin.H){
	go services.SendSeenMessageEvent(messageID.String(), channelID.String(), receiveID.String(), actorActionID.String())
	return http.StatusOK, gin.H{"message": "success"}
}

func (controller *Message) GetMessageById(messageID string) (int, gin.H){
	if id, err := gocql.ParseUUID(messageID); err != nil {
		return http.StatusBadRequest, gin.H{"message": err.Error()}
	} else {
		if message, err := controller.MessageRepo.FindById(id); err != nil {
			return http.StatusInternalServerError, gin.H{"message": err.Error()}
		} else {
			return http.StatusOK, gin.H{"message": "success", "data": message}
		}
	}
}

func (controller *Message) SendMessageToChannel(strChannelID string, strAuthorID string, content string, typeMessage int, isSystem bool) (int, gin.H){
	var channelID gocql.UUID
	var authorID gocql.UUID
	var messageID gocql.UUID
	var err error

	if channelID, err = gocql.ParseUUID(strChannelID); err != nil {
		return http.StatusInternalServerError, gin.H{"message": err.Error()}
	}

	if authorID, err = gocql.ParseUUID(strAuthorID); err != nil {
		return http.StatusInternalServerError, gin.H{"message": err.Error()}
	}

	messageID = gocql.TimeUUID()

	message := entities.Message{
		ID: messageID,
		AuthorID: authorID,
		ChannelID: channelID,
		Content: content,
		Type: typeMessage,
		IsSystemMessage: isSystem,
		CreatedAt: time.Now(),
	}

	if memberIds, err := controller.ParticipantsRepo.UserIdsInChannel(channelID); err != nil {
		println(err.Error())
	} else {
		for _, id := range memberIds {
			go services.SendMessageViaCentrigufo(&message, id.String())
		}
	}

	if _, err := controller.MessageRepo.SaveMessage(&message); err != nil {
		println(err.Error())
	}

	if _, err := controller.ChannelRepo.UpdateLastMessage(channelID, messageID); err != nil {
		println(err.Error())
	}


	return http.StatusOK, gin.H{
		"message": "success",
		"data": message,
	}

}

