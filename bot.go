package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strings"
)

type Bot struct {
	port     string
	server   string
	username string
	realname string
	channel  string
	ident    string
	conn     net.Conn
}

func (b *Bot) Connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", b.server, b.port))
	b.conn = conn

	defer b.conn.Close()

	if err != nil {
		panic(err)
	}

	pingRegex, _ := regexp.Compile("PING :(.*)$")

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	// identify bot
	b.Identify()
	// join channel

	// go func() {
	for {
		line, _ := tp.ReadLine()
		log.Printf("%s\n", line)
		if pingRegex.Match([]byte(line)) {
			text := strings.Replace(line, "PING", "PONG", 1)
			fmt.Fprintf(b.conn, "%s\r\n", text)
		}
		b.Join()
	}
	// }()

}

func (b *Bot) Identify() {
	// fmt.Fprintf(b.conn, "PASS irc.foonet.com\r\n")
	fmt.Fprintf(b.conn, "NICK %s\r\n", b.username)
	fmt.Fprintf(b.conn, "USER %s 1 1 1:%s\r\n", b.ident, b.realname)
}

func (b *Bot) Join() {
	fmt.Fprintf(b.conn, "JOIN %s\r\n", b.channel)
}

// Bot constructor
func NewBot() *Bot {
	return &Bot{
		username: "botko",
		realname: "Testbot",
		ident:    "testbot",
		channel:  "#test",
		server:   "irc.foonet.com",
		port:     "6667",
		conn:     nil,
	}
}

func main() {
	bot := NewBot()
	bot.Connect()
	bot.Join()
}
