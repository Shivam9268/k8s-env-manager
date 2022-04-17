package health

import (
	"github.com/Shivam9268/k8s-env-manager/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (h *Controller) HealthHandler(c *gin.Context) {
	sessions.HealthOrPanic()
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}
