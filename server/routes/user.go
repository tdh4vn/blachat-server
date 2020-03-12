package routes

import (
	"blachat-server/controllers"
	"blachat-server/db"
	"blachat-server/middlewares"
	"blachat-server/repo/repo_impl"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
	"strconv"
)

func CreateUserRoute(v1User *gin.RouterGroup) {
	v1User.Use(middlewares.UserAuthMiddleware())
	{

		membersGroup := v1User.Group("members")
		{
			userController := controllers.UserController{
				UserRepo: repo_impl.NewUserRepo(db.GetSession()),
			}

			membersGroup.POST("/gets", func (c *gin.Context) {
				type Params struct {
					Ids []string `json:"ids"`
				}

				var params Params

				if err := c.BindJSON(&params); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
					c.Abort()
					return
				}

				ok, body := userController.GetUserByIds(params.Ids)

				c.JSON(ok, body)
				c.Abort()
				return

			})

		}

		channelGroup := v1User.Group("channels")
		{
			channelController := controllers.Channels{
				UserRepo: repo_impl.NewUserRepo(db.GetSession()),
				ParticipantsRepo: repo_impl.NewParticipantsRepo(db.GetSession()),
				ChannelRepo: repo_impl.NewChanelRepo(db.GetSession()),
			}

			//Create channel with users invite
			channelGroup.POST("/create", func (c *gin.Context) {
				type CreateChannelParam struct{
					Name string `json:"name"`
					UserIds []string `json:"userIds" binding:"required"`
					Avatar string `json:"avatar"`
					Type int `json:"type" binding:"required"`
				}

				var param CreateChannelParam

				if err := c.BindJSON(&param); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				}

				param.UserIds = append(param.UserIds, c.GetString("userID"))

				status, body := channelController.CreateChannelWithUsers(param.Name, param.Avatar, param.Type, param.UserIds)

				c.JSON(status, body)
				c.Abort()
			})

			channelGroup.PUT("/events/typing/:id", func(c *gin.Context) {
				id := c.Param("id")
				statusCode, body := channelController.UserTyping(id, c.GetString("userID"))
				c.JSON(statusCode, body)
				c.Abort()
			})

			channelGroup.PUT("/events/stop-typing/:id", func(c *gin.Context) {
				id := c.Param("id")
				statusCode, body := channelController.UserStopTyping(
					id, c.GetString("userID"))
				c.JSON(statusCode, body)
				c.Abort()
			})

			//Get members id of channel
			channelGroup.GET("/members/:id", func(c *gin.Context) {
				id := c.Param("id")
				statusCode, body := channelController.GetMembersOfChannel(id)
				c.JSON(statusCode, body)
				c.Abort()
			})

			//Invite users to channel
			channelGroup.POST("/invite/:id", func(c *gin.Context) {
				type Param struct {
					UserIds []string `json:"userIds" binding:"required"`
				}
				var param Param
				id := c.Param("id")
				if err := c.BindJSON(&param); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
					c.Abort()
					return
				}

				status, body := channelController.InviteUsersToChannel(id, param.UserIds)

				c.JSON(status, body)
				c.Abort()
			})

			//get channel user joined
			channelGroup.GET("/me", func (c *gin.Context) {
				if pageSize, err := strconv.Atoi(c.Query("pageSize")); err != nil {
					pageSize = 20
				} else {
					if lastId, err := gocql.ParseUUID(c.Query("lastId")); err != nil {
						status, body := channelController.GetChannelsUserJoined(c.GetString("userID"), nil, pageSize)
						c.JSON(status, body)
					} else {
						status, body := channelController.GetChannelsUserJoined(c.GetString("userID"), &lastId, pageSize)
						c.JSON(status, body)
					}
				}
			})

		}
	}
}