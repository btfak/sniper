package main

import (
	"syscall"
	"time"
)

/*
	c.sa = new(syscall.SockaddrInet4)
	c.sa.Addr = [4]byte{10, 1, 1, 131}
	c.sa.Port = port
*/
func (c *Worker) socket() {
	begin := time.Now()
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		syscall.Close(s)
		return
	}

	err = syscall.Connect(s, nil)//TODO
	if err != nil {
		lg(err)
		syscall.Close(s)
		return
	}

	_, err = syscall.Write(s, message.content)
	if err != nil {
		syscall.Close(s)
		return
	}

	_, err = syscall.Read(s, c.data)
	if err != nil {
		syscall.Close(s)
		return
	}

	syscall.Close(s)
	time.Now().Sub(begin)
}
