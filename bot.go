package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

var pingRegex = regexp.MustCompile("PING :(.*)$")

type Bot struct {
	port     string
	server   string
	username string
	realname string
	channel  string
	ident    string
	conn     net.Conn
}

// listenToIRCMessages returns the channel which receives lines from the irc connection
func listenToIRCMessages(conn net.Conn) <-chan []byte {
	c := make(chan []byte)
	reader := bufio.NewReader(conn)
	var line []byte
	var err error
	go func() {
		for {
			if line, err = reader.ReadBytes('\n'); err != nil {
				log.Println(err)
			}
			// send the line through the channel
			c <- line
			// respond to the server PING request with a PONG by doing a simple replace
			if pingRegex.Match([]byte(line)) {
				conn.Write(bytes.Replace(line, []byte("PING"), []byte("PONG"), 1))
			}
		}
	}()
	return c
}

// Connect to the IRC server
func (b *Bot) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", b.server, b.port))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Identify the bot
	fmt.Fprintf(conn, "NICK %s\r\n", b.username)
	fmt.Fprintf(conn, "USER %s 1 1 1:%s\r\n", b.ident, b.realname)
	// Join the provided channel
	fmt.Fprintf(conn, "JOIN %s\r\n", b.channel)

	c := listenToIRCMessages(conn)

	for {
		select {
		case line := <-c:
			log.Printf("Bot said: %q", line)
		case <-time.After(time.Second * 5):
			log.Println("No lines to report in the last 5 seconds")
		}
	}

}

// NewBot returns a new Bot
func NewBot() *Bot {
	return &Bot{
		username: "botko",
		realname: "Testbot",
		ident:    "testbot",
		channel:  "#ggtst",
		server:   "irc.freenode.com",
		port:     "6667",
		conn:     nil,
	}
}

func main() {
	if err := NewBot().Connect(); err != nil {
		log.Fatal(err)
	}
}
