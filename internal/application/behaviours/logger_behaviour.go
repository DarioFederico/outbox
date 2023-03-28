package behaviours

import (
	"context"
	"reflect"

	"outbox/internal/infrastructure/log"
)

type loggerBehaviour struct {
}

func NewLoggerBehaviour() PipelineBehavior {
	return &loggerBehaviour{}
}

func (r *loggerBehaviour) Handle(c context.Context, request interface{}, next func(context.Context) (interface{}, error)) (interface{}, error) {
	requestType := reflect.TypeOf(request).String()
	ctx := log.WithValue(c, "handler", requestType)
	log.For(ctx).Infof("initialize process request %s with values %+v", requestType, request)

	response, err := next(ctx)
	if err != nil {
		log.For(ctx).Errorf("an error was occurred when process request %s. %+v", requestType, err)
		return nil, err
	}

	log.For(ctx).Infof("finalize request %s with response %+v", requestType, response)
	return response, nil
}
