package app

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type SockServer struct {
	ln  net.Listener
	cmd *cmds
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

	m.cmd = newCmds()
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
			go m.clientRpcHandler(conn)
			continue
		}

		if err == net.ErrClosed {
			return err
		}

		return err
	}
}

func (m *SockServer) clientTestHandler(c net.Conn) {
	var clientId = c.RemoteAddr().Network()

	gLog.Info().Str("client", clientId).Msg("socket server: client connected")
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

func (m *SockServer) clientRpcHandler(c net.Conn) {
	var clientId = c.RemoteAddr().Network()

	gLog.Info().Str("client", clientId).Msg("socket server: client connected")
	defer gLog.Info().Str("client", c.RemoteAddr().Network()).Msg("client disconnected")
	defer c.Close()

	for {
		msg, err := ioutil.ReadAll(c)
		if err != nil {
			gLog.Warn().Err(err).Str("client", clientId).Msg("there are some errors with client communication")
			return
		}

		gLog.Info().Str("client", clientId).Str("cmd", string(msg)).Msg("received a cmd from the client")

		var clientCmd rpcCommand
		if clientCmd = m.parseClientCmd(string(msg)); clientCmd == cmdRpcUndefined {
			gLog.Warn().Str("client", clientId).Str("cmd", string(msg)).Msg("received cmd is undefined")

			var buf = bytes.NewBufferString("command not found")
			if n, err := io.Copy(c, buf); m.checkRespondErrors(n, err, string(msg), clientId) != nil {
				return
			}
		}

		var buf io.Reader
		if buf, err = m.runClientCmd(clientCmd); err != nil {
			gLog.Warn().Str("client", clientId).Str("cmd", string(msg)).Err(err).Msg("could not run received cmd because of internal errors")

			var buf = bytes.NewBufferString("internal server error")
			if n, err := io.Copy(c, buf); m.checkRespondErrors(n, err, string(msg), clientId) != nil {
				return
			}
		}

		if n, err := io.Copy(c, buf); m.checkRespondErrors(n, err, string(msg), clientId) != nil {
			return
		}
	}
}

func (m *SockServer) checkRespondErrors(written int64, e error, cmd, id string) error {
	if e != nil {
		gLog.Warn().Err(e).Str("client", id).Str("cmd", cmd).Msg("there are some errors with client communication")
		return e
	}

	gLog.Debug().Str("client", id).Int64("response_length", written).Msg("successfully responded to the client")
	return nil
}

func (m *SockServer) parseClientCmd(cmd string) rpcCommand {
	switch cmd {
	case "getTorrents":
		return cmdsRpcGetTorrents
	default:
		return cmdRpcUndefined
	}
}

func (m *SockServer) runClientCmd(cmd rpcCommand) (io.Reader, error) {
	switch cmd {
	case cmdsRpcGetTorrents:
		return m.cmd.getAvaliableTorrentHashes()

	default:
		gLog.Error().Msg("golang internal error; given cmd is undefined in runClientCmd()")
	}

	// !!
	return nil, nil
}
