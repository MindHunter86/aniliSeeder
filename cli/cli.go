package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func TestDial(c *cli.Context, _ string) (err error) {
	log.Println("trying to connect via unix socket")

	conn, err := net.Dial("unix", c.String("socket-path"))
	if err != nil {
		return
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

	var buf = bytes.NewBuffer(nil)

	for {
		buf.Reset()

		pr := promptui.Prompt{
			Label: ":>",
			// IsVimMode: true,
			Templates: &promptui.PromptTemplates{
				Prompt:  "{{ . }} ",
				Valid:   "{{ . | green }} ",
				Invalid: "{{ . | red }} ",
				Success: "{{ . | bold }} ",
			},
			AllowEdit: true,
		}

		var data string
		data, err = pr.Run()

		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print(":> ")
		// data, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		// log.Println("readed input")

		buf.WriteString(data + "\n")
		_, err = io.Copy(conn, buf)
		if err != nil {
			return
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
			return
		}

		// log.Println("received response")
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}
