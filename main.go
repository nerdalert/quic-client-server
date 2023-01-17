package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// https://www.rfc-editor.org/rfc/rfc9221.html
func main() {

	if len(os.Args) != 3 {
		fmt.Printf("Usage:\n")
		fmt.Printf("mode    - \"server\" or \"client\"\n")
		fmt.Printf("ip:port - local or remote depending on the mode\n\n")
		fmt.Printf("Example from two different hosts:\n")
		fmt.Printf("server: quic-udp-linux server 192.168.1.201:34000\n")
		fmt.Printf("client: quic-udp-linux client 192.168.1.201:34000\n")
		os.Exit(1)
	}
	mode := os.Args[1]
	socket := os.Args[2]
	_, _, err := net.SplitHostPort(socket)
	if err != nil {
		log.Fatalf("unable to split the IP:Port argument %s: %v", socket, err)
	}
	// server mode
	if mode == "server" {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		s := NewServer(socket)
		s.SetHandler(func(c Ctx) error {
			msg := c.String()
			log.Printf("Client [ %s ] sent a message [ %s ]", c.RemoteAddr().String(), msg)
			return nil
		})
		log.Fatal(s.StartServer(ctx))
	}
	// client mode
	if mode == "client" {
		hostname := getHostname()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		c := NewClient(socket)
		err := c.Dial()
		if err != nil {
			panic(err)
		}

		// send a msg every 10 seconds
		relayStateTicker := time.NewTicker(time.Second * 10)
		for range relayStateTicker.C {
			udpMsg := fmt.Sprintf("UDP message sent from client %s", hostname)
			err = c.Send(udpMsg)
			if err != nil {
				fmt.Printf("failed to send client message: %v", err)
			}
		}
		<-ctx.Done()
	}
}

func getHostname() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return name
}
