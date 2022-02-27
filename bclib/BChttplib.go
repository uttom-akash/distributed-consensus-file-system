package bclib

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func HttpListen(serverConf http.Server) {

	peerServer := serverConf

	go func() {
		err := peerServer.ListenAndServe()
		log.Println("Started the peer server on port: ", 8080)

		if err != nil {
			log.Fatalf("Error : ", err)
		}
	}()

	interruptChan := make(chan os.Signal, 1)

	signal.Notify(interruptChan, os.Interrupt)
	signal.Notify(interruptChan, os.Kill)

	sig := <-interruptChan
	log.Println("Got Interrupt: ", sig)

	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	peerServer.Shutdown(ctx)
}
