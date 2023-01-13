package ddns

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	CronExpression string
	Service        Service
}

func (scheduler Scheduler) Start() {

	c := cron.New()
	job := cron.FuncJob(func() {
		scheduler.Service.Run()
	})

	c.AddJob(scheduler.CronExpression, job)
	c.Start()

	fmt.Println("Started DynDNS")

	// Do the first run instantly instead of waiting for cron
	scheduler.Service.Run()

	wait()
	fmt.Println("Stopped DynDNS")
	c.Stop()
}

func wait() {
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
