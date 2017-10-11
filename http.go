package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

var message Message

type Message struct {
	content []byte
}

func (m *Message) prepareContent() {
	m.buildHeader()
	if config.Basic.method == "POST" {
		m.buildBody()
	}
}

func (m *Message) buildHeader() {
	var msg, path string
	//first line
	path = config.Basic.target[0].path
	if path == "" {
		path = "/"
	}
	line := fmt.Sprintf("%s%s%s", config.Basic.method+" ", path+"?"+config.Basic.target[0].param+" ", config.Protocol.version+"\r\n")
	msg = msg + line

	//host
	line = "Host: " + config.Basic.target[0].ip + "\r\n"
	msg = msg + line

	//accept
	line = "Accept: */*" + "\r\n"
	msg = msg + line

	//Accept-Encoding
	line = "Accept-Encoding: " + config.Protocol.acceptEncoding + "\r\n"
	msg = msg + line

	//User-Agent
	line = "User-Agent: " + config.Protocol.userAgent + "\r\n"
	msg = msg + line

	//user-header
	for k, v := range config.Header.head {
		msg = msg + k + ": " + v + "\r\n"
	}

	//post
	if config.Basic.method == "POST" {
		//content-length
		line = "Content-Length:" + strconv.Itoa(len(config.Command.postData)) + "\r\n"
		msg = msg + line

		//content-type
		line = "Content-Type:" + config.Command.contentType + "\r\n"
		msg = msg + line
	}

	//Authenticate
	if login := config.Authenticate.login; login != "" {
		line = "Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte(login)) + "\r\n"
		msg = msg + line
	}

	//Connection
	line = "Connection: " + config.Protocol.connection + "\r\n"
	msg = msg + line

	//blank line
	line = "\r\n"
	msg = msg + line

	m.content = []byte(msg)
}

func (m *Message) buildBody() {
	m.content = append(m.content, config.Command.postData...)
	if config.Protocol.connection == "close" {
		//end of post data
		//m.content = append(m.content, []byte("\r\n")...)

		//blank line
		//m.content = append(m.content, []byte("\r\n")...)
	}
}
