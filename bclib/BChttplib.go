package bclib

import (
	"log"
	"net/http"
)

func HttpListen(serverConf http.Server) {

	peerServer := serverConf

	// go func() {
	log.Println("starting the peer server on: ", peerServer.Addr)
	err := peerServer.ListenAndServe()
	log.Println("Started the peer server on: ", peerServer.Addr)

	if err != nil {
		log.Fatalf("Error : ", err)
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
