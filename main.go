package main

import (
	"net/http"
	"rfs/handler/minehandler"
	"rfs/models/entity"
	"rfs/models/modelconst"
)

type MinerHandler struct{}

func (handler *MinerHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

}

func NewMinerHandler() {

}

func main() {

	minehandler := minehandler.NewMinerHandler()

	minehandler.AddNewOperation(entity.NewOperation("first.txt", modelconst.CREATE_FILE, nil))

	minehandler.AddNewOperation(entity.NewOperation("first.txt", modelconst.APPEND_RECORD, []byte("Append please")))

	minehandler.GenerateOpBlock()
	minehandler.GenerateNoOpBlock()

	// // a := time.Now()

	// minerServer := http.Server{
	// 	Addr:    ":8080",
	// 	Handler: &MinerHandler{},
	// }

	// go func() {
	// 	err := minerServer.ListenAndServe()
	// 	log.Println("Started the server on port: ", 8080)

	// 	if err != nil {
	// 		log.Fatalf("Error : ", err)
	// 	}
	// }()

	// interruptChan := make(chan os.Signal, 1)

	// signal.Notify(interruptChan, os.Interrupt)
	// signal.Notify(interruptChan, os.Kill)

	// sig := <-interruptChan
	// log.Println("Got Interrupt: ", sig)

	// ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	// minerServer.Shutdown(ctx)
}
