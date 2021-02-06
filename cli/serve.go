package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"

	"microtecture/infrastructure/config"
	"microtecture/interface/router"
	"microtecture/registry"
	"microtecture/usecase/controllers"
)

func serveAPI(ctx context.Context, controller controllers.Root) {
	r := httprouter.New()

	// cors config (github.com/rs/cors)
	//c := cors.New(cors.Options{
	//	AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "PATCH", "DELETE"},
	//	AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "X-Requested-With", "Authorization"},
	//	AllowCredentials: true,
	//	Debug:            true,
	//	MaxAge:           12 * 60 * 60, // 12 hour
	//})

	router.Route(r, controller)

	// set cors to handler
	//handler := c.Handler(r)

	s := &http.Server{
		Addr: fmt.Sprintf("%s%d", ":", controller.GetBase().Application.Config.Port),
		// use cors in server
		//Handler:           handler,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Second,
	}
	http2.ConfigureServer(s, nil)

	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			logrus.Error(err)
		}
	}()

	logrus.Infof("Serving at http://localhost:%d", controller.GetBase().Application.Config.Port)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Panic(err)
	}
}

var serveCli = &cobra.Command{
	Use:   "serve",
	Short: "Serve the application.",
	RunE: func(cli *cobra.Command, args []string) error {
		if os.Getenv(config.ENVIRONMENT_NAME) == "dev" {
			logrus.Warning(
				`Running in "dev" mode. use `,
				config.ENVIRONMENT_NAME,
				` in "dev" or "prod" for development or production mode.`,
			)
		}

		reg, err := registry.New()
		if err != nil {
			errMsg := fmt.Sprintf("%+v\n", err)
			return fmt.Errorf(errMsg)
		}
		c := reg.NewRootController()

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(
				sigint,
				os.Interrupt,
				syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT,
			)
			<-sigint
			logrus.Info("Signal caught. Shutting down...")
			cancel()
		}()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			serveAPI(ctx, c)
		}()

		wg.Wait()
		return nil
	},
}

func init() {
	rootCli.AddCommand(serveCli)
}
