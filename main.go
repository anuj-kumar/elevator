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
	var liftCount uint = 1
	lifts := make([]lift.ILift, liftCount)
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	var i uint
	for i = range liftCount {
		lifts[i] = lift.NewLiftImpl(i+1, 100)
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
