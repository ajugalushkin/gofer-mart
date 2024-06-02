package workerpool

import (
	"context"
	"fmt"
	"time"

	"github.com/ajugalushkin/gofer-mart/internal/dto"
)

type Result struct {
	WorkerID int
	Order    string
	Data     dto.Accrual
	Err      error
}

type WorkerPool struct {
	Workers       []*Worker
	taskQueue     chan string
	resultChan    chan Result
	workerCount   int
	runBackground chan bool
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		taskQueue:   make(chan string),
		resultChan:  make(chan Result),
		workerCount: workerCount,
	}
}

func (wp *WorkerPool) RunBackground(ctx context.Context) {
	go func() {
		for {
			fmt.Print("⌛ Waiting for tasks to come in ...\n")
			time.Sleep(10 * time.Second)
		}
	}()

	for i := 1; i <= wp.workerCount; i++ {
		worker := NewWorker(wp.taskQueue, i, wp.resultChan)
		wp.Workers = append(wp.Workers, worker)
		go worker.StartBackground(ctx)
	}

	wp.runBackground = make(chan bool)
	<-wp.runBackground
}

func (wp *WorkerPool) AddTask(url string) {
	wp.taskQueue <- url
}

func (wp *WorkerPool) GetResult() Result {
	return <-wp.resultChan
}

func (p *WorkerPool) Stop() {
	for i := range p.Workers {
		p.Workers[i].Stop()
	}
	p.runBackground <- true
}

type ctxWorkerPool struct{}

func ContextWorkerPool(ctx context.Context, pool *WorkerPool) context.Context {
	return context.WithValue(ctx, ctxWorkerPool{}, pool)
}

func WorkerPoolFromContext(ctx context.Context) *WorkerPool {
	if pool, ok := ctx.Value(ctxWorkerPool{}).(*WorkerPool); ok {
		return pool
	}
	return nil
}
