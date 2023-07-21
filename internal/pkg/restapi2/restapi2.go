package restapi2

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/trace"
	"sync"
)

func Run(addr, wsAddr string) error {

	f, _ := os.Create("trace.out")
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()

	router, srv := setupRouter(wsAddr)

	srv.Run()
	defer srv.Wait()
	defer srv.Stop()
	sr := http.Server{Addr: addr, Handler: router}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//http server runs in goroutine, main thread is responsible for os signals handling
		sr.ListenAndServe()
	}()

	go func() {
		http.ListenAndServe("localhost:6060", nil)

	}()

	// Setting up signal capturing
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stopSignal
	return sr.Shutdown(context.Background())
}
