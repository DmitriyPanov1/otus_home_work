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
	wg := sync.WaitGroup{}

	if m == 0 {
		return ErrErrorsLimitExceeded
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

	go func() {
		defer func() {
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

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(taskCh, errorCh, doneCh, &wg)
	}

	wg.Wait()

	close(errorCh)

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(taskCh <-chan Task, errorCh chan<- error, doneCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskCh {
		select {
		case <-doneCh:
			return
		default:
			errorCh <- task()
		}
	}
}
