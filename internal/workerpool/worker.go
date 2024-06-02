package workerpool

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
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

func (w *Worker) StartBackground(ctx context.Context) {
	go func() {
		for url := range w.taskQueue {
			data, err := process(ctx, url)
			w.resultChan <- Result{WorkerID: w.ID, Data: data, Err: err}
		}
	}()
}

func process(ctx context.Context, URL string) (dto.Accrual, error) {
	var accrual dto.Accrual

	req := resty.New()

	req.GetClient().Transport.(*http.Transport).DialContext = nil

	resp, err := req.
		R().
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetResult(&accrual).
		Get(URL)

	if err != nil {
		logger.LogFromContext(ctx).Debug(
			"worker.process: http request error",
			zap.Error(err),
			zap.String("Status", strconv.Itoa(resp.StatusCode())),
			zap.String("Status", resp.Status()),
		)
		return accrual, err
	}

	logger.LogFromContext(ctx).Debug(
		"worker.process: http response Status",
		zap.String("StatusCode", strconv.Itoa(resp.StatusCode())),
		zap.String("Status", resp.Status()))

	switch resp.StatusCode() {
	case http.StatusNoContent:
		err := errors.New("worker.process: status no content")
		logger.LogFromContext(ctx).Debug("worker.process: status no content", zap.Error(err))
		return accrual, err
	case http.StatusTooManyRequests:
		delay, err := strconv.ParseInt(resp.Header().Get(echo.HeaderRetryAfter), 10, 64)
		if err != nil {
			logger.LogFromContext(ctx).Debug("worker.process: Error parsing delay in header", zap.Error(err))
			return accrual, err
		}

		logger.LogFromContext(ctx).Debug(
			"worker.process: Sleeping",
			zap.String("delay", strconv.FormatInt(delay, 10)),
		)
		time.Sleep(time.Duration(delay) * time.Second)
	case http.StatusInternalServerError:
		err := errors.New("worker.process: internal server error")
		logger.LogFromContext(ctx).Debug("worker.process: internal server error", zap.Error(err))
		return accrual, err
	}

	return accrual, nil
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
