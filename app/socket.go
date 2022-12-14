package app

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"os"
	"strings"
)

type SockServer struct {
	ln  net.Listener
	cmd *cmds
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

func (m *SockServer) clientRpcHandler(c net.Conn) {
	var clientId = c.RemoteAddr().Network()

	gLog.Info().Str("client", clientId).Msg("socket server: client connected")
	defer gLog.Info().Str("client", c.RemoteAddr().Network()).Msg("client disconnected")
	defer c.Close()

	var reader = bufio.NewReader(c)
	for {
		msg, err := reader.ReadString('\n')
		// msg, err := ioutil.ReadAll(c)
		if err != nil {
			gLog.Warn().Err(err).Str("client", clientId).Msg("there are some errors with client communication")
			return
		}

		gLog.Info().Str("client", clientId).Str("cmd", msg).Msg("received a cmd from the client")

		// !!
		// TODO
		// refactoring is needed
		var buf io.ReadWriter
		var clientCmd rpcCommand
		if clientCmd = m.parseClientCmd(msg); clientCmd == cmdRpcUndefined {
			gLog.Warn().Str("client", clientId).Str("cmd", msg).Msg("received cmd is undefined")

			buf = bytes.NewBufferString("command not found")
			if n, err := io.Copy(c, m.getResponseMessage(buf)); m.checkRespondErrors(n, err, msg, clientId) != nil {
				return
			}

			continue
		}

		if buf, err = m.runClientCmd(clientCmd); err != nil {
			gLog.Warn().Str("client", clientId).Str("cmd", msg).Err(err).Msg("could not run received cmd because of internal errors")

			buf = bytes.NewBufferString("internal server error: " + err.Error())
			if n, err := io.Copy(c, m.getResponseMessage(buf)); m.checkRespondErrors(n, err, msg, clientId) != nil {
				return
			}
		}

		// TODO
		// !! remove
		// if buf == nil {
		// 	buf = bytes.NewBuffer(nil)
		// }

		if n, err := io.Copy(c, m.getResponseMessage(buf)); m.checkRespondErrors(n, err, msg, clientId) != nil {
			return
		}
	}
}

func (*SockServer) getResponseMessage(rw io.ReadWriter) io.ReadWriter {
	_, err := rw.Write([]byte("\n\n"))
	if err != nil {
		gLog.Warn().Err(err).Msg("could not prepare response message because of internal golang error")
	}

	return rw
}

func (*SockServer) checkRespondErrors(written int64, e error, cmd, id string) error {
	if e != nil {
		gLog.Warn().Err(e).Str("client", id).Str("cmd", cmd).Msg("there are some errors with client communication")
		return e
	}

	gLog.Debug().Str("client", id).Int64("response_length", written).Msg("successfully responded to the client")
	return e
}

func (*SockServer) parseClientCmd(cmd string) rpcCommand {
	switch strings.TrimSpace(cmd) {
	case "getTorrents":
		return cmdsGetTorrents
	case "listWorkers":
		return cmdWorkersList
	case "aniUpdates":
		return cmdLoadAniUpdates
	case "aniChanges":
		return cmdLoadAniChanges
	case "aniSchedule":
		return cmdLoadAniSchedule
	case "dryDeployAniUpdates":
		return cmdDryDeployAniUpdates
	case "deployAniUpdates":
		return cmdDeployAniUpdates
	case "dryDeployAniChanges":
		return cmdDryDeployAniChanges
	case "deployAniChanges":
		return cmdDeployAniChanges
	case "getActiveSessions":
		return cmdGetActiveSessions
	case "dropActiveSessions":
		return cmdDropAllActiveSessions
	case "dryDeployFailedAnnounces":
		return cmdDryDeployFailedAnnounces
	case "deployFailedAnnounces":
		return cmdDeployFailedAnnounces
	default:
		gLog.Debug().Str("cmd", strings.TrimSpace(cmd)).Msg("trimmed")
		return cmdRpcUndefined
	}
}

func (m *SockServer) runClientCmd(cmd rpcCommand) (io.ReadWriter, error) {
	switch cmd {
	case cmdsGetTorrents:
		return m.cmd.getMasterTorrents()
	case cmdWorkersList:
		return m.cmd.listWorkers()
	case cmdLoadAniUpdates:
		return m.cmd.loadAniUpdates()
	case cmdLoadAniChanges:
		return m.cmd.loadAniChanges()
	case cmdLoadAniSchedule:
		return m.cmd.loadAniSchedule()
	case cmdDryDeployAniUpdates:
		return m.cmd.deployAniUpdates()
	case cmdDeployAniUpdates:
		return m.cmd.deployAniUpdates(false)
	case cmdDryDeployAniChanges:
		return m.cmd.deployAniChanges()
	case cmdDeployAniChanges:
		return m.cmd.deployAniChanges(false)
	case cmdGetActiveSessions:
		return m.cmd.getActiveSessions()
	case cmdDropAllActiveSessions:
		return m.cmd.dropAllActiveSessions()
	case cmdDryDeployFailedAnnounces:
		return m.cmd.deployFailedAnnounces(true)
	case cmdDeployFailedAnnounces:
		return m.cmd.deployFailedAnnounces(false)

	default:
		gLog.Error().Msg("golang internal error; given cmd is undefined in runClientCmd()")
	}

	// !!
	return nil, nil
}
