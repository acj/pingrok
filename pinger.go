package main

import (
	"log"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const protocolICMPIPv4 = 1

type pinger struct {
	connection       *icmp.PacketConn
	quit             chan int
	dataPoints       chan latencyDataPoint
	messagesInFlight *pendingEchos
	startTime        time.Time
	targetHost       string
}

func newPinger(targetHost string, dataPoints chan latencyDataPoint) *pinger {
	return &pinger{
		targetHost: targetHost,
		dataPoints: dataPoints,
	}
}

func (p *pinger) start() {
	p.quit = make(chan int)
	p.messagesInFlight = newPendingEchoes()
	p.startTime = time.Now()

	var err error
	p.connection, err = icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	go p.consumer()
	go p.producer(p.targetHost, 10*time.Millisecond)
}

func (p *pinger) stop() {
	close(p.quit)
	p.connection.Close()
}

func (p *pinger) producer(destinationIP string, interval time.Duration) {
	body := &icmp.Echo{
		ID:   os.Getpid() & 0xffff,
		Seq:  0,
		Data: []byte("Now is the time for all good homo sapiens to come to the aid of their species"),
	}
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: body,
	}

	for {
		wb, err := msg.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}

		p.messagesInFlight.start(body.Seq)
		if _, err := p.connection.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(destinationIP)}); err != nil {
			log.Printf("error sending echo request: %v", err)
		}

		body.Seq++
		time.Sleep(interval)
	}
}

func (p *pinger) consumer() {
	rb := make([]byte, 1500)

	for {
		n, peer, err := p.connection.ReadFrom(rb)
		if err != nil {
			log.Fatal(err)
		}

		candidateReceiptTime := time.Now()
		rm, err := icmp.ParseMessage(protocolICMPIPv4, rb[:n])
		if err != nil {
			log.Fatal(err)
		}
		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			echoReply := rm.Body.(*icmp.Echo)

			echoRequestSentTime, ok := p.messagesInFlight.resolve(echoReply.Seq)
			if !ok {
				log.Printf("unexpected message #%d, sent at %v", echoReply.Seq, echoRequestSentTime)
				continue
			}

			timeOffset := echoRequestSentTime.Sub(p.startTime).Seconds()
			latency := float64(candidateReceiptTime.Sub(echoRequestSentTime).Nanoseconds()) / 1e6
			p.dataPoints <- latencyDataPoint{timeOffset: timeOffset, latency: latency}
		default:
			log.Printf("unexpected message from %v: got %+v, want echo reply", peer, rm)
		}
	}
}

type pendingEchos struct {
	mux   sync.Mutex
	times map[int]time.Time
}

func newPendingEchoes() *pendingEchos {
	return &pendingEchos{
		times: make(map[int]time.Time),
	}
}

func (mt *pendingEchos) start(sequenceNumber int) {
	mt.mux.Lock()
	mt.times[sequenceNumber] = time.Now()
	mt.mux.Unlock()
}

func (mt *pendingEchos) resolve(sequenceNumber int) (time.Time, bool) {
	mt.mux.Lock()
	defer mt.mux.Unlock()

	time, ok := mt.times[sequenceNumber]
	delete(mt.times, sequenceNumber)

	return time, ok
}
