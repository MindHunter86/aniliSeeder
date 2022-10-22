package anilibria

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type session struct{}

func newSession() *session {
	return &session{}
}

func (m *session) getActiveAniSessions(body *[]byte) (_ *map[string][]string, e error) {
	return m.parseRawCPPage(bytes.NewBuffer(*body))
}

func (*session) parseRawCPPage(buf io.Reader) (_ *map[string][]string, e error) {

	var z = html.NewTokenizer(buf)

	var isSessTable, isTableBody, isTd bool
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
						gLog.Debug().Bool("isSessTable", isSessTable).Msg("session html table has been found")
						break
					}
				}
			case atom.Tbody:
				isTableBody = isSessTable
				gLog.Debug().Bool("isTableBpdy", isTableBody).Msg("html table body has been found")
			case atom.Td:
				isTd = isTableBody
				gLog.Trace().Bool("isTd", isTd).Msg("html table raw has been found")
			case atom.A:
				if !isTd {
					continue
				}

				gLog.Trace().Bool("isTd", isTd).Msg("html table raw link has been found")

				for _, attr := range tkn.Attr {
					gLog.Trace().Str("key", attr.Key).Str("value", attr.Val).Msg("")

					if attr.Key == "data-session-id" && attr.Val != "" {
						if tdCount%3 != 0 {
							gLog.Debug().Msg("html parser found data-session-id but there is no session details")
						}

						sessions[strings.TrimSpace(attr.Val)] = sessTdBuf
						sessTdBuf, isTd = nil, false

						gLog.Debug().Str("session", strings.TrimSpace(attr.Val)).Msg("session has been collected")
					}
				}
			}

		case html.TextToken:
			tkn := z.Token()

			if !isTd {
				continue
			}

			data := strings.TrimSpace(tkn.Data)

			if strings.TrimSpace(tkn.Data) == "" {
				continue
			}

			gLog.Trace().Str("shit", data).Msg("found some text shit")
			sessTdBuf = append(sessTdBuf, data)
			isTd = false
		}
	}

	// for id, sess := range sessions {
	// 	fmt.Println(id)
	// 	fmt.Println(sess)
	// }

	gLog.Debug().Int("sessions_length", len(sessions)).Msg("parsed sessions count")
	return &sessions, e
}
