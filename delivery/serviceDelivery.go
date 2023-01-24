package delivery

import (
	"net/http"
	conf "parkDB/config"
	"parkDB/models"
	"parkDB/usecase"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceURL     string
	ServiceUsecase usecase.ServiceUsecaseInterface
}

func NewServiceHandler(router *gin.RouterGroup, serviceURL string, serviceUsecase usecase.ServiceUsecaseInterface) {
	handler := &ServiceHandler{
		ServiceURL:     serviceURL,
		ServiceUsecase: serviceUsecase,
	}

	service := router.Group(handler.ServiceURL)
	{
		service.POST("/clear", handler.Clear)
		service.GET("/status", handler.GetStatus)
	}
}

func (serviceHandler *ServiceHandler) Clear(c *gin.Context) {
	err := serviceHandler.ServiceUsecase.Clear()
	if err != nil {
		// c.Data(errors.PrepareErrorResponse(err))
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (serviceHandler *ServiceHandler) GetStatus(c *gin.Context) {
	status, err := serviceHandler.ServiceUsecase.GetStatus()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())

		// c.Data(errors.PrepareErrorResponse(err))
		return
	}

	statusJSON, err := status.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())

		// c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, conf.Headers, statusJSON)
}
