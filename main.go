package main

import (
	"context"
	"github.com/neda1985/getir/pkg/apis"
	"github.com/neda1985/getir/pkg/dataManager/cache"
	mongodb "github.com/neda1985/getir/pkg/dataManager/mongo"
	"github.com/neda1985/getir/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	logger.Start()
	m := apis.MongoMangerService{MongoManager: mongodb.NewMongoInstance()}
	in := apis.InMemoryManagerService{InMemory: cache.NewInMemoryService()}
	http.HandleFunc("/fetch", m.FetchData)
	http.HandleFunc("/in-memory", in.SwitchRout)
	logger.LogInfo("started listening on port 8080")
	logger.LogInfo("Waiting request...")
	httpServer := &http.Server{
		Addr: ":8080",
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	ctx, signalCancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer signalCancel()

	<-ctx.Done()

	log.Print("os.Interrupt - shutting down...\n")

	Context, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(Context); err != nil {
		log.Printf("shutdown error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Context stopped\n")

}
