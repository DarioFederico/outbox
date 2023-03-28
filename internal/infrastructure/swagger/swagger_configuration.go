package swagger

import (
	"fmt"

	"outbox/config"
	_ "outbox/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Configuration interface {
	Configure()
}

type swagger struct {
	engine *gin.Engine
	cfg    *config.AppConfig
}

func NewSwagger(engine *gin.Engine, cfg *config.AppConfig) Configuration {
	return &swagger{engine: engine, cfg: cfg}
}

func (c *swagger) Configure() {
	c.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", c.cfg.Port)),
		ginSwagger.DefaultModelsExpandDepth(-1)))
}
