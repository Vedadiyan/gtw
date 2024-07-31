package gtw

import (
	"time"

	"github.com/vedadiyan/gtw/internal/di"
)

type (
	Service[T any] struct {
		name     string
		hasScope bool
		scopeId  uint64
		ttl      time.Duration
	}
)

func AddSingleton[T any](fn func() (instance *T, err error)) {
	di.AddSinleton[T](fn)
}

func AddSingletonWithName[T any](name string, fn func() (instance *T, err error)) {
	di.AddSinletonWithName[T](name, fn)
}

func AddTransient[T any](fn func() (instance *T, err error)) {
	di.AddTransient[T](fn)
}

func AddTransientWithName[T any](name string, fn func() (instance *T, err error)) {
	di.AddTransientWithName[T](name, fn)
}

func AddScoped[T any](fn func() (instance *T, err error)) {
	di.AddScoped[T](fn)
}

func AddScopedWithName[T any](name string, fn func() (instance *T, err error)) {
	di.AddScopedWithName[T](name, fn)
}

func (i *Service[T]) Value() *T {
	var options *di.Options
	if i.hasScope {
		options = di.NewOptions(i.scopeId, i.ttl)
	}
	if len(i.name) == 0 {
		return di.ResolveOrPanic[T](options)
	}
	return di.ResolveWithNameOrPanic[T](i.name, options)
}

func (i *Service[T]) Scope(scopeId uint64, ttl time.Duration) *Service[T] {
	copy := *i
	copy.hasScope = true
	copy.scopeId = scopeId
	copy.ttl = ttl
	return &copy
}