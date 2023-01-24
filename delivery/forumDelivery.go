package delivery

import (
	"net/http"
	conf "parkDB/config"
	"parkDB/models"
	"parkDB/usecase"

	"strconv"

	"github.com/mailru/easyjson"

	"github.com/gin-gonic/gin"
)

type ForumHandler struct {
	ForumURL     string
	ForumUsecase usecase.ForumUsecaseInterface
}

func NewForumHandler(router *gin.RouterGroup, forumURL string, forumUsecase usecase.ForumUsecaseInterface) {
	handler := &ForumHandler{
		ForumURL:     forumURL,
		ForumUsecase: forumUsecase,
	}

	forums := router.Group(handler.ForumURL)
	{
		forums.POST("/create", handler.CreateForum)
		forums.GET("/:slug/details", handler.GetDetails)
		forums.POST("/:slug/create", handler.CreateThread)
		forums.GET("/:slug/users", handler.GetForumUsers)
		forums.GET("/:slug/threads", handler.GetForumThreads)
	}
}

func (forumHandler *ForumHandler) CreateForum(c *gin.Context) {
	forum := new(models.Forum)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, forum); err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	err := forumHandler.ForumUsecase.CreateForum(forum)
	if err != nil {
		if err == conf.ConflictError {
			forumJSON, errInt := forum.MarshalJSON()
			if errInt != nil {
				models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
				return
			}
			c.Data(conf.GetErrorCode(err), conf.Headers, forumJSON)
		} else {
			models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		}
		return
	}

	forumJSON, err := forum.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusCreated, conf.Headers, forumJSON)
}

func (forumHandler *ForumHandler) GetDetails(c *gin.Context) {
	slug := c.Param("slug")

	forum, err := forumHandler.ForumUsecase.GetForumBySlug(slug)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	forumJSON, err := forum.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, forumJSON)
}

func (forumHandler *ForumHandler) CreateThread(c *gin.Context) {
	slug := c.Param("slug")

	thread := new(models.Thread)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, thread); err != nil {
		models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	thread.Forum = slug

	err := forumHandler.ForumUsecase.CreateThread(thread)
	if err != nil {
		if err == conf.ConflictError {
			threadJSON, errInt := thread.MarshalJSON()
			if errInt != nil {
				models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
				return
			}
			c.Data(conf.GetErrorCode(err), conf.Headers, threadJSON)
		} else {
			models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		}
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusCreated, conf.Headers, threadJSON)
}

func (forumHandler *ForumHandler) GetForumUsers(c *gin.Context) {
	slug := c.Param("slug")

	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	since := c.Query("since")
	descStr := c.Query("desc")
	desc := false
	if descStr != "" {
		var err error
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	users, err := forumHandler.ForumUsecase.GetUsersOfForum(slug, limit, since, desc)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	usersJSON, err := users.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, usersJSON)
}

func (forumHandler *ForumHandler) GetForumThreads(c *gin.Context) {
	slug := c.Param("slug")

	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	since := c.Query("since")
	descStr := c.Query("desc")
	desc := false
	if descStr != "" {
		var err error
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	threads, err := forumHandler.ForumUsecase.GetThreadsOfForum(slug, limit, since, desc)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	threadsJSON, err := threads.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, threadsJSON)
}
