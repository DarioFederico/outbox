package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Configuration interface {
	Configure()
}

type healthcheck struct {
	engine *gin.Engine
}

func NewHealthcheck(engine *gin.Engine) *healthcheck {
	return &healthcheck{engine: engine}
}

func (c *healthcheck) ConfigureHealthCheck() {
	// HealthCheck
	c.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Healthy!")
	})
}
