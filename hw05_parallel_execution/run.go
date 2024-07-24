package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task, len(tasks))
	errorCh := make(chan error, len(tasks))
	doneCh := make(chan struct{})
	var completedTasks int

	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		go worker(taskCh, errorCh, doneCh, &completedTasks, len(tasks))
	}

	go func() {
		defer func() {
			close(taskCh)
		}()
		for _, task := range tasks {
			select {
			case taskCh <- task:
			case <-doneCh:
				return
			}
		}
	}()

	var errCount int

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(doneCh)
		}()

		for err := range errorCh {
			if err != nil {
				errCount++

				if errCount >= m {
					return
				}
			}
		}
	}()

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(taskCh <-chan Task, errorCh chan<- error, doneCh <-chan struct{}, completedTasks *int, totalTasks int) {
	select {
	case <-doneCh:
		close(errorCh)
		return
	default:
		for task := range taskCh {
			if task != nil {
				err := task()
				errorCh <- err
				*completedTasks++
			}

			if *completedTasks >= totalTasks {
				close(errorCh)

				return
			}
		}
	}
}
