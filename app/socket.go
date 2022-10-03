package app

import (
	"io"
	"net"
	"os"
)

type SockServer struct {
	ln net.Listener
}

func NewSockServer() *SockServer {
	return &SockServer{}
}

func (m *SockServer) Bootstrap() (e error) {
	if e = os.RemoveAll(gCli.String("socket-path")); e != nil {
		return
	}

	var lc net.ListenConfig
	if m.ln, e = lc.Listen(gCtx, "unix", gCli.String("socket-path")); e != nil {
		return
	}

	return
}

func (m *SockServer) Serve(done func()) error {
	defer done()

	for {
		conn, err := m.ln.Accept()
		if err != nil {
			go m.clientHandler(conn)
			continue
		}

		if err == net.ErrClosed {
			return nil
		}

		return err
	}
}

func (m *SockServer) clientHandler(c net.Conn) {
	gLog.Info().Str("client", c.RemoteAddr().Network()).Msg("socket server: client connected")
	io.Copy(c, c)
	c.Close()
}
