package delivery

import (
	"net/http"
	conf "parkDB/config"
	"parkDB/models"
	"parkDB/usecase"

	"github.com/mailru/easyjson"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserURL     string
	UserUsecase usecase.UserUsecaseInterface
}

func NewUserHandler(router *gin.RouterGroup, userURL string, userUsecase usecase.UserUsecaseInterface) {
	handler := &UserHandler{
		UserURL:     userURL,
		UserUsecase: userUsecase,
	}

	users := router.Group(handler.UserURL)
	{
		users.POST("/:nickname/create", handler.CreateUser)
		users.GET("/:nickname/profile", handler.GetUser)
		users.POST("/:nickname/profile", handler.UpdateUser)
	}
}

func (userHandler *UserHandler) CreateUser(c *gin.Context) {
	nickname := c.Param("nickname")

	userUpdate := new(models.UserUpdate)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, userUpdate); err != nil {
		models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := &models.User{
		Nickname: nickname,
		Fullname: userUpdate.Fullname,
		About:    userUpdate.About,
		Email:    userUpdate.Email,
	}

	users, err := userHandler.UserUsecase.CreateUser(user)
	if err != nil {
		if err == conf.ConflictError {
			usersJSON, errInt := users.MarshalJSON()
			if errInt != nil {
				models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
				return
			}
			c.Data(conf.GetErrorCode(err), conf.Headers, usersJSON)
		} else {
			models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		}
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusCreated, "application/json; charset=utf-8", userJSON)
}

func (userHandler *UserHandler) GetUser(c *gin.Context) {
	nickname := c.Param("nickname")

	user, err := userHandler.UserUsecase.GetUser(nickname)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
}

func (userHandler *UserHandler) UpdateUser(c *gin.Context) {
	nickname := c.Param("nickname")

	userUpdate := new(models.UserUpdate)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, userUpdate); err != nil {
		user, err := userHandler.UserUsecase.GetUser(nickname)
		if err != nil {
			models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
			return
		}

		userJSON, err := user.MarshalJSON()
		if err != nil {
			models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
			return
		}

		c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
		return
	}

	user := &models.User{
		Nickname: nickname,
		Fullname: userUpdate.Fullname,
		About:    userUpdate.About,
		Email:    userUpdate.Email,
	}

	err := userHandler.UserUsecase.UpdateUser(user)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, userJSON)
}
