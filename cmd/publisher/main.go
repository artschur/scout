package main

import "go-observability-tool/internal/publisher"

func main() {
	// import from config afterwards
	sender := publisher.NewPublisher("localhost:8082", "my-machine")
	sender.Run()
}
