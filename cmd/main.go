package main

import (
	"fmt"
	"log"

	"outbox/config"
	"outbox/internal/application"
	"outbox/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	//Configure infrastructure
	conf := config.GetConf()
	engine := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	configurator := infrastructure.NewConfigurator(engine, conf)
	infCfg := configurator.Configure()

	//Configure application
	app := application.NewApplication(infCfg)
	if err := app.Setup(); err != nil {
		panic(err)
	}

	// Start server
	if err := engine.Run(fmt.Sprintf(":%s", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
