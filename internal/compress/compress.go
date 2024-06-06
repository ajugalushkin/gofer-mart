package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	Skipper func(c echo.Context) bool

	GzipConfig struct {
		Skipper Skipper
	}
)

var (
	DefaultGzipConfig = GzipConfig{
		Skipper: func(c echo.Context) bool {
			return false
		},
	}
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func Gzip() echo.MiddlewareFunc {
	return GzipWithConfig(DefaultGzipConfig)
}

func GzipWithConfig(config GzipConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultGzipConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			if config.Skipper(context) {
				return next(context)
			}

			ow := context.Response().Writer

			acceptEncoding := context.Request().Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := newCompressWriter(context.Response().Writer)
				ow = cw
				defer cw.Close()
			}

			contentEncoding := context.Request().Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(context.Request().Body)
				if err != nil {
					context.Response().WriteHeader(http.StatusInternalServerError)
					return nil
				}

				context.Request().Body = cr
				defer cr.Close()
			}

			context.Response().Writer = ow
			if err := next(context); err != nil {
				context.Error(err)
			}
			return nil
		}
	}
}
