package delivery

import (
	"net/http"
	conf "parkDB/config"
	"parkDB/models"
	"parkDB/usecase"
	"strconv"
	"strings"

	"github.com/mailru/easyjson"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	PostURL     string
	PostUsecase usecase.PostUsecaseInterface
}

func NewPostHandler(router *gin.RouterGroup, postURL string, postUsecase usecase.PostUsecaseInterface) {
	handler := &PostHandler{
		PostURL:     postURL,
		PostUsecase: postUsecase,
	}

	posts := router.Group(handler.PostURL)
	{
		posts.GET("/:id/details", handler.GetPost)
		posts.POST("/:id/details", handler.UpdatePost)
	}
}

func (postHandler *PostHandler) GetPost(c *gin.Context) {
	postIDstr := c.Param("id")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		models.GetErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	relatedData := c.Query("related")
	var relatedDataArr []string
	if relatedData != "" {
		relatedDataArr = strings.Split(relatedData, ",")
	}

	postFull, err := postHandler.PostUsecase.Get(int64(postID), &relatedDataArr)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	postFullJSON, err := postFull.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, postFullJSON)
}

func (postHandler *PostHandler) UpdatePost(c *gin.Context) {
	postIDstr := c.Param("id")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	postUpdate := new(models.PostUpdate)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, postUpdate); err != nil {
		models.GetErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	post := &models.Post{
		ID:      int64(postID),
		Message: postUpdate.Message,
	}
	err = postHandler.PostUsecase.Update(post)
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	postJSON, err := post.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, postJSON)
}
