package application

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

type statusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

type Controller struct {
	Application application
}

type RestController struct {
	Controller
}

// NewController creates and returns controller
func NewController(a application) (Controller, error) {
	controller := Controller{Application: a}
	return controller, nil
}

// NewRestController creates and returns restController
func NewRestController(controller Controller) RestController {
	return RestController{controller}
}

type action func(*Context) error

func (self RestController) Handle(f action) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beginTime := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)
		defer r.Body.Close()

		ctx := NewContext().WithRequest(r).WithResponseWriter(w)

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)

			logger := self.Application.Logger.WithFields(logrus.Fields{
				"duration":    duration,
				"status_code": statusCode,
				"remote":      ctx.RemoteAddress,
			})
			logger.Info(r.Method + " " + r.URL.RequestURI())
		}()

		defer func() {
			if r := recover(); r != nil {
				self.Application.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		w.Header().Set("Content-Type", "application/json")

		httperror := func(w http.ResponseWriter, code int, message string) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(code)
			fmt.Fprint(w, message)
		}

		if err := f(ctx); err != nil {
			switch e := err.(type) {
			case ErrHTTP:
				httperror(w, e.Code(), e.Error())
			default:
				self.Application.Logger.Error(fmt.Sprintf("%+v\n", err))
				httperror(
					w,
					http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
				)
			}
		}
	})
}
