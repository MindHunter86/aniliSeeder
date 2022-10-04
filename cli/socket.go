package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/urfave/cli/v2"
)

func TestDial(c *cli.Context, _ string) {
	log.Println("trying to connect via unix socket")

	conn, err := net.Dial("unix", c.String("socket-path"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("connection successfull; trying to sen some data")

	// var buf = bytes.NewBufferString(test)
	// if n, err := conn.Write(buf.Bytes()); err != nil {
	// !!
	// log.Fatal will exit, and defer conn.Close() will not run
	// CRT-D0011
	// !!
	// 	log.Fatal(err)
	// } else {
	// 	log.Printf("good, close; written bytes %d\n", n)
	// }

	// !!
	// TODO
	// github.com/tcnksm/go-input

	var buf bytes.Buffer
	for {
		buf.Reset()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(":> ")
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// log.Println("readed input")

		buf.WriteString(data)
		_, err = io.Copy(conn, &buf)
		if err != nil {
			log.Fatal(err)
		}

		// log.Println("sent input")

		buf.Reset()

		scanner := bufio.NewScanner(conn)
		var lines []string
		for {
			scanner.Scan()
			line := scanner.Text()
			if len(line) == 0 {
				break
			}

			lines = append(lines, line)
		}

		if scanner.Err() != nil {
			log.Fatal(scanner.Err())
		}

		// log.Println("received response")
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}
