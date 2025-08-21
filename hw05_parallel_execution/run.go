package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type taskRunner struct {
	errorsCount   int32
	runTasksCount int32
	tasksCount    int
}

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	tr := &taskRunner{
		tasksCount: len(tasks),
	}
	tasksChan := make(chan Task, len(tasks))
	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	go func() {
		fmt.Printf("[PRODUCER] Started, total tasks: %d\n", len(tasks))
		start := time.Now()
		for i, task := range tasks {
			tasksChan <- task
			fmt.Printf("[PRODUCER] Sent task %d\n", i+1)
		}
		close(tasksChan)
		fmt.Printf("[PRODUCER] Finished in %v\n", time.Since(start))
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("[WORKER %02d] Started\n", workerID)
			defer fmt.Printf("[WORKER %02d] Stopped\n", workerID)
			for task := range tasksChan {
				if atomic.LoadInt32(&tr.errorsCount) >= int32(m) {
					fmt.Printf("[WORKER %02d] Error limit reached, stopping\n", workerID)
					return
				}
				taskStart := time.Now()
				fmt.Printf("[WORKER %02d] Started task\n", workerID)
				err := task()
				atomic.AddInt32(&tr.runTasksCount, 1)
				duration := time.Since(taskStart)
				if err != nil {
					atomic.AddInt32(&tr.errorsCount, 1)
					fmt.Printf("[WORKER %02d] Task failed (%v) after %v\n", workerID, err, duration)
				} else {
					fmt.Printf("[WORKER %02d] Task completed in %v\n", workerID, duration)
				}
			}
		}(i + 1)
	}

	go func() {
		wg.Wait()
		close(doneChan)
		fmt.Println("[RUNNER] All workers finished")
	}()
	<-doneChan
	if atomic.LoadInt32(&tr.errorsCount) >= int32(m) {
		fmt.Printf("[RUNNER] Error limit exceeded (%d/%d)\n", tr.errorsCount, m)
		return ErrErrorsLimitExceeded
	}
	fmt.Printf("[RUNNER] All tasks completed (%d/%d)\n", tr.runTasksCount, tr.tasksCount)
	return nil
}
