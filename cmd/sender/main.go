package main

import "go-observability-tool/internal/publisher"

func main() {
	// import from config afterwards
	sender := publisher.NewSender("localhost:8082")
	sender.Run()
}
