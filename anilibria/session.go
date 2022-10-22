package anilibria

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type session struct{}

func newSession() *session {
	return &session{}
}

func (*session) getActiveAniSessions() (_ []string, e error) {
	return
}

func (*session) parseRawCPPage(buf io.ReadCloser) (e error) {

	var z = html.NewTokenizer(buf)

	var isSessTable, isParsable, isTd bool
	var tdCount int

	var sessTdBuf []string
	var sessions = make(map[string][]string)

loop:
	for {
		switch z.Next() {
		case html.ErrorToken:
			if z.Err() != io.EOF {
				return
			}
			break loop

		case html.StartTagToken:
			tkn := z.Token()

			switch tkn.DataAtom {
			case atom.Table:
				for _, attr := range tkn.Attr {
					if isSessTable = attr.Key == "id" && attr.Val == "tableSess"; isSessTable {
						continue
					}
				}
			case atom.Tbody:
				isParsable = isSessTable
			case atom.Td:
				isTd = isParsable
			case atom.A:
				if !isTd {
					continue
				}

				for _, attr := range tkn.Attr {
					if attr.Key == "data-session-id" && attr.Val != "" {
						if tdCount%3 != 0 {
							gLog.Debug().Msg("html parser found data-session-id but there is no session details")
						}

						sessions[strings.TrimSpace(attr.Val)] = sessTdBuf
						sessTdBuf, isTd = nil, false
					}
				}
			}

		case html.TextToken:
			tkn := z.Token()

			if !isTd {
				continue
			}

			sessTdBuf = append(sessTdBuf, strings.TrimSpace(tkn.Data))
			isTd = false
		}
	}

	gLog.Debug().Int("sessions_length", len(sessions)).Msg("parsed sessions count")

	for id, sess := range sessions {
		fmt.Println(id)
		fmt.Println(sess)
	}

	return
}
