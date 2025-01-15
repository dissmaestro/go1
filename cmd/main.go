package main

import (
	"fmt"
	"~/Documents/projects/golang/1/cmd/server"

	"github.com/dissmaestro/go1/cmd/agent"
)

func main() {
	fmt.Println("Запускаем server")
	go func() {
		server.Main()
	}()

	fmt.Println("Запускаем агент")
	agent.RunAgent()

	select {} // Блок
}
