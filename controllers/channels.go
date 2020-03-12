package controllers

import (
	"blachat-server/dto"
	"blachat-server/repo"
	"blachat-server/repo/repo_impl"
	"blachat-server/services"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
)

type Channels struct {
	UserRepo         *repo_impl.UserRepoImpl
	ParticipantsRepo *repo_impl.ParticipantsRepoImpl
	ChannelRepo      *repo_impl.ChannelRepoImpl
	ContactsRepo	 *repo.ContactsRepoImpl
}

func (controller *Channels) CreateChannelWithUsers(name string, avatar string, channelType int, users []string) (int, gin.H) {
	if channelType == 2 && len(users) != 2 {
		return http.StatusInternalServerError, gin.H{"message": "cannot create direct channel with users list"}
	}

	if channelType > 2 || channelType < 1 {
		return http.StatusInternalServerError, gin.H{"message": "not support channel type"}
	}

	if channel, err := controller.ChannelRepo.Create(name, avatar, channelType); err != nil {
		return http.StatusInternalServerError, gin.H{"message": err.Error()}
	} else {
		var userIds []gocql.UUID

		for _, id := range users {
			if uuid, err := gocql.ParseUUID(id); err == nil {
				userIds = append(userIds, uuid)
			}
		}

		go func(){
			for i := 0; i < len(userIds); i++ {
				var secondaryUserIds []gocql.UUID
				for j := 0; j < len(userIds); j++ {
					if i != j {
						secondaryUserIds = append(secondaryUserIds, userIds[j])
					}
				}
				_ = controller.ContactsRepo.CreateContactsForUser(userIds[i], secondaryUserIds);
			}
		}()

		if err = controller.ParticipantsRepo.InviteToChannel(channel.ID, userIds); err != nil {
			return http.StatusBadRequest, gin.H{"message": err.Error()}
		} else {
			for _, id := range userIds {
				go services.SendNewChannel(channel, id.String())
			}
			return http.StatusOK, gin.H{
				"message": "success",
				"data": dto.MapToChannelDTO(*channel, userIds),
			}
		}

	}
}

func (controller *Channels) InviteUsersToChannel(channelId string, ids []string) (int, gin.H) {
	var userIds []gocql.UUID
	for _, id := range ids {
		if uuid, err := gocql.ParseUUID(id); err == nil {
			userIds = append(userIds, uuid)
		}
	}

	if cId, err := gocql.ParseUUID(channelId); err != nil {
		return http.StatusBadRequest, gin.H{"message": "id not uuid format"}
	} else {

		if currentUserInChannel, err := controller.ParticipantsRepo.UserIdsInChannel(cId); err == nil {
			go func() {
				for _, idOfNewUser := range userIds {
					_ = controller.ContactsRepo.CreateContactsForUser(idOfNewUser, currentUserInChannel)
				}
				for _, idOfCurrentUser := range currentUserInChannel {
					_ = controller.ContactsRepo.AddUserToContactOfUsers(currentUserInChannel, idOfCurrentUser)
				}
			}()
		}

		if err := controller.ParticipantsRepo.InviteToChannel(cId, userIds); err != nil {
			return http.StatusBadRequest, gin.H{"message": err.Error()}
		} else {
			return http.StatusOK, gin.H{
				"message": "success",
				"data": userIds,
			}
		}
	}
}

func (controller *Channels) UserTyping(cID string, uID string) (int, gin.H) {
	if channelID, err := gocql.ParseUUID(cID); err != nil {
		return http.StatusBadRequest, gin.H{"message": err.Error()}
	} else {
		if userIDs, err := controller.ParticipantsRepo.UserIdsInChannel(channelID); err != nil {
			return http.StatusBadRequest, gin.H{"message": err.Error()}
		} else {
			for _,id := range userIDs {
				if id.String() != uID {
					go services.SendTypingEvent(id.String(), cID, uID, true)
				}
			}
		}
		return http.StatusOK, gin.H{"message": "success"}
	}
}

func (controller *Channels) UserStopTyping(cID string, uID string) (int, gin.H) {
	if channelID, err := gocql.ParseUUID(cID); err != nil {
		return http.StatusBadRequest, gin.H{"message": err.Error()}
	} else {

		if userIDs, err := controller.ParticipantsRepo.UserIdsInChannel(channelID); err != nil {
			return http.StatusBadRequest, gin.H{"message": err.Error()}
		} else {
			for _,id := range userIDs {
				if id.String() != uID {
					go services.SendTypingEvent(id.String(), cID, uID, false)
				}
			}
		}
		return http.StatusOK, gin.H{"message": "success"}
	}
}

func (controller *Channels) GetMembersOfChannel(id string) (int, gin.H) {
	if cId, err := gocql.ParseUUID(id); err != nil {
		return http.StatusBadRequest, gin.H{"message": err.Error()}
	} else {
		if memberIds, err := controller.ParticipantsRepo.UserIdsInChannel(cId); err != nil {
			return http.StatusInternalServerError, gin.H{"message": err.Error()}
		} else {
			return http.StatusOK, gin.H{"message": "success", "data": memberIds}
		}
	}
}

func (controller *Channels) GetChannelsUserJoined(uID string, lastId *gocql.UUID, pageSize int) (int, gin.H) {
	if uID, err := gocql.ParseUUID(uID); err != nil {
		return http.StatusBadRequest, gin.H{"message": err.Error()}
	} else {
		if channels, err := controller.ChannelRepo.ChannelsUserJoin(uID, lastId, pageSize); err != nil {
			return http.StatusInternalServerError, gin.H{"message": err.Error()}
		} else {
			return http.StatusOK, gin.H{"message": "success", "data": channels}
		}
	}
}
