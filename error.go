package di

import (
	"fmt"
	"sync"
)

type Errorer interface {
	ErrorCh() <-chan error
}

type errorerShutdowner interface {
	shutdownErrorer()
}

type ErrorEmitter struct {
	ch    chan error
	oncer *sync.Once
	mx    *sync.Mutex
}

func NewErrorEmitter() *ErrorEmitter {
	return &ErrorEmitter{
		ch:    make(chan error, 1),
		oncer: new(sync.Once),
		mx:    new(sync.Mutex),
	}
}

func (emitter *ErrorEmitter) ErrorCh() <-chan error {
	return emitter.ch
}

func (emitter *ErrorEmitter) FireError(err error) {
	emitter.mx.Lock()
	defer emitter.mx.Unlock()

	if emitter.oncer != nil {
		emitter.oncer.Do(func() { emitter.ch <- err })
	}
}

func (emitter *ErrorEmitter) shutdownErrorer() {
	emitter.mx.Lock()
	defer emitter.mx.Unlock()

	emitter.oncer = nil
	close(emitter.ch)
}

type ErrUnknownBean struct {
	beanName string
}

func (err ErrUnknownBean) Error() string { return fmt.Sprintf("unkown bean `%s`", err.beanName) }
