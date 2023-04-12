package provider

import (
    "sync"
    
    "github.com/pkg/errors"
    
    "github.com/go4s/configuration"
)

type (
    Factory[T any]  func(configuration.Configuration) (T, error)
    Provider[T any] interface {
        Implement(string, Factory[T])
        New(env configuration.Configuration) (T, error)
    }
    provider[T any] struct {
        driverTypeKey string
        registry      *sync.Map
    }
)

func New[T any](key string) Provider[T] {
    return &provider[T]{key, new(sync.Map)}
}

func (p *provider[T]) Implement(key string, fn Factory[T]) { p.registry.Store(key, fn) }
func (p *provider[T]) New(env configuration.Configuration) (T, error) {
    var (
        key            = stringOrDefault(env, p.driverTypeKey)
        fn             interface{}
        factory        Factory[T]
        found, support bool
    )
    if len(key) == 0 {
        return *new(T), errors.WithMessagef(ErrEnvironmentNotFound, "key : %s", p.driverTypeKey)
    }
    if fn, found = p.registry.Load(key); !found {
        return *new(T), errors.WithMessagef(ErrClassNotFound, "driver name : %s", key)
    }
    if factory, support = fn.(Factory[T]); !support {
        return *new(T), errors.WithMessagef(ErrClassNotImported, "driver name : %s", key)
    }
    return factory(env)
}

func stringOrDefault(env configuration.Configuration, driverKey string) string {
    if driverName, imported := env[driverKey]; imported {
        return driverName.(string)
    }
    return ""
}

var ErrClassNotImported = errors.New("err class not loaded consider")
var ErrClassNotFound = errors.New("err driver not found")
var ErrEnvironmentNotFound = errors.New("err driver environment not found")
