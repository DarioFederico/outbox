package categories

import (
	"log"

	"outbox/internal/application/commons"
	"outbox/internal/application/modules/categories/controllers"
	"outbox/internal/application/modules/categories/handlers/commands"
	"outbox/internal/application/modules/categories/handlers/queries"
	"outbox/internal/application/modules/categories/jobs"
	"outbox/internal/infrastructure"
)

type categoryModule struct {
	infCfg   *infrastructure.Configuration
	mediator commons.Mediator
}

func NewCategoryModule(infCfg *infrastructure.Configuration, mediator commons.Mediator) commons.Module {
	return &categoryModule{infCfg: infCfg, mediator: mediator}
}

func (m *categoryModule) Configure() {
	//Create outbox job
	outboxJob := jobs.NewOutboxJob(m.infCfg.Cfg, m.infCfg.MBClient, m.infCfg.MySql)

	//Create queries
	getCategoryByIdHandler := queries.NewGetCategoryByIdHandler(m.infCfg.MySql)
	err := m.mediator.RegisterRequestHandler(queries.GetCategoryById{}, getCategoryByIdHandler)
	if err != nil {
		log.Fatal(err)
	}

	//Create commands
	createCategoryHandler := commands.NewCreateCategoryCommandHandler(m.infCfg.MySql)
	err = m.mediator.RegisterRequestHandler(commands.CreateCategoryCommand{}, createCategoryHandler)
	if err != nil {
		log.Fatal(err)
	}

	//Create controllers
	controller := controllers.NewCategoryController(m.mediator)

	//Create routes
	m.infCfg.Engine.GET("/categories/:id", controller.GetCategoryById)
	m.infCfg.Engine.POST("/categories", controller.CreateCategory)

	//Run jobs
	outboxJob.Run()
}
