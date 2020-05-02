package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	h "github.com/rinosukmandityo/maknews/api"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go run main.go

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	set url=root:Password.1@tcp(127.0.0.1:3306)/news
	set timeout=10
	set db=news
	set driver=mysql
	set redis_url=redis://:@localhost:6379/0
	set redis_timeout=10
	set elastic_url=http://localhost:9200
	set elastic_timeout=10
	set elastic_index=news
	set kafka_url=localhost:9092
	set kafka_timeout=10
	set kafka_topic=news
*/

func main() {
	r := h.RegisterHandler()

	errs := make(chan error, 2)
	go func() {
		log.Printf("Listening on port %s\n", httpPort())
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)

	}()
	log.Printf("Terminated %s", <-errs)

}

func httpPort() string {
	port := "8000"
	if os.Getenv("port") != "" {
		port = os.Getenv("port")
	}
	return fmt.Sprintf(":%s", port)
}
