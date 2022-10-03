package cmd

import (
	"bytes"
	"io"
	"log"
	"net"

	"github.com/urfave/cli/v2"
)

func TestDial(c *cli.Context, test string) {
	conn, err := net.Dial("unix", c.String("socket-path"))
	if err != nil {
		log.Fatal(err)
	}

	var buf = bytes.NewBuffer([]byte(test))
	io.Copy(buf, conn)
}
