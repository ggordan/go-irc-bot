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
	port     int
	server   string
	username string
	realname string
	channel  string
	ident    string
	conn     net.Conn
}

func (b Bot) close() {
	b.conn.Close()
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

// Connect to the IRC server and returns a channel
func (b *Bot) Connect() (<-chan []byte, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", b.server, b.port))
	if err != nil {
		return nil, err
	}
	b.conn = conn

	// Identify the bot
	fmt.Fprintf(conn, "NICK %s\r\n", b.username)
	fmt.Fprintf(conn, "USER %s 1 1 1:%s\r\n", b.ident, b.realname)
	// Join the provided channel
	fmt.Fprintf(conn, "JOIN %s\r\n", b.channel)

	return listenToIRCMessages(conn), nil
}

// NewBot returns a new Bot
func NewBot(userName, realName, ident, channel, server string, port int) *Bot {
	return &Bot{
		username: userName,
		realname: realName,
		ident:    ident,
		channel:  channel,
		server:   server,
		port:     port,
	}
}

func main() {

	// Bot 1
	bot1 := NewBot("botko1", "botko1", "botko1", "#ggtst", "irc.freenode.net", 6667)
	botChan1, err := bot1.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer bot1.close()

	// Bot 2
	bot2 := NewBot("botko2", "botko2", "botko2", "#ggtst", "irc.freenode.net", 6667)
	botChan2, err := bot2.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer bot2.close()

	for {
		select {
		case line := <-botChan1:
			log.Printf("Bot 1 said: %q", line)
		case line := <-botChan2:
			log.Printf("Bot 2 said: %q", line)
		case <-time.After(time.Second * 5):
			log.Println("No lines to report in the last 5 seconds")
		}
	}

}
