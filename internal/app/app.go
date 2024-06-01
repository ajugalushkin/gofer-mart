package app

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/compress"
	"github.com/ajugalushkin/gofer-mart/internal/cookies"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
	"github.com/ajugalushkin/gofer-mart/internal/service"
	"github.com/ajugalushkin/gofer-mart/internal/storage"
	"github.com/ajugalushkin/gofer-mart/internal/userrors"
)

type App struct {
	ctx     context.Context
	cache   map[string]dto.User
	service *service.Service
}

const cookieName string = "User"

func NewApp(ctx context.Context, db *sqlx.DB) *App {
	cfg := config.FlagsFromContext(ctx)

	ctx = storage.ContextWithConnect(ctx, db)

	log, _ := logger.Initialize(cfg.LogLevel)
	ctx = logger.ContextWithLogger(ctx, log)

	return &App{
		ctx,
		make(map[string]dto.User),
		service.NewService(db)}
}

func (a App) Routes(r *echo.Echo) {
	r.POST("/api/user/register", a.register)
	r.POST("/api/user/login", a.login)

	r.POST("/api/user/orders", a.authorized(a.postOrders))
	r.GET("/api/user/orders", a.authorized(a.getOrders))

	r.POST("/api/user/accrual/withdraw", a.authorized(a.postBalanceWithdraw))
	r.GET("/api/user/accrual", a.authorized(a.getBalance))

	r.GET("/api/user/withdrawal", a.authorized(a.getWithdrawals))
	r.POST("/api/user/balance/withdraw", a.authorized(a.postBalanceWithdraw))

	//Middleware
	r.Use(logger.MiddlewareLogger(a.ctx))
	r.Use(compress.GzipWithConfig(compress.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger") ||
				strings.Contains(c.Request().URL.Path, "debug")
		},
	}))
}

func (a App) register(echoCtx echo.Context) error {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	loginData := dto.User{}
	err = loginData.UnmarshalJSON(body)
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	err = a.service.AddNewUser(a.ctx, dto.User{
		Login:    loginData.Login,
		Password: loginData.Password,
	})
	if err != nil {
		if errors.Is(err, userrors.ErrorDuplicateLogin) {
			return echoCtx.JSON(http.StatusConflict, err.Error())
		}
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := cookies.Create(a.ctx, cookieName, loginData)
	a.cache[cookie.Value] = loginData
	echoCtx.SetCookie(cookie)

	return echoCtx.JSON(http.StatusOK, "")
}

func (a App) login(echoCtx echo.Context) error {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	loginData := dto.User{}
	err = loginData.UnmarshalJSON(body)
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	err = a.service.LoginUser(a.ctx, dto.User{
		Login:    loginData.Login,
		Password: loginData.Password,
	})
	if err != nil {
		if errors.Is(err, userrors.ErrorIncorrectLoginPassword) {
			return echoCtx.JSON(http.StatusUnauthorized, err.Error())
		}
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}
	cookie := cookies.Create(a.ctx, cookieName, loginData)
	a.cache[cookie.Value] = loginData
	echoCtx.SetCookie(cookie)

	return echoCtx.JSON(http.StatusOK, "")
}

func (a App) authorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		cookie, err := echoCtx.Cookie(cookieName)
		if err != nil {
			return echoCtx.JSON(http.StatusUnauthorized, err.Error())
		}

		if _, ok := a.cache[cookie.Value]; !ok {
			return echoCtx.JSON(http.StatusUnauthorized, "")
		}

		return next(echoCtx)
	}
}

func (a App) postOrders(echoCtx echo.Context) error {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	order := string(body)
	err = goluhn.Validate(order)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	cookie, err := echoCtx.Cookie(cookieName)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}
	login := cookies.GetLogin(a.ctx, cookie.Value)

	err = a.service.AddNewOrder(a.ctx, order, login.Login)
	if err != nil {
		if errors.Is(err, userrors.ErrorOrderAlreadyUploadedAnotherUser) {
			return echoCtx.JSON(http.StatusConflict, err.Error())
		} else if errors.Is(err, userrors.ErrorOrderAlreadyUploadedThisUser) {
			return echoCtx.JSON(http.StatusOK, err.Error())
		}
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusAccepted, "")
}

func (a App) getOrders(echoCtx echo.Context) error {
	cookie, err := echoCtx.Cookie(cookieName)
	if err != nil {
		logger.LogFromContext(a.ctx).Debug("app.getOrders: Error getting cookie data")
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}
	login := cookies.GetLogin(a.ctx, cookie.Value)
	logger.LogFromContext(a.ctx).Debug("get Login", zap.String("login", login.Login))

	orderList, err := a.service.GetOrders(a.ctx, login.Login)
	if err != nil {
		logger.LogFromContext(a.ctx).Debug("app.getOrders: Error getting order list")
		return echoCtx.JSON(http.StatusNoContent, err.Error())
	}
	logger.LogFromContext(a.ctx).Debug("app.getOrders: Ok")
	return echoCtx.JSON(http.StatusOK, orderList)
}

func (a App) getBalance(echoCtx echo.Context) error {
	cookie, err := echoCtx.Cookie(cookieName)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}
	login := cookies.GetLogin(a.ctx, cookie.Value)

	balance, err := a.service.GetBalance(a.ctx, login.Login)
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}
	return echoCtx.JSON(http.StatusOK, balance)
}

func (a App) postBalanceWithdraw(echoCtx echo.Context) error {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	cookie, err := echoCtx.Cookie(cookieName)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}
	login := cookies.GetLogin(a.ctx, cookie.Value)

	withdraw := dto.Withdraw{}
	err = withdraw.UnmarshalJSON(body)
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	err = a.service.AddNewWithdrawal(a.ctx, withdraw, login.Login)
	if err != nil {
		if errors.Is(err, userrors.ErrorInsufficientFunds) {
			return echoCtx.JSON(http.StatusPaymentRequired, err.Error())
		}
		if errors.Is(err, userrors.ErrorIncorrectOrderNumber) {
			return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
		}
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}
	return echoCtx.JSON(http.StatusAccepted, "")
}

func (a App) getWithdrawals(echoCtx echo.Context) error {
	cookie, err := echoCtx.Cookie(cookieName)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}
	login := cookies.GetLogin(a.ctx, cookie.Value)

	list, err := a.service.GetWithdrawalList(a.ctx, login.Login)
	if err != nil {
		return echoCtx.JSON(http.StatusNoContent, err.Error())
	}
	return echoCtx.JSON(http.StatusAccepted, list)
}
