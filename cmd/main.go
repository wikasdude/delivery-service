package main

import (
	"delivery-service/handler"
	"delivery-service/utils"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	fmt.Println("inside main .go")
	utils.InitRedis()

	http.HandleFunc("/v1/delivery", handler.Gethandler)

	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8080", nil)
}
