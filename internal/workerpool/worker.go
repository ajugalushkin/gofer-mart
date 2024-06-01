package workerpool

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"

	"github.com/ajugalushkin/gofer-mart/internal/dto"
)

type Worker struct {
	ID         int
	taskQueue  <-chan string
	resultChan chan<- Result
	quit       chan bool
}

func NewWorker(channel <-chan string, ID int, resultChan chan<- Result) *Worker {
	return &Worker{
		ID:         ID,
		taskQueue:  channel,
		quit:       make(chan bool),
		resultChan: resultChan,
	}
}

func (w *Worker) StartBackground() {
	go func() {
		for url := range w.taskQueue {
			data, err := process(url)
			w.resultChan <- Result{WorkerID: w.ID, Data: data, Err: err}
		}
	}()
}

func process(URL string) (dto.Accrual, error) {
	var accrual dto.Accrual
	resp, err := resty.New().
		R().
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetResult(accrual).
		Get(URL)

	if err != nil {
		return accrual, err
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		delay, err := strconv.ParseInt(resp.Header().Get(echo.HeaderRetryAfter), 10, 64)
		if err != nil {
			return accrual, err
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}

	return accrual, nil
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
