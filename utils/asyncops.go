package utils

import (
	"errors"
	"time"
)

type AsyncResult struct {
	Result any
	Error  error
}

func ExecuteAsync(task func() (any, error), timeout time.Duration) (*AsyncResult, error) {
	resultChan := make(chan *AsyncResult)

	go func() {
		result, err := task()
		resultChan <- &AsyncResult{Result: result, Error: err}
	}()

	select {
	case res := <-resultChan:
		return res, nil
	case <-time.After(timeout):
		return nil, errors.New("operation timed out")
	}
}
