package main

import (
	agent_app "lms-1/internal/agent/app"
	orchestrator_app "lms-1/internal/orchestrator/app"
	"time"
)

func main() {
	go func() {
		orchestrator_app.Run()
	}()
	time.Sleep(1 * time.Second)
	go func() {
		agent_app.Run()
	}()
	select {} 
}
