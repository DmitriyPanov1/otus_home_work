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
		return errors.New("m must be greater than 0")
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

		for {
			select {
			case err, ok := <-errorCh:

				if !ok {
					return
				}

				if err != nil {
					errCount++

					if errCount >= m {
						return
					}
				}
			}
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {

				select {
				case <-doneCh:
					return
				default:
				}

				err := task()

				if err != nil {
					errorCh <- err
				}
			}

			return

		}()
	}

	wg.Wait()

	close(errorCh)

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
