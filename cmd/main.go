package main

import "todo-service/internal/api"

func main() {
	api := api.InitAPI()
	api.StartServer()

}
