package exec

import (
	"context"
	"io"
	"sync"
)

type RunContext struct {
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	Ctx    *context.Context
	Err    error
}

func (r *RunContext) Wait() *RunContext {
	if r.Ctx != nil {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			<-(*r.Ctx).Done()
		}()
	}

	return r
}

func (r *RunContext) ToString() string {
	data, err := io.ReadAll(r.Stdout)
	if err != nil {
		return ""
	}
	return string(data)
}

func NewContext() *RunContext {
	return &RunContext{}
}
