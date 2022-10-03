package cmd

import (
	"bytes"
	"log"
	"net"

	"github.com/urfave/cli/v2"
)

func TestDial(c *cli.Context, test string) {
	log.Println("trying to connect via unix socket")

	conn, err := net.Dial("unix", c.String("socket-path"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("connection successfull; trying to sen some data")

	var buf = bytes.NewBufferString(test)
	if n, err := conn.Write(buf.Bytes()); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("good, close; written bytes %d\n", n)
	}
}
