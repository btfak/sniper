package main

import (
	"fmt"
	"runtime"
)

func main() {
	initConfig()
	prepare()
	if err := taste(); err != nil {
		fmt.Println(err)
		return
	}

	go manager.Monitor()
	go Snipe()

	<-over
}

func prepare() {
	over = make(chan struct{})
	//TODO:multicore but low performance
	runtime.GOMAXPROCS(1) //runtime.NumCPU()
	record.trans = make(chan *transaction, config.Command.concurrent)
	message.prepareContent()
}
