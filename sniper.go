package main

import (
	"runtime"
)

func main() {
	initConfig()
	prepBattle()
	if err := taste(); err != nil {
		lg(err)
		return
	}
	go manager.Monitor()
	go Snipe()

	<-over
}

func prepBattle() {
	over = make(chan struct{})
	//runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(1)
	record.trans = make(chan *transaction, config.Command.concurrent)
	message.prepareContent()
}