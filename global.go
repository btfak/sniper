package main

import (
	"fmt"
	"os"
)

//work variant
var over chan struct{}
var home = os.Getenv("HOME")

var maxRecvSize int

//var maxReqSize int

const (
	recvSizeOffset = 64
)

//config variant
var sniperLog = "/usr/local/var/sniper.log"
var sniperc = home + "/.sniperc"

var record Record

type Record struct {
	trans chan *transaction
}

//function variant
var sniperUsage = func() {
	fmt.Println(sniperUsageInfo)
	os.Exit(2)
}

var sniperVersion = func() {
	fmt.Println(sniperVersionInfo)
	os.Exit(2)
}

var sniperConfig = func() {
	config.Command.sniperc = sniperc
	config.parseConfigFile()
	sniperShowConfig()
	os.Exit(2)
}

var sniperShowConfig = func() {
	fmt.Println("CURRENT  SNIPER  CONFIGURATION")
	fmt.Println("Edit the config file to change the settings")
	fmt.Println(`--------------------------------------------`)
	fmt.Printf("version:                          %s\n", config.Protocol.version)
	fmt.Printf("connection:                       %s\n", config.Protocol.connection)
	fmt.Printf("accept-encoding:                  %s\n", config.Protocol.acceptEncoding)
	fmt.Printf("user-agent:                       %s\n", config.Protocol.userAgent)
	fmt.Printf("timeout:                          %d\n", config.Process.timeout)
	fmt.Printf("failures:                         %d\n", config.Process.failures)
	fmt.Printf("login:                            %s\n", config.Authenticate.login)
	fmt.Printf("ssl-cert:                         %s\n", config.Ssl.cert)
	fmt.Printf("ssl-key:                          %s\n", config.Ssl.key)
	fmt.Printf("ssl-timeout:                      %d\n", config.Ssl.timeout)
	//fmt.Printf("ssl-ciphers:                      %s\n", config.Ssl.ciphers)
	//fmt.Printf("proxy-host:                       %s\n", config.Proxy.host)
	//fmt.Printf("proxy-port:                       %s\n", config.Proxy.port)
	//fmt.Printf("proxy-login:                      %s\n", config.Proxy.login)
	fmt.Println(`--------------------------------------------`)
}

var sniperNoData = func() {
	fmt.Println(sniperNoDataInfo)
	os.Exit(2)
}

var sniperOpenFileError = func() {
	fmt.Println(sniperOpenFileErrorInfo)
	os.Exit(2)
}
