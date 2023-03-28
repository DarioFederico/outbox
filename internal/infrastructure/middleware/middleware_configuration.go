package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	xRequestId   string = "x-request-id"
	xOperationId string = "x-operation-id"
	traceId      string = "trace_id"
)

type Configuration interface {
	Configure()
}

type middleware struct {
	engine *gin.Engine
}

func NewMiddleware(engine *gin.Engine) Configuration {
	return &middleware{engine: engine}
}

func (c *middleware) Configure() {
	c.setCorrelationHeader()
	c.setTracer()
}

func (c *middleware) setCorrelationHeader() {
	c.engine.Use(func(c *gin.Context) {
		operationId := uuid.New().String()
		requestId := c.GetHeader(xRequestId)
		if len(requestId) == 0 {
			requestId = uuid.New().String()
		}
		c.Header(xRequestId, requestId)
		c.Set(xRequestId, requestId)
		c.Header(xOperationId, operationId)
		c.Set(xOperationId, operationId)
		c.Next()
	})
}

func (c *middleware) setTracer() {
	c.engine.Use(func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set(traceId, traceID)
		c.Next()
	})
}
