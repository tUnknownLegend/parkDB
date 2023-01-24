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

type ThreadHandler struct {
	ThreadURL     string
	ThreadUsecase usecase.ThreadUsecaseInterface
}

func NewThreadHandler(router *gin.RouterGroup, threadURL string, threadUsecase usecase.ThreadUsecaseInterface) {
	handler := &ThreadHandler{
		ThreadURL:     threadURL,
		ThreadUsecase: threadUsecase,
	}

	threads := router.Group(handler.ThreadURL)
	{
		threads.POST("/:slug_or_id/create", handler.CreatePosts)
		threads.GET("/:slug_or_id/details", handler.GetDetails)
		threads.POST("/:slug_or_id/details", handler.UpdateDetails)
		threads.GET("/:slug_or_id/posts", handler.GetThreadPosts)
		threads.POST("/:slug_or_id/vote", handler.Vote)
	}
}

func (threadHandler *ThreadHandler) CreatePosts(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

	posts := new(models.Posts)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, posts); err != nil {
		models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := threadHandler.ThreadUsecase.CreatePosts(slugOrID, posts)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	postsJSON, err := posts.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusCreated, conf.Headers, postsJSON)
}

func (threadHandler *ThreadHandler) GetDetails(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

	thread, err := threadHandler.ThreadUsecase.GetPost(slugOrID)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, threadJSON)
}

func (threadHandler *ThreadHandler) UpdateDetails(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

	threadUpdate := new(models.ThreadUpdate)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, threadUpdate); err != nil {
		models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	thread := &models.Thread{
		Title:   threadUpdate.Title,
		Message: threadUpdate.Message,
	}
	err := threadHandler.ThreadUsecase.UpdatePost(slugOrID, thread)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, threadJSON)
}

func (threadHandler *ThreadHandler) GetThreadPosts(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

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
	sinceStr := c.Query("since")
	since := -1
	if sinceStr != "" {
		var err error
		since, err = strconv.Atoi(sinceStr)
		if err != nil {
			models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	sort := c.Query("sort")
	if sort == "" {
		sort = "flat"
	}
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

	posts, err := threadHandler.ThreadUsecase.GetPosts(slugOrID, limit, since, sort, desc)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	postsJSON, err := posts.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, postsJSON)
}

func (threadHandler *ThreadHandler) Vote(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

	vote := new(models.Vote)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, vote); err != nil {
		models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	thread, err := threadHandler.ThreadUsecase.Vote(slugOrID, vote)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, threadJSON)
}
