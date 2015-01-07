package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
)

func Snipe() {
	for i := 0; i < config.Command.concurrent; i++ {
		go run()
	}
}

func run() {
	c := &Worker{}

	//record
	c.trans = &transaction{}

	//receive
	if maxRecvSize < 1024 {
		c.data = make([]byte, maxRecvSize)
	} else {
		c.data = make([]byte, 1024)
	}

	//https tls
	if config.Basic.https {
		cert, err := tls.LoadX509KeyPair(config.Ssl.cert, config.Ssl.key)
		if err != nil {
			fmt.Println("open ssl key err:", err)
			sniperUsage()
		}
		c.tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	} else {
		//http
		ip := net.ParseIP(config.Basic.target[0].ip)
		port, _ := strconv.Atoi(config.Basic.target[0].port)
		c.tcpAddr = &net.TCPAddr{ip, port, ""}
	}
	if config.Protocol.connection == "close" {
		for {
			c.fire()
		}
	} else {
		//keep-alive
		c.shoot()
	}
}

type Transactions []*transaction

type transaction struct {
	currentTime    time.Time
	totalTime      time.Duration
	connectionTime time.Duration
	responseTime   time.Duration
	code           int
}

func (t Transactions) Len() int           { return len(t) }
func (t Transactions) Less(i, j int) bool { return t[i].currentTime.Before(t[j].currentTime) }
func (t Transactions) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type Worker struct {
	data      []byte
	tcpAddr   *net.TCPAddr
	conn      net.Conn
	trans     *transaction
	tlsConfig *tls.Config
}

func (c *Worker) fire() {
	var err error
	var deadline int
	c.trans.currentTime = time.Now()

	//http & https
	if config.Basic.https {
		c.conn, err = tls.Dial("tcp", config.Basic.target[0].ip+":"+config.Basic.target[0].port, c.tlsConfig)
		deadline = config.Ssl.timeout
	} else {
		c.conn, err = net.DialTCP("tcp", nil, c.tcpAddr)
		deadline = config.Process.timeout
	}
	if err != nil {
		c.trans.totalTime = 0
		record.trans <- c.trans
		return
	}
	c.conn.SetDeadline(c.trans.currentTime.Add(time.Second * time.Duration(deadline)))

	connection := time.Now()
	c.trans.connectionTime = connection.Sub(c.trans.currentTime)
	_, err = c.conn.Write(message.content)
	if err != nil {
		c.trans.totalTime = 0
		record.trans <- c.trans
		c.conn.Close()
		return
	}
	sum := 0
	var isFirstLine = true
	for {
		n, err := c.conn.Read(c.data)
		if err != nil {
			if err == io.EOF {
				break
			}
			c.trans.totalTime = 0
			record.trans <- c.trans
			c.conn.Close()
			return
		}
		if isFirstLine && len(c.data) > 12 {
			if c.data[9] == 50 && c.data[10] == 48 && c.data[11] == 48 {
				c.trans.code = 200
                isFirstLine = false
			}
		}
		sum += n
		if sum == maxRecvSize {
			break
		}
	}

	response := time.Now()
	c.trans.responseTime = response.Sub(connection)
	c.trans.totalTime = response.Sub(c.trans.currentTime)
	record.trans <- c.trans
	c.conn.Close()
	return
}

func (c *Worker) shoot() {
	c.connect()
	for {
		success := c.request()
		if !success {
			c.conn.Close()
			c.connect()
		}
	}
}

func (c *Worker) connect() {
	var err error
	var deadline int
	ip := config.Basic.target[0].ip
	port := config.Basic.target[0].port

	begin := time.Now()
	c.trans.currentTime = begin

	//http & https
	if config.Basic.https {
		c.conn, err = tls.Dial("tcp", ip+":"+port, c.tlsConfig)
		deadline = config.Ssl.timeout
	} else {
		c.conn, err = net.DialTCP("tcp", nil, c.tcpAddr)
		deadline = config.Process.timeout
	}

	if err != nil {
		c.trans.totalTime = 0
		record.trans <- c.trans
		c.connect()
	}
	c.conn.SetDeadline(c.trans.currentTime.Add(time.Second * time.Duration(deadline)))

	connection := time.Now()
	c.trans.connectionTime = connection.Sub(begin)
}

func (c *Worker) request() bool {
	connection := time.Now()
	c.trans.currentTime = connection
	_, err := c.conn.Write(message.content)
	if err != nil {
		c.trans.totalTime = 0
		record.trans <- c.trans
		return true
	}

	sum := 0
	var isFirstLine = true
	for {
		n, err := c.conn.Read(c.data)
		if err != nil {
			// When server close the connection or server reset the connection
			// We don't reguard it as an request error, just reconnect server.
			if err == io.EOF {
				return false
			} else if opErr, ok := err.(*net.OpError); ok {
				if opErr.Err.Error() == "connection reset by peer" {
					return false
				}
			}

			// Other situation we reguard it as request error.
			c.trans.totalTime = 0
			record.trans <- c.trans
			return true
		}
		if isFirstLine && len(c.data) > 12 {
			if c.data[9] == 50 && c.data[10] == 48 && c.data[11] == 48 {
				c.trans.code = 200
			}
			isFirstLine = false
		}
		sum += n
		if sum == maxRecvSize {
			break
		}
	}
	response := time.Now()
	c.trans.responseTime = response.Sub(connection)
	c.trans.totalTime = response.Sub(connection)
	record.trans <- c.trans
	c.trans.connectionTime = 0
	return true
}

func taste() error {
	var conn net.Conn
	var err error
	ip := config.Basic.target[0].ip
	port := config.Basic.target[0].port

	if config.Basic.https {
		cert, err := tls.LoadX509KeyPair(config.Ssl.cert, config.Ssl.key)
		if err != nil {
			fmt.Println("open ssl key err:", err)
			return err
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		conn, err = tls.Dial("tcp", ip+":"+port, tlsConfig)
		if err != nil {
			return err
		}
	} else {
		conn, err = net.Dial("tcp", ip+":"+port)
		if err != nil {
			return err
		}
	}

	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	addr := strings.Split(remoteAddr, ":")
	config.Basic.target[0].ip = addr[0]

	_, err = conn.Write(message.content)
	if err != nil {
		return err
	}
	if config.Protocol.connection == "close" {
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		maxRecvSize = len(b)
		manager.HttpLen.total = len(b)
	} else {
		var chunked bool
		var content_len int
		var head_len int

		buffer := bufio.NewReader(conn)
		var hasContentLength bool
		for {
			line, _ := buffer.ReadString('\n')
			head_len += len(line)
			//end of header
			if len(line) <= 2 {
				c := strings.Replace(line, "\n", "", -1)
				c = strings.Replace(c, "\r", "", -1)
				if c == "" {
					if hasContentLength == false && chunked == false {
						return errors.New("http stack err")
					}
					break
				}
			}

			//Content-Length
			r := strings.Split(line, ":")
			if r[0] == "Content-Length" {
				content_len, err = strconv.Atoi(eatCRLF(r[1]))
				if err != nil {
					return err
				}
				hasContentLength = true
			}

			//Transfer-Encoding
			if r[0] == "Transfer-Encoding" {
				c := eatCRLF(r[1])
				if c == "chunked" {
					chunked = true
				}
			}
		}

		if chunked {
			for {
				//chunk-size
				line, _ := buffer.ReadString('\n')
				content_len += len(line)

				//chunk-ext
				tempLine := strings.Split(line, ";")
				line = tempLine[0]

				//CRLF
				line = strings.Replace(line, "\n", "", -1)
				line = strings.Replace(line, "\r", "", -1)
				chunk_len, err := strconv.ParseInt(line, 16, 64)
				if err != nil {
					return errors.New("http stack err")
				}

				if chunk_len == 0 {
					content_len += buffer.Buffered()
					break
				}

				var b byte
				var count int64
				for {
					b, err = buffer.ReadByte()
					if err != nil {
						return errors.New("http stack err")
					}
					count++
					content_len += 1
					if count == chunk_len+1 && b == '\n' {
						break
					}
					if count == chunk_len+2 {
						break
					}
				}
			}
		}
		maxRecvSize = head_len + content_len
		manager.HttpLen = httpLen{maxRecvSize, head_len, content_len}
	}
	return nil
}

func eatCRLF(str string) string {
	c := strings.Replace(str, " ", "", -1)
	c = strings.Replace(c, "\n", "", -1)
	c = strings.Replace(c, "\r", "", -1)
	return c
}
