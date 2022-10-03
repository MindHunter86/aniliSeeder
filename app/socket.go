package app

import (
	"io/ioutil"
	"log"
	"net"
	"os"
)

type SockServer struct {
	ln net.Listener
}

// !!
// socket.Close()

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

	go m.close()

	gLog.Debug().Msg("socket server has been bootstrapped successfully")
	return
}

func (m *SockServer) close() {
	<-gCtx.Done()
	gLog.Debug().Msg("context close() event has been caught; closing unix socket")
	m.ln.Close()
}

func (m *SockServer) Serve(done func()) error {
	defer done()

	gLog.Info().Msg("initiating unix socket serving...")
	defer gLog.Info().Msg("unix socket has been closed")

	for {
		conn, err := m.ln.Accept()
		if err == nil {
			go m.clientHandler(conn)
			continue
		}

		if err == net.ErrClosed {
			return err
		}

		return err
	}
}

func (m *SockServer) clientHandler(c net.Conn) {
	gLog.Info().Str("client", c.RemoteAddr().Network()).Msg("socket server: client connected")
	defer gLog.Info().Str("client", c.RemoteAddr().Network()).Msg("client disconnected")
	defer c.Close()

	msg, err := ioutil.ReadAll(c)
	if err != nil {
		gLog.Warn().Err(err).Msg("there are some errors with client communication")
		return
	}
	log.Println(string(msg))
	gLog.Debug().Str("message", string(msg)).Int("message_lentgh", len(msg)).Msg("there is new message from unix socket server client")

	gLog.Debug().Msg("trying to respond the client's message to the client")
	if n, err := c.Write(msg); err != nil {
		gLog.Warn().Err(err).Msg("there are some errors with client communication")
		return
	} else {
		gLog.Debug().Int("bytes_count", n).Msg("the server has been successfully responed")
	}
}
