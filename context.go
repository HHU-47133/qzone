package qzone

// TODO: 没想好怎么设计...

import "math"

const abortIndex int8 = math.MaxInt8 >> 1

type HandlerFunc func(c *Context)

type HandlersChain []HandlerFunc

type Context struct {
	handlers HandlersChain
	index    int8
	// Manager pointer
	manager *Manager
}

func newContext() *Context {
	return &Context{
		index: -1,
	}
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}
