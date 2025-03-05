package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"anujkumar.com/elevator/lift"
	"anujkumar.com/elevator/models"
	"anujkumar.com/elevator/server"
)

func main() {
	lifts := make([]lift.ILift, 4)
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	var i uint
	for i = range 4 {
		lifts[i] = lift.NewModel(i+1, 100)
		go lifts[i].Operate((timeoutCtx))
	}

	assigner := server.NewService(lifts)
	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-stopChan:
			return
		default:
			var source, dest models.Floor
			fmt.Println("Enter source and destination floors:")
			fmt.Scanf("%d %d", &source, &dest)
			assigner.Assign(ctx, models.Request{
				Source: source,
				Dest:   dest,
			})
		}

	}
}
