package accrualclient

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
)

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrualclient"`
}

func GetAccrual(ctx context.Context, number string) (*Accrual, error) {
	var accrual Accrual
	req := resty.New().
		SetBaseURL(config.FlagsFromContext(ctx).AccrualSystemAddress).
		R().
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetResult(&accrual)

	resp, err := req.Get(fmt.Sprintf("/api/orders/%s", number))
	if err != nil {
		logger.LogFromContext(ctx).Debug("accrualclient.GetAccrual Error:",
			zap.Error(err))
		return &accrual, err
	}

	logger.LogFromContext(ctx).Debug("accrualclient.GetAccrual Status:",
		zap.String("Status Code", strconv.Itoa(resp.StatusCode())),
		zap.String("Status", resp.Status()))

	if resp.StatusCode() == http.StatusTooManyRequests {
		delay, err := strconv.ParseInt(resp.Header().Get(echo.HeaderRetryAfter), 10, 64)
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
	return &accrual, nil
}
