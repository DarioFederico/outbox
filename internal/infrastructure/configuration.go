package infrastructure

import (
	"database/sql"

	"outbox/config"
	"outbox/internal/infrastructure/db"
	"outbox/internal/infrastructure/healthcheck"
	"outbox/internal/infrastructure/mbaas"
	"outbox/internal/infrastructure/middleware"
	"outbox/internal/infrastructure/swagger"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	Cfg      *config.AppConfig
	Engine   *gin.Engine
	MySql    *sql.DB
	MBClient mbaas.MBaaS
}

type Configurator interface {
	Configure() *Configuration
}

type configurator struct {
	cfg    *config.AppConfig
	engine *gin.Engine
}

func NewConfigurator(engine *gin.Engine, cfg *config.AppConfig) Configurator {
	return &configurator{engine: engine, cfg: cfg}
}

func (c *configurator) Configure() *Configuration {
	c.configureMiddleware()
	c.configureHealthcheck()
	c.configureSwagger()
	c.configureLog()

	return &Configuration{
		MySql:    c.configureDb(),
		MBClient: c.configureRabbit(),
		Engine:   c.engine,
		Cfg:      c.cfg,
	}
}

func (c *configurator) configureDb() *sql.DB {
	mysql := db.NewMySql(c.cfg)
	db, err := mysql.Configure()
	if err != nil {
		//panic(err)
	}
	return db
}

func (c *configurator) configureMiddleware() {
	mid := middleware.NewMiddleware(c.engine)
	mid.Configure()
}

func (c *configurator) configureHealthcheck() {
	hc := healthcheck.NewHealthcheck(c.engine)
	hc.ConfigureHealthCheck()
}

func (c *configurator) configureSwagger() {
	sw := swagger.NewSwagger(c.engine, c.cfg)
	sw.Configure()
}

func (c *configurator) configureLog() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
}

func (c *configurator) configureRabbit() mbaas.MBaaS {
	rabbit := mbaas.NewMbaas(c.cfg)
	if err := rabbit.Connect(); err != nil {
		panic(err)
	}
	return rabbit
}
