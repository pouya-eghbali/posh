package exec

import (
	"bufio"
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
	scanner := bufio.NewScanner(r.Stdout)

	result := ""
	for scanner.Scan() {
		result += scanner.Text() + "\n"
	}

	return result
}

func NewContext() *RunContext {
	return &RunContext{}
}
