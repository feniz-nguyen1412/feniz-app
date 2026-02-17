package main

import (
	"gomono_template/internal/bootstrap"
	"log"
	"net/http"
	"sync"
)

func main() {
	container, err := bootstrap.BuildContainer()
	if err != nil {
		panic(err)
	}
	defer container.RedisClient.Close()
	defer container.KafkaProducer.Close()
	defer container.KafkaConsumer.Close()

	// Start HTTP server
	httpServer := bootstrap.NewHTTPServer(container)
	
	// Start gRPC server
	grpcServer, grpcListener := bootstrap.NewGRPCServer(container)

	var wg sync.WaitGroup

	// HTTP Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("HTTP server starting on port %s", container.Config.AppPort)
		if err := httpServer.Run(":" + container.Config.AppPort); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// gRPC Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("gRPC server starting on port %s", container.Config.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	wg.Wait()
}
