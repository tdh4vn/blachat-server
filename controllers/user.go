package controllers

import (
	"blachat-server/config"
	"blachat-server/entities"
	"blachat-server/repo/repo_impl"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
	"time"
)

type UserController struct {
	UserRepo *repo_impl.UserRepoImpl
}

func (controller *UserController) CreateUser(c *gin.Context) {
	name := c.PostForm("name")
	avatar := c.PostForm("avatar")

	if name == "" || avatar == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	}

	if user, err := controller.UserRepo.Insert(&entities.User{
		Name: name,
		Avatar: avatar,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server cannot create user"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User created", "data": user})
	}
	c.Abort()
	return
}

func (controller *UserController) UpdateUser(c *gin.Context){
	name := c.PostForm("name")
	avatar := c.PostForm("avatar")
	if id, err := gocql.ParseUUID(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else {
		if user, err := controller.UserRepo.Update(&entities.User{
			ID: id,
			Name: name,
			Avatar: avatar,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "user info updated", "data": user})
		}
	}
	c.Abort()
	return
}

func (controller *UserController) CreateTokenForUser(c *gin.Context) {
	id := c.PostForm("id")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,
		"sub": id,
		"client": id,
		"channel": "$chat:" + id,
		"exp": time.Now().Unix() + 30 * 86400,
	})

	if tokenString, err := token.SignedString([]byte(config.GetConfig().GetString("service_sceret"))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": tokenString})
	}

	c.Abort()
	return
}

func (controller *UserController) GetUserByIds(_ids []string) (int, gin.H){
	var ids []gocql.UUID
	for _, idString := range _ids {
		if id, err := gocql.ParseUUID(idString); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return http.StatusBadRequest, gin.H{"message": "List id of user not found"}
	}

	if users, err := controller.UserRepo.FindByIDs(ids); err != nil {
		return http.StatusInternalServerError, gin.H{"message": "cannot find user via id"}
	} else {
		return http.StatusOK, gin.H{"message": "success", "data": users}
	}
}