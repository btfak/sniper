package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/msbranco/goconfig"
)

var config Config

type targetSet []target

type target struct {
	ip    string
	port  string
	path  string
	param string
}

func (t *targetSet) set(a target) {
	*t = append(*t, a)
}

type protocol struct {
	version        string
	cache          bool
	connection     string
	expireSession  bool
	acceptEncoding string
	userAgent      string
}

type header struct {
	head map[string]string
}

type process struct {
	timeout     int
	failures    int
	redirection bool
}

type authenticate struct {
	login string
}

type ssl struct {
	cert    string
	key     string
	timeout int
	ciphers string
}

type proxy struct {
	host  string
	port  string
	login string
}

type command struct {
	concurrent int
	request    int
	time       int
	reps       int

	urlFile     string
	sniperc     string
	postFile    string
	postData    []byte
	contentType string
	plot        string
}

type basic struct {
	https  bool
	method string
	target targetSet
}

type Config struct {
	Basic        basic
	Command      command
	Protocol     protocol
	Header       header
	Process      process
	Authenticate authenticate
	Ssl          ssl
	Proxy        proxy
}

func (c *Config) parseCommandLine() {
	c.prepareFlag(os.Args)
	fs := flag.NewFlagSet("sniper", flag.ExitOnError)
	fs.Usage = sniperUsage
	fs.IntVar(&c.Command.concurrent, "c", 1, "number of multiple requests to make.")
	fs.IntVar(&c.Command.request, "n", 0, "number of requests to perform.")
	fs.IntVar(&c.Command.reps, "r", 0, "number of times to run the test.")
	fs.IntVar(&c.Command.time, "t", 0, "testing time, 30 mean 30 seconds.")
	fs.StringVar(&c.Command.urlFile, "f", "", "select a specific URLS file.")
	fs.StringVar(&c.Command.postFile, "p", "", "POST file name.")
	fs.StringVar(&c.Command.plot, "s", "true", "plot detail transactions' info.")
	fs.StringVar(&c.Command.sniperc, "R", sniperc, "specify an sniperc file.")
	//fs.Var(&c.header, "H", "Add a header to request (can be many).")
	//fs.StringVar(&user_agent, "A", "", "Set User-Agent in request.")
	fs.StringVar(&c.Command.contentType, "T", "", "Set Content-Type in request.")
	fs.Parse(os.Args[1:])
	if c.Command.reps != 0 {
		c.Command.request = c.Command.concurrent * c.Command.reps
	}
	c.checkCommandLine()
	if c.Command.postFile != "" {
		c.Basic.method = "POST"
		c.getPostData()
	} else {
		c.Basic.method = "GET"
	}

	if c.Command.urlFile == "" {
		c.parseUrl(fs)
	} else {
		c.parseUrlFile()
	}
	c.parseConfigFile()
}

func (c *Config) parseConfigFile() error {
	c.setDefaultConfig()
	file, err := goconfig.ReadConfigFile(config.Command.sniperc)
	if err != nil {
		fmt.Println("read config file error:", err, "use default setting.")
		//sniperUsage()
		return nil
	}
	config.Header.head = map[string]string{}
	sections := file.GetSections()
	for _, section := range sections {
		options, _ := file.GetOptions(section)
		for _, option := range options {
			if section == "header" {
				value, _ := file.GetString(section, option)
				config.Header.head[option] = value
			} else {
				switch option {
				case "version":
					value, _ := file.GetString(section, option)
					if value != "HTTP/1.1" && value != "HTTP/1.0" {
						fmt.Println("http version must HTTP/1.1 or HTTP/1.0, check .sniperc file.")
						sniperUsage()
					}
					config.Protocol.version = value
				case "cache":
					value, _ := file.GetBool(section, option)
					config.Protocol.cache = value
				case "connection":
					value, _ := file.GetString(section, option)
					if value != "close" && value != "keep-alive" {
						fmt.Println("connection mode must close or keep-alive, check .sniperc file.")
						sniperUsage()
					}
					config.Protocol.connection = value
				case "expire-session":
					value, _ := file.GetBool(section, option)
					config.Protocol.expireSession = value
				case "accept-encoding":
					value, _ := file.GetString(section, option)
					config.Protocol.acceptEncoding = value
				case "user-agent":
					value, _ := file.GetString(section, option)
					config.Protocol.userAgent = value
				case "timeout":
					value, err := file.GetInt64(section, option)
					if err != nil {
						fmt.Println("set timeout err:", err)
						sniperUsage()
					}
					config.Process.timeout = int(value)
				case "failures":
					value, err := file.GetInt64(section, option)
					if err != nil {
						fmt.Println("set failures err:", err)
						sniperUsage()
					}
					config.Process.failures = int(value)
				case "redirection":
					value, _ := file.GetBool(section, option)
					config.Process.redirection = value
				case "login":
					value, _ := file.GetString(section, option)
					config.Authenticate.login = value
				case "ssl-cert":
					value, _ := file.GetString(section, option)
					config.Ssl.cert = value
				case "ssl-key":
					value, _ := file.GetString(section, option)
					config.Ssl.key = value
				case "ssl-timeout":
					value, _ := file.GetInt64(section, option)
					config.Ssl.timeout = int(value)
				case "ssl-ciphers":
					value, _ := file.GetString(section, option)
					config.Ssl.ciphers = value
				case "proxy-host":
					value, _ := file.GetString(section, option)
					config.Proxy.host = value
				case "proxy-port":
					value, _ := file.GetString(section, option)
					config.Proxy.port = value
				case "proxy-login":
					value, _ := file.GetString(section, option)
					config.Proxy.login = value
				}
			}
		}
	}
	return nil
}

func (c *Config) setDefaultConfig() {
	c.Protocol.version = "HTTP/1.1"
	c.Protocol.cache = false
	c.Protocol.connection = "close"
	c.Process.timeout = 30
	c.Process.failures = 128
	c.Protocol.acceptEncoding = "*"
	c.Protocol.connection = "close"
	c.Protocol.userAgent = "Golang & Sniper"
	c.Command.contentType = "text/plain"
}

func (c *Config) prepareFlag(args []string) {
	for _, v := range args {
		switch v {
		case "-h":
			sniperUsage()
		case "-V":
			sniperVersion()
		case "-C":
			sniperConfig()
		}
	}
}

func (c *Config) checkCommandLine() {
	if c.Command.request == 0 && c.Command.reps == 0 && c.Command.time == 0 {
		sniperUsage()
	}
	if c.Command.sniperc == "" {
		sniperUsage()
	}
	if c.Command.request != 0 && c.Command.request < c.Command.concurrent {
		sniperUsage()
	}
}

func (c *Config) getPostData() {
	bt, err := ioutil.ReadFile(c.Command.postFile)
	if err != nil {
		fmt.Println(err)
		sniperUsage()
	}
	c.Command.postData = bt
}

func (c *Config) parseUrlFile() {
	f, err := os.Open(c.Command.urlFile)
	if err != nil {
		return
	}
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		isURL, _ := regexp.MatchString(`http.*?://.*`, s.Text())
		if !isURL {
			sniperUsage()
		}
		u, err := url.Parse(s.Text())
		if err != nil {
			fmt.Println(err)
			sniperUsage()
		}

		if u.Scheme == "https" {
			c.Basic.https = true
		}

		var ip, port string
		addr := strings.Split(u.Host, ":")
		ip = addr[0]
		if len(addr) == 2 {
			port = addr[1]
		} else if len(addr) == 1 {
			port = "80"
		} else {
			sniperUsage()
		}
		c.Basic.target.set(target{ip, port, u.Path, u.RawQuery})
	}
}

func (c *Config) parseUrl(fs *flag.FlagSet) {
	if fs.NArg() != 1 {
		sniperUsage()
	}
	isURL, _ := regexp.MatchString(`http.*?://.*`, fs.Arg(0))
	if !isURL {
		sniperUsage()
	}

	u, err := url.Parse(fs.Arg(0))
	if err != nil {
		lg(err)
		sniperUsage()
	}

	if u.Scheme == "https" {
		c.Basic.https = true
	}

	var ip, port string
	addr := strings.Split(u.Host, ":")
	ip = addr[0]
	if len(addr) == 2 {
		port = addr[1]
	} else if len(addr) == 1 {
		port = "80"
	} else {
		sniperUsage()
	}
	c.Basic.target.set(target{ip, port, u.Path, u.RawQuery})
}

func initConfig() {
	config.parseCommandLine()
}

func eatSpace(old string) string {
	return strings.Replace(old, " ", "", -1)
}

func lg(l ...interface{}) {
	log.Println(l)
}
