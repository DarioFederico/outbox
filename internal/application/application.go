package application

import (
	"outbox/internal/application/behaviours"
	"outbox/internal/application/commons"
	"outbox/internal/application/modules/categories"
	"outbox/internal/infrastructure"

	log "github.com/sirupsen/logrus"
)

type Application interface {
	Setup() error
}

type application struct {
	InfCfg *infrastructure.Configuration
}

func NewApplication(infCfg *infrastructure.Configuration) Application {
	return &application{InfCfg: infCfg}
}

func (app *application) Setup() error {
	mediator := commons.NewMediator()

	//Setup behaviours
	loggerPipeline := behaviours.NewLoggerBehaviour()
	metricPipeline := behaviours.NewMetricBehaviour()

	err := mediator.RegisterRequestPipelineBehaviors(loggerPipeline, metricPipeline)
	if err != nil {
		log.Fatalf("an error was occurred when modules was configured. %+v", err)
		return err
	}

	//Setup Modules
	categoryModule := categories.NewCategoryModule(app.InfCfg, mediator)

	setupModules(categoryModule)

	return nil
}

func setupModules(modules ...commons.Module) {
	for _, m := range modules {
		m.Configure()
	}
}
