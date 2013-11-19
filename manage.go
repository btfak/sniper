package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"time"
)

var manager Manager
var temSuccess int

type httpLen struct {
	total int
	head  int
	body  int
}

type Manager struct {
	Transactions    int
	Availability    float32
	ElapsedTime     time.Duration
	TotalTransfer   float64
	HTMLTransfer    float32
	TransactionRate float32
	Throughput      float64
	Successful      int
	Failed          int
	StateCode       int

	TransactionTime time.Duration
	ConnectionTime  time.Duration
	ResponseTime    time.Duration

	TotalTransTime time.Duration
	TotalConnTime  time.Duration
	TotalRespTime  time.Duration

	HttpLen httpLen
}

var transactions Transactions

func (m *Manager) Monitor() {
	userInterrupt := make(chan os.Signal, 1)
	signal.Notify(userInterrupt, os.Interrupt)

	temSuccess = config.Command.request / 10
	timer := new(time.Timer)
	if config.Command.time > 0 {
		timer = time.NewTimer(time.Duration(config.Command.time) * time.Second)
	}

	var identy int
	var mark int = config.Command.request/1000 + 1
	printTitle()
	begin := time.Now()
loop:
	for {
		select {
		case trans := <-record.trans:
			m.Transactions++
			if trans.totalTime != 0 {
				identy++
				m.Successful++
				if trans.code == 200 {
					m.StateCode += 1
				}
				if identy == mark && config.Command.time == 0 {
					identy = 0
					m.AppendTrans(trans)
					printCompleteReq()
				} else if config.Command.time != 0 {
					m.AddRecord(trans)
				}
				if m.Successful == config.Command.request {
					break loop
				}
			} else {
				m.Failed++
				if m.Failed >= config.Process.failures {
					break loop
				}
			}
		case <-userInterrupt:
			break loop
		case <-timer.C:
			break loop
		}
	}
	m.ElapsedTime = time.Now().Sub(begin)
	m.Summary()
	m.Dump()
	m.Plot()
	close(over)
}

func (m *Manager) AppendTrans(trans *transaction) {
	temp := *trans
	transactions = append(transactions, &temp)
}

func (m *Manager) AddRecord(trans *transaction) {
	m.TotalTransTime += trans.totalTime
	m.TotalConnTime += trans.connectionTime
	m.TotalRespTime += trans.responseTime
}

func (m *Manager) Summary() {
	if m.Transactions != 0 {
		m.Availability = float32(m.Successful) / float32(m.Transactions)
	} else {
		m.Availability = 0
	}
	m.TransactionRate = float32(m.Successful) / float32(m.ElapsedTime.Seconds()) //TODO
	if config.Protocol.connection == "close" {
		m.TotalTransfer = float64(m.HttpLen.total) * float64(m.Successful) / float64(1024*1024)
	} else {
		m.TotalTransfer = (float64(m.HttpLen.head) + float64(m.HttpLen.body)*float64(m.Successful)) / float64(1024*1024)
	}
	m.Throughput = m.TotalTransfer / float64(m.ElapsedTime.Seconds())
	if config.Command.time == 0 {
		var total, conn, res time.Duration
		for _, trans := range transactions {
			total += trans.totalTime
			conn += trans.connectionTime
			res += trans.responseTime
		}
		if size := time.Duration(len(transactions)); size != 0 {
			m.TransactionTime = total / size
			m.ConnectionTime = conn / size
			m.ResponseTime = res / size
		}
	} else if size := time.Duration(m.Transactions); size != 0 {
		m.TransactionTime = m.TotalTransTime / size
		m.ConnectionTime = m.TotalConnTime / size
		m.ResponseTime = m.TotalRespTime / size
	}
}

func (m *Manager) Dump() {
	fmt.Println(fmt.Sprintf("%s%d%s", "Transactions:                   ", m.Transactions, " hits"))
	fmt.Println(fmt.Sprintf("%s%5.2f%s", "Availability:                   ", m.Availability*100, " %"))
	fmt.Println(fmt.Sprintf("%s%7.2f%s", "Elapsed time:                ", m.ElapsedTime.Seconds(), " secs"))
	fmt.Println(fmt.Sprintf("%s%d%s", "Document length:               ", maxRecvSize, " Bytes"))
	fmt.Println(fmt.Sprintf("%s%8.2f%s", "TotalTransfer:              ", m.TotalTransfer, " MB"))
	//fmt.Println(fmt.Sprintf("%s%8.2f%s", "HTMLTransfer:               ", m.HTMLTransfer, " MB"))
	fmt.Println(fmt.Sprintf("%s%7.2f%s", "Transaction rate:            ", m.TransactionRate, " trans/sec"))
	fmt.Println(fmt.Sprintf("%s%5.2f%s", "Throughput:                    ", m.Throughput, " MB/sec"))
	fmt.Println(fmt.Sprintf("%s%d%s", "Successful:                     ", m.Successful, " hits"))
	fmt.Println(fmt.Sprintf("%s%d%s", "Failed:                           ", m.Failed, " hits"))
	fmt.Println(fmt.Sprintf("%s%8.3f%s", "TransactionTime:            ", m.TransactionTime.Seconds()*1000, " ms(mean)"))
	fmt.Println(fmt.Sprintf("%s%8.3f%s", "ConnectionTime:             ", m.ConnectionTime.Seconds()*1000, " ms(mean)"))
	fmt.Println(fmt.Sprintf("%s%8.3f%s", "ProcessTime:                ", m.ResponseTime.Seconds()*1000, " ms(mean)"))
	fmt.Println(fmt.Sprintf("%s%d%s", "StateCode:                    ", m.StateCode, "(code 200)"))
}

func (m *Manager) Plot() {
	if config.Command.plot != "true" || config.Command.time != 0 {
		return
	}
	sort.Sort(transactions)
	out := &bytes.Buffer{}
	for _, result := range transactions {
		fmt.Fprintf(out, "[%.3f,%f,%f,%f],",
			result.currentTime.Sub(transactions[0].currentTime).Seconds(),
			result.totalTime.Seconds()*1000,
			result.connectionTime.Seconds()*1000,
			result.responseTime.Seconds()*1000)
	}
	if out.Len() > 0 {
		out.Truncate(out.Len() - 1) // Remove trailing comma
	}
	buffer := fmt.Sprintf(plotsTemplate, dygraphJSLibSrc(), out)
	perm := os.FileMode(0666)
	f, err := os.OpenFile("plot.html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write([]byte(buffer))
	if err != nil {
		return
	}
}

var printTitle = func() {
	fmt.Println(sniperVersionInfo + "\r\n")
	fmt.Println("The server is now under snipe ...\r\n")
}

var printCompleteReq = func() {
	if manager.Successful >= temSuccess {
		if config.Command.request > 10 {
			fmt.Println(fmt.Sprintf("%s %d %s", "Completed", temSuccess, "requests"))
		}
		temSuccess += config.Command.request / 10
	}
}
