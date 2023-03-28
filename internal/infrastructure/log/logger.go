package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

func For(ctx context.Context) *logrus.Entry {
	return logrus.WithFields(Values(ctx))
}

// Key defines the key in the context baggage where values are stored as map[string]interface{}
const Key = "logctx-data-map-string-interface"

// baggageKey is the same key but typed as an  interface{}, because lint doesn't like strings as keys
var baggageKey interface{} = Key

// WithValue returns a context with the added key value pair in the baggage store.
func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	oldBaggage, ok := ctx.Value(baggageKey).(map[string]interface{})

	if !ok {
		return context.WithValue(ctx, baggageKey, map[string]interface{}{key: value})
	}

	newBaggage := make(map[string]interface{}, len(oldBaggage)+1)
	for oldbaggageKey, oldValue := range oldBaggage {
		newBaggage[oldbaggageKey] = oldValue
	}
	newBaggage[key] = value

	return context.WithValue(ctx, baggageKey, newBaggage)
}

// WithValues returns a context with all key value pairs added to the baggage store.
func WithValues(ctx context.Context, keyValue map[string]interface{}) context.Context {
	oldBaggage, ok := ctx.Value(baggageKey).(map[string]interface{})
	if !ok {
		return context.WithValue(ctx, baggageKey, map[string]interface{}(keyValue))
	}

	newBaggage := make(map[string]interface{}, len(oldBaggage)+len(keyValue))
	for oldbaggageKey, oldValue := range oldBaggage {
		newBaggage[oldbaggageKey] = oldValue
	}

	for newbaggageKey, newValue := range keyValue {
		newBaggage[newbaggageKey] = newValue
	}

	return context.WithValue(ctx, baggageKey, newBaggage)
}

// NewFrom returns a new context.Background() with baggage values obtained from another context
func NewFrom(ctx context.Context) context.Context {
	oldBaggage, ok := ctx.Value(baggageKey).(map[string]interface{})
	if !ok {
		return context.Background()
	}
	return context.WithValue(context.Background(), baggageKey, oldBaggage)
}

// Values returns the values stored in the baggage, or an empty map if there are none
func Values(ctx context.Context) map[string]interface{} {
	if values, ok := ctx.Value(baggageKey).(map[string]interface{}); ok {
		return values
	}
	return map[string]interface{}{}
}
