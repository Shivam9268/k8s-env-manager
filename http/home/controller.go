package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (h *Controller) HomeHandler(c *gin.Context) {
	var envVars = map[string]string{}
	for _, env := range os.Environ() {
		envPair := strings.SplitN(env, "=", 2)
		key := envPair[0]
		value := envPair[1]
		envVars[key] = value
	}
	c.JSON(http.StatusOK, envVars)
}
