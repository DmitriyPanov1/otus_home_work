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

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("m is more than the count of errors in the tasks", func(t *testing.T) {
		tests := []struct {
			name           string
			workers        int
			maxErrorsCount int
			errorsCount    int
			tasksCount     int
		}{
			{
				"m > errorsCount",
				5,
				5,
				4,
				7,
			},
			{
				"m = errorsCount",
				5,
				5,
				5,
				7,
			},
			{
				"m = 0",
				5,
				0,
				5,
				7,
			},
			{
				"m = 0 and errorCount = 0",
				5,
				0,
				0,
				7,
			},
		}

		for _, tc := range tests {
			tasks := make([]Task, 0, tc.tasksCount)

			var runTasksCount int32
			var countErrors int

			for i := 0; i < tc.tasksCount; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)

					if countErrors < tc.errorsCount {
						countErrors++

						return err
					}

					return nil
				})
			}

			err := Run(tasks, tc.workers, tc.maxErrorsCount)

			switch {
			case tc.maxErrorsCount == 0:
			case tc.errorsCount >= tc.maxErrorsCount:
				require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
			default:
				require.NoError(t, err)
			}

			require.LessOrEqual(t, runTasksCount, int32(tc.workers+tc.maxErrorsCount), "extra tasks were started")
		}
	})
}
