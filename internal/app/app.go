package app

import (
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/ajugalushkin/gofer-mart/internal/cookies"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/storage"
	"github.com/ajugalushkin/gofer-mart/internal/user_errors"
)

//	func Run(ctx context.Context) error {
//		flags := config.FlagsFromContext(ctx)
//
//		server := echo.New()
//
//		handler := userHandler.NewHandler()
//
//		//Handlers
//		server.POST("/api/user/register", handler.Register)
//		server.POST("/api/user/login", handler.Login)
//		server.POST("/api/user/orders", handler.PostOrders)
//		server.POST("/api/user/balance/withdraw", handler.PostBalanceWithdraw)
//
//		server.GET("/api/user/orders", handler.GetOrders)
//		server.GET("/api/user/balance", handler.GetBalance)
//		server.GET("/api/user/withdrawals", handler.GetWithdrawals)
//
//		log, err := logger.Initialize(flags.FlagLogLevel)
//		if err != nil {
//			return err
//		}
//
//		ctx = logger.ContextWithLogger(ctx, log)
//
//		// Middleware
//		server.Use(logger.MiddlewareLogger(ctx))
//
//		err = server.Start(flags.RunAddr)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
type App struct {
	ctx context.Context
	//repo  *repository.Repository
	//cache map[string]repository.User
}

func NewApp(ctx context.Context) *App {
	return &App{ctx}
}

func (a App) Routes(r *echo.Echo) {
	r.POST("/api/user/register", a.register)
	r.POST("/api/user/login", a.login)
	r.POST("/api/user/orders", a.postOrders)
	r.POST("/api/user/balance/withdraw", a.postBalanceWithdraw)

	r.GET("/api/user/orders", a.getOrders)
	r.GET("/api/user/balance", a.getBalance)
	r.GET("/api/user/withdrawals", a.getWithdrawals)
}

func (a App) register(echoCtx echo.Context) error {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	loginData := dto.Login{}
	err = loginData.UnmarshalJSON(body)
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	err = storage.AddNewUser(a.ctx, dto.User{
		Login:    loginData.Login,
		Password: loginData.Password,
	})
	if err != nil {
		if errors.Is(err, user_errors.ErrorDuplicateLogin) {
			return echoCtx.JSON(http.StatusConflict, err.Error())
		}
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := cookies.Create(a.ctx, "User_registered", loginData.Login)
	echoCtx.SetCookie(cookie)

	return echoCtx.JSON(http.StatusOK, "")
}

func (a App) login(echoCtx echo.Context) error {
	body, err := io.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	loginData := dto.Login{}
	err = loginData.UnmarshalJSON(body)
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	err = storage.LoginUser(a.ctx, dto.User{
		Login:    loginData.Login,
		Password: loginData.Password,
	})
	if err != nil {
		if errors.Is(err, user_errors.ErrorLoginAlreadyTaken) {
			return echoCtx.JSON(http.StatusUnauthorized, err.Error())
		}
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}
	cookie := cookies.Create(a.ctx, "User_login", loginData.Login)
	echoCtx.SetCookie(cookie)

	return echoCtx.JSON(http.StatusOK, "")
}

func (a App) getOrders(echoCtx echo.Context) error {
	return nil
}

func (a App) postOrders(echoCtx echo.Context) error {
	return nil
}

func (a App) getBalance(echoCtx echo.Context) error {
	return nil
}

func (a App) postBalanceWithdraw(echoCtx echo.Context) error {
	return nil
}

func (a App) getWithdrawals(echoCtx echo.Context) error {
	return nil
}
