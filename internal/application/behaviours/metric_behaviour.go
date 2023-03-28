package behaviours

import (
	"context"
	"reflect"
	"time"

	"outbox/internal/infrastructure/log"
)

type metricBehaviour struct {
}

func NewMetricBehaviour() PipelineBehavior {
	return &metricBehaviour{}
}

func (r *metricBehaviour) Handle(ctx context.Context, request interface{}, next func(context.Context) (interface{}, error)) (interface{}, error) {
	start := time.Now()
	response, err := next(ctx)
	elapsed := time.Since(start)
	requestType := reflect.TypeOf(request).String()
	log.For(ctx).Infof("%s elapsed %s", requestType, elapsed)

	if err != nil {
		return nil, err
	}

	return response, nil
}
