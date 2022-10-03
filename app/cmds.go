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

func (*cmds) getAvaliableTorrentHashes() (io.ReadWriter, error) {
	var buf = bytes.NewBufferString("")

	hashes, err := gDeluge.GetTorrentsHashes()
	if err != nil {
		return nil, err
	}

	for _, hash := range hashes {
		buf.WriteString(hash + "\n")
	}

	return buf, nil
}
