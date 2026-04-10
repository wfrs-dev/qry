package console

import "sync"

type Context struct {
	data map[string]any
}

var (
	instance *Context
	once     sync.Once
)

func GetContext() *Context {
	once.Do(func() {
		instance = &Context{
			data: make(map[string]any),
		}
	})
	return instance
}

func (c *Context) Set(key string, value any) {
	c.data[key] = value
}

func Get[T any](c *Context, key string) (T, bool) {
	val, ok := c.data[key]
	if !ok {
		var zero T
		return zero, false
	}
	typed, ok := val.(T)
	return typed, ok
}
