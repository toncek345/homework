package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/client"
	"github.com/spacelift-io/homework-object-storage/api"
	"github.com/spacelift-io/homework-object-storage/storage"
)

func main() {
	dockerCli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("getting docker client: %s", err)
	}

	storagePool, err := storage.NewMinioPool(dockerCli, "objects")
	if err != nil {
		log.Fatalf("creating storage pool: %s", err)
	}

	api := api.NewApi(storagePool)
	log.Println("starting api....")
	api.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT)

	<-signalChan
	log.Println("shutting down...")
	api.Stop()
}
