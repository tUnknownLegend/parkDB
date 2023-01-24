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
		service.POST("/clear", handler.ClearService)
		service.GET("/status", handler.GetServiceStatus)
	}
}

func (serviceHandler *ServiceHandler) ClearService(c *gin.Context) {
	err := serviceHandler.ServiceUsecase.ClearDB()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (serviceHandler *ServiceHandler) GetServiceStatus(c *gin.Context) {
	status, err := serviceHandler.ServiceUsecase.GetStatusOfDB()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	statusJSON, err := status.MarshalJSON()
	if err != nil {
		models.GetErrorResponse(c, conf.GetErrorCode(err), err.Error())
		return
	}

	c.Data(http.StatusOK, conf.Headers, statusJSON)
}
