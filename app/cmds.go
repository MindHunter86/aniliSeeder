package app

import (
	"bytes"
	"io"
)

type rpcCommand uint8

const (
	cmdRpcUndefined rpcCommand = iota
	cmdsRpcGetTorrents
)

type cmds struct{}

func newCmds() *cmds { return &cmds{} }

func (*cmds) getAvaliableTorrentHashes() (io.Reader, error) {
	var buf *bytes.Buffer
	buf = bytes.NewBufferString("")
	buf.WriteString("fuckyou nigga\n\n")
	return buf, nil
}
