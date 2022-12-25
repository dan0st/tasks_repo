package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidWorkerCount  = errors.New("worker count should be grater than 0")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m <= 0 means no limit for failed task.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrInvalidWorkerCount
	}
	if m <= 0 {
		m = len(tasks) + 1
	}

	var errorsCount int64
	taskCh := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := task(); err != nil {
					atomic.AddInt64(&errorsCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt64(&errorsCount) >= int64(m) {
			break
		}
		taskCh <- task
	}

	close(taskCh)
	wg.Wait()

	if errorsCount >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
