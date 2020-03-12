package routes

import (
	"blachat-server/controllers"
	"blachat-server/db"
	"blachat-server/middlewares"
	"blachat-server/repo/repo_impl"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
)

func CreateMessageRoute(v1Message *gin.RouterGroup) {
	v1Message.Use(middlewares.UserAuthMiddleware())
	{

		message := controllers.Message{
			MessageRepo: repo_impl.NewMessageRepo(db.GetSession()),
			UserRepo: repo_impl.NewUserRepo(db.GetSession()),
			ParticipantsRepo: repo_impl.NewParticipantsRepo(db.GetSession()),
			ChannelRepo: repo_impl.NewChanelRepo(db.GetSession()),
		}

		v1Message.POST("/create", func (c *gin.Context) {
			type Params struct {
				Message string `json:"message"`
				Type int `json:"type"`
				ChannelID string `json:"channel_id"`
			}

			var param Params

			if err := c.BindJSON(&param); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			status, body := message.SendMessageToChannel(param.ChannelID, c.GetString("userID"), param.Message, param.Type, false)

			c.JSON(status, body)
			c.Abort()
			return

		})

		v1Message.POST("/mark-receive", func(c *gin.Context) {
			type Params struct {
				MessageID string `json:"message_id"`
				ChannelID string `json:"channel_id"`
				ReceiveID string `json:"receive_id"`
			}

			var param Params

			if err := c.BindJSON(&param); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			var messageID gocql.UUID
			var channelID gocql.UUID
			var userID gocql.UUID
			var receiveID gocql.UUID
			var err error

			if messageID, err = gocql.ParseUUID(param.MessageID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}


			if channelID, err = gocql.ParseUUID(param.ChannelID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			if userID, err = gocql.ParseUUID(c.GetString("userID")); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			if receiveID, err = gocql.ParseUUID(param.ReceiveID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			status, body := message.SendMarkReceive(channelID, messageID, receiveID, userID)

			c.JSON(status, body)
			c.Abort()
			return

		})

		v1Message.POST("/mark-seen", func(c *gin.Context) {
			type Params struct {
				MessageID string `json:"message_id"`
				ChannelID string `json:"channel_id"`
				ReceiveID string `json:"receive_id"`
			}

			var param Params

			if err := c.BindJSON(&param); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			var messageID gocql.UUID
			var channelID gocql.UUID
			var userID gocql.UUID
			var receiveID gocql.UUID
			var err error

			if messageID, err = gocql.ParseUUID(param.MessageID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}


			if channelID, err = gocql.ParseUUID(param.ChannelID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			if userID, err = gocql.ParseUUID(c.GetString("userID")); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			if receiveID, err = gocql.ParseUUID(param.ReceiveID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			}

			status, body := message.SendMarkSeen(channelID, messageID, receiveID, userID)

			c.JSON(status, body)
			c.Abort()
			return

		})

		v1Message.GET("/channel/:id", func (c *gin.Context) {
			if cID, err := gocql.ParseUUID(c.Param("id")); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				c.Abort()
				return
			} else {
				var status int
				var body gin.H
				if lastID, err := gocql.ParseUUID(c.Query("lastId")); err == nil {
					status, body = message.GetMessages(cID, &lastID, 20)
				} else if latestID, err := gocql.ParseUUID(c.Query("latestId")); err == nil {
					status, body = message.GetNewMessages(cID, &latestID)
				} else {
					status, body = message.GetMessages(cID, nil, 20)
				}

				c.JSON(status, body)
				c.Abort()
				return
			}
		})

		v1Message.GET("/get-by-id/:id", func (c *gin.Context) {
			status, body := message.GetMessageById(c.Param("id"))
			c.JSON(status, body)
			c.Abort()
		})

	}
}