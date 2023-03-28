package commons

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
)

type Request interface {
}

type Response interface {
}

type Notification interface {
}

type RequestHandler interface {
	Handle(ctx context.Context, request Request) (Response, error)
}

type NotificationHandler interface {
	Handle(ctx context.Context, notification Notification) error
}

type Mediator interface {
	Send(ctx context.Context, request Request) (Response, error)
	RegisterRequestHandler(requestType Request, handler RequestHandler) error
	RegisterRequestPipelineBehaviors(behaviours ...pipelineBehavior) error
}

// requestHandlerFunc is a continuation for the next task to execute in the pipeline
type requestHandlerFunc func(ctx context.Context) (interface{}, error)

// pipelineBehavior is a Pipeline behavior for wrapping the inner handler.
type pipelineBehavior interface {
	Handle(ctx context.Context, request interface{}, next func(context.Context) (interface{}, error)) (interface{}, error)
}

type mediator struct {
	requestHandlersRegistrations      map[string]RequestHandler
	pipelineBehaviours                []pipelineBehavior
	notificationHandlersRegistrations map[string][]interface{}
}

func NewMediator() Mediator {
	return &mediator{
		requestHandlersRegistrations: map[string]RequestHandler{},
	}
}

func (m *mediator) RegisterRequestHandler(requestType Request, handler RequestHandler) error {
	name := reflect.TypeOf(requestType).String()
	_, exist := m.requestHandlersRegistrations[name]
	if exist {
		// each request in request/response strategy should have just one handler
		return errors.Errorf("registered handler already exists in the registry for message %s", name)
	}

	m.requestHandlersRegistrations[name] = handler

	return nil
}

// RegisterRequestPipelineBehaviors register the request behaviors to mediatr registry.
func (m *mediator) RegisterRequestPipelineBehaviors(behaviours ...pipelineBehavior) error {
	for _, behavior := range behaviours {
		behaviorType := reflect.TypeOf(behavior).String()

		if m.existsPipeType(behaviorType) {
			return errors.Errorf("registered behavior already exists in the registry.")
		}
		m.pipelineBehaviours = append(m.pipelineBehaviours, behavior)
	}
	return nil
}

func (m *mediator) Send(ctx context.Context, request Request) (Response, error) {
	requestType := reflect.TypeOf(request).String()
	handler, ok := m.requestHandlersRegistrations[requestType]
	if !ok {
		// request-response strategy should have exactly one handler and if we can't find a corresponding handler, we should return an error
		return nil, errors.Errorf("no handler for request %T", request)
	}

	handlerValue, ok := handler.(RequestHandler)
	if !ok {
		return nil, errors.Errorf("handler for request %T is not a Handler", request)
	}

	if len(m.pipelineBehaviours) > 0 {
		var lastHandler requestHandlerFunc = func(context.Context) (interface{}, error) {
			return handlerValue.Handle(ctx, request)
		}

		wrapped := buildChain(ctx, request, lastHandler, m.pipelineBehaviours)
		res, err := wrapped(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "error handling request")
		}
		return res, nil

	} else {
		res, err := handlerValue.Handle(ctx, request)
		if err != nil {
			return nil, errors.Wrap(err, "error handling request")
		}
		return res, nil
	}
}

// Publish the notification event to its corresponding notification handler.
func (m *mediator) Publish(ctx context.Context, notification Notification) error {
	eventType := reflect.TypeOf(notification).String()

	handlers, ok := m.notificationHandlersRegistrations[eventType]
	if !ok {
		// notification strategy should have zero or more handlers, so it should run without any error if we can't find a corresponding handler
		return nil
	}

	for _, handler := range handlers {
		handlerValue, ok := handler.(NotificationHandler)
		if !ok {
			return errors.Errorf("handler for notification %T is not a Handler", notification)
		}

		err := handlerValue.Handle(ctx, notification)
		if err != nil {
			return errors.Wrap(err, "error handling notification")
		}
	}

	return nil
}

// buildChain builds the middlware chain recursively, functions are first class
func buildChain(ctx context.Context, request Request, f requestHandlerFunc, pipes []pipelineBehavior) requestHandlerFunc {
	// if our chain is done, use the original handlerfunc
	if len(pipes) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return func(context.Context) (interface{}, error) {
		return pipes[0].Handle(ctx, request, buildChain(ctx, request, f, pipes[1:]))
	}
}

func (m *mediator) existsPipeType(p string) bool {
	for _, pipe := range m.pipelineBehaviours {
		if reflect.TypeOf(pipe).String() == p {
			return true
		}
	}
	return false
}
