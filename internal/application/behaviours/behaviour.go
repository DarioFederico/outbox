package behaviours

import "context"

type PipelineBehavior interface {
	Handle(ctx context.Context, request interface{}, next func(context.Context) (interface{}, error)) (interface{}, error)
}
