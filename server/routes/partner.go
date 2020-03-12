package routes

import (
	"blachat-server/controllers"
	"blachat-server/db"
	"blachat-server/middlewares"
	"blachat-server/repo/repo_impl"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreatePartnerRoute(v1Partner *gin.RouterGroup){
	v1Partner.Use(middlewares.PartnerAuthMiddleware())
	{
		userGroup := v1Partner.Group("users")
		{
			user := controllers.UserController{
				UserRepo: repo_impl.NewUserRepo(db.GetSession()),
			}

			userGroup.POST("/create", user.CreateUser)
			userGroup.PUT("/update/:id", user.UpdateUser)
			userGroup.POST("/create-token", user.CreateTokenForUser)
		}

		channelGroup := v1Partner.Group("channels")
		{
			channelController := controllers.Channels{
				UserRepo: repo_impl.NewUserRepo(db.GetSession()),
				ParticipantsRepo: repo_impl.NewParticipantsRepo(db.GetSession()),
				ChannelRepo: repo_impl.NewChanelRepo(db.GetSession()),
			}

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
					c.Abort()
					return
				}

				status, body := channelController.CreateChannelWithUsers(param.Name, param.Avatar, param.Type, param.UserIds)

				c.JSON(status, body)
				c.Abort()
			})
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

			channelGroup.GET("/members/:id", func(c *gin.Context) {
				id := c.Param("id")
				statusCode, body := channelController.GetMembersOfChannel(id)
				c.JSON(statusCode, body)
				c.Abort()
			})
		}
	}
}
