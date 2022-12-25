package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("invalid worker count", func(t *testing.T) {
		err := Run(nil, 0, 1)

		require.ErrorIs(t, err, ErrInvalidWorkerCount)
	})

	t.Run("unlimited errors", func(t *testing.T) {
		var runTaskCount int32

		taskCount := 50
		tasks := make([]Task, 0, taskCount)

		for i := 0; i < taskCount; i++ {
			err := fmt.Errorf("error in task with id - %d", i)
			tasks = append(tasks, func() error {
				defer atomic.AddInt32(&runTaskCount, 1)
				time.Sleep(200 * time.Millisecond)
				return err
			})
		}

		workerCount := 10

		err := Run(tasks, workerCount, 0)

		require.NoError(t, err)
		require.Equal(t, int32(taskCount), runTaskCount, "not all tasks were completed")
	})

	t.Run("tasks without errors but without sleep", func(t *testing.T) {
		workerCount := 5
		tasks := make([]Task, 0, workerCount)

		var runTasksCount int32
		waitChan := make(chan struct{})

		for i := 0; i < workerCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				<-waitChan
				return nil
			})
		}

		errChan := make(chan error, 1)
		go func() {
			errChan <- Run(tasks, workerCount, 1)
		}()

		require.Eventually(t, func() bool {
			return int32(workerCount) == atomic.LoadInt32(&runTasksCount)
		}, time.Second*3, time.Millisecond, "tasks were run sequentially?")

		close(waitChan)

		var err error
		require.Eventually(t, func() bool {
			select {
			case err = <-errChan:
				return true
			default:
				return false
			}
		}, time.Second, time.Millisecond*100)

		require.NoError(t, err, "error was received")
	})
}
