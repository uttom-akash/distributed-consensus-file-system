package cfslib

import (
	"log"
	"net/http"
)

func HttpListen(serverConf http.Server) {

	peerServer := serverConf

	// go func() {
	log.Println("bclib/HttpListen - starting the peer server on: ", peerServer.Addr)

	httpErr := peerServer.ListenAndServe()

	log.Println("bclib/HttpListen - started the peer server on: ", peerServer.Addr)

	if httpErr != nil {
		log.Fatalln("bclib/HttpListen - error : ", httpErr)
	}

	// interruptChan := make(chan os.Signal, 1)

	// signal.Notify(interruptChan, os.Interrupt)
	// signal.Notify(interruptChan, os.Kill)

	// sig := <-interruptChan
	// log.Println("Got Interrupt: ", sig)

	// ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	// peerServer.Shutdown(ctx)

	// }()
}
