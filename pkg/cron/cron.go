package cron

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/robfig/cron/v3"
)

func StartCronScheduler(cronExpression string, cronJob cron.FuncJob) {
	waitUntilStopped := func() {
		var endWaiter sync.WaitGroup
		endWaiter.Add(1)
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		go func() {
			<-signalChannel
			endWaiter.Done()
		}()
		endWaiter.Wait()
	}
	cron := cron.New()
	cron.AddJob(cronExpression, cronJob)
	cron.Start()
	fmt.Println("Started DynDNS")
	waitUntilStopped()
	fmt.Println("Stopped DynDNS")
	cron.Stop()
}
