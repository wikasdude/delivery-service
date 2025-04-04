package main

import (
	"delivery-service/handler"
	"delivery-service/utils"
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("inside main .go")
	utils.InitRedis()

	http.HandleFunc("/v1/delivery", handler.Gethandler)
	http.ListenAndServe(":8080", nil)
}
