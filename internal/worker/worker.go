package worker

import (
	"context"
	"time"

	"github.com/ajugalushkin/gofer-mart/internal/accrualclient"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
	"github.com/ajugalushkin/gofer-mart/internal/queue"
	"github.com/ajugalushkin/gofer-mart/internal/service"
	"github.com/ajugalushkin/gofer-mart/internal/storage"
)

func doWork(ctx context.Context) {
	order, ok := queue.FetchOrder()
	if !ok {
		logger.LogFromContext(ctx).Debug("worker.doWork FetchOrder() failed")
		time.Sleep(100 * time.Millisecond)
		return
	}

	newAccrual, err := accrualclient.GetAccrual(ctx, order.Number)
	if err != nil {
		queue.AddOrder(order)
		logger.LogFromContext(ctx).Debug("worker.doWork accrualclient.GetAccrual() failed")
		time.Sleep(100 * time.Millisecond)
		return
	}

	db := storage.ConnectFromContext(ctx)
	newService := service.NewService(db)

	err = newService.UpdateOrder(ctx, dto.Order{
		Number:     newAccrual.Order,
		UploadedAt: time.Time{},
		Status:     newAccrual.Status,
		Accrual:    newAccrual.Accrual,
	})
	if err != nil {
		queue.AddOrder(order)
		logger.LogFromContext(ctx).Debug("worker.doWork storage.UpdateOrder() failed")
		time.Sleep(100 * time.Millisecond)
		return
	}

	if newAccrual.Status != "INVALID" && newAccrual.Status != "PROCESSED" {
		logger.LogFromContext(ctx).Debug("worker.doWork add order to queue")
		queue.AddOrder(order)
	}
}

func Start(ctx context.Context) {
	go func() {
		for {
			doWork(ctx)
		}
	}()
}
