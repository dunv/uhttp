package uhttp

import (
	"context"
	"time"

	"github.com/dunv/ulog"
)

type HandlerOption interface {
	apply(*handlerOptions)
}

type handlerOptions struct {
	Get          HandlerFunc
	GetWithModel HandlerFuncWithModel
	GetModel     interface{}

	Post          HandlerFunc
	PostWithModel HandlerFuncWithModel
	PostModel     interface{}

	Delete          HandlerFunc
	DeleteWithModel HandlerFuncWithModel
	DeleteModel     interface{}

	RequiredGet    R
	OptionalGet    R
	Middlewares    []Middleware
	PreProcess     func(ctx context.Context) error
	Timeout        time.Duration
	TimeoutMessage string
}
type funcHandlerOption struct {
	f func(*handlerOptions)
}

func (fdo *funcHandlerOption) apply(do *handlerOptions) {
	fdo.f(do)
}

func newFuncHandlerOption(f func(*handlerOptions)) *funcHandlerOption {
	return &funcHandlerOption{f: f}
}

func WithGet(h HandlerFunc) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.GetWithModel != nil {
			ulog.Error("cannot use WithGetModel in conjunction with WithGet. WithGet will supercede this assignment")
		}

		o.Get = h
	})
}

func WithGetModel(m interface{}, h HandlerFuncWithModel) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.Get != nil {
			ulog.Error("cannot use WithGetModel in conjunction with WithGet. WithGet will supercede this assignment")
		}

		o.GetModel = m
		o.GetWithModel = h
	})
}

func WithPost(h HandlerFunc) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.PostWithModel != nil {
			ulog.Error("cannot use WithPostModel in conjunction with WithPost. WithPost will supercede this assignment")
		}

		o.Post = h
	})
}

func WithPostModel(m interface{}, h HandlerFuncWithModel) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.Post != nil {
			ulog.Error("cannot use WithPostModel in conjunction with WithPost. WithPost will supercede this assignment")
		}

		o.PostModel = m
		o.PostWithModel = h
	})
}

func WithDelete(h HandlerFunc) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.DeleteWithModel != nil {
			ulog.Error("cannot use WithDeleteModel in conjunction with WithDelete. WithDelete will supercede this assignment")
		}

		o.Delete = h
	})
}

func WithDeleteModel(m interface{}, h HandlerFuncWithModel) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.Delete != nil {
			ulog.Error("cannot use WithDeleteModel in conjunction with WithDelete. WithDelete will supercede this assignment")
		}

		o.DeleteModel = m
		o.DeleteWithModel = h
	})
}

func WithRequiredGet(r R) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.RequiredGet = r
	})
}

func WithOptionalGet(r R) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.OptionalGet = r
	})
}

func WithMiddlewares(m ...Middleware) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.Middlewares = m
	})
}

func WithPreProcess(p PreProcessFunc) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.PreProcess = p
	})
}

func WithTimeout(timeout time.Duration, timeoutMessage string) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.Timeout = timeout
		o.TimeoutMessage = timeoutMessage
	})
}
