package run

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"web-app-go/internal/config"
	"web-app-go/internal/loggers/apachelogger"
	"web-app-go/internal/router"
	"web-app-go/internal/utils"
	"web-app-go/pkg/exit"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "run",
	Short: "Run server.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		//fmt.Println("RUN.....")
		log.Println("[INFO    ] Loaded Configuration:\n", utils.Sdump(config.GetConfig()))

		svc := &http.Server{
			Addr:    config.GetConfig().GetListenConfig(),
			Handler: apachelogger.NewApacheLoggingHandler(os.Stdout, router.NewRouter()),
		}
		go func() {
			if err := svc.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("[CRITICAL] Unable to start HTTP Service: ", err.Error())
			}
		}()
		log.Printf("[INFO    ] HTTP Service strated.")

		// Register pre-quit handling
		exit.Listen(func(s os.Signal) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer func() { cancel() }()

			if err := svc.Shutdown(ctx); err != nil {
				log.Fatalf("[CRITICAL] Server Shutdown Failed:%+v", err)
			}
			log.Println("[INFO    ] Exiting.....")
		})

		// Wait to quit.
		exit.Wait()
		return nil
	},
}
