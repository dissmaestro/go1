package main

import (
	"fmt"

	"github.com/dissmaestro/go1/cmd/agent"
	"github.com/dissmaestro/go1/cmd/server"
)

func main() {
	fmt.Println("Запускаем server")
	go func() {
		server.StartServer()
	}()

	fmt.Println("Запускаем агент")
	agent.RunAgent()

	select {} // Блок
}
