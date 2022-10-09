package app

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type rpcCommand uint8

const (
	cmdRpcUndefined rpcCommand = iota

	cmdsRpcGetTorrents
	cmdsRpcStatTorrents

	cmdWorkersList
	cmdLoadAniUpdates
	cmdLoadAniChanges
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

func (*cmds) statCurrentTorrents() (io.ReadWriter, error) {
	tb := table.NewWriter()
	defer tb.Render()

	var buf = bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{"Hash", "Name", "TotalSize", "Ratio", "Uploaded", "Seedtime", "VKScore"})

	trrs, e := gDeluge.GetTorrents()
	if e != nil {
		return nil, e
	}

	for hash, torrent := range trrs {
		name, _, _ := strings.Cut(torrent.Name, "- AniLibria.TV")
		seedTime := time.Duration(torrent.SeedingTime) * time.Second
		tb.AppendRow([]interface{}{
			hash[0:9], name, torrent.TotalSize / 1024 / 1024, torrent.Ratio, torrent.TotalUploaded / 1024 / 1024, seedTime.String(),
			// ??
			// todo optimize
			gDeluge.GetScoreFromInput(float64(torrent.TotalUploaded), seedTime.Hours()/float64(24), float64(torrent.TotalSize)),
		})
	}

	tb.SetRowPainter(func(raw table.Row) text.Colors {
		if raw[6].(float64) >= float64(gCli.Int("torrents-vkscore-line")) {
			return text.Colors{text.FgGreen}
		}
		if raw[3].(float32) < 1 {
			return text.Colors{text.FgRed}
		}
		return text.Colors{text.FgYellow}
	})

	tb.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, WidthMax: 60},
	})

	tb.SortBy([]table.SortBy{
		{Name: "VKScore", Mode: table.DscNumeric},
	})

	return buf, nil
}

func (*cmds) listWorkers() (_ io.ReadWriter, e error) {
	tb := table.NewWriter()
	defer tb.Render()

	var buf = bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{"ID", "Version", "FreeSpaceMB", "ActiveTorrents"})

	for id, wrk := range gSwarm.GetConnectedWorkers() {
		tb.AppendRow([]interface{}{
			id, wrk.Version, wrk.FreeSpace / 1024 / 1024, len(wrk.ActiveTorrents),
		})
	}

	tb.SortBy([]table.SortBy{
		{Name: "ID", Mode: table.Dsc},
	})

	return buf, e
}

func (*cmds) loadAniUpdates() (_ io.ReadWriter, e error) {
	tb := table.NewWriter()
	defer tb.Render()

	buf := bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{
		"ID", "Name", "Status", "Type", "Series", "Torrent", "Size", "Seeders", "Leechers",
	})

	var titles []*anilibria.Title
	if titles, e = gAniApi.GetLastUpdates(); e != nil {
		return
	}

	for _, tl := range titles {
		for _, tr := range tl.Torrents.List {
			tb.AppendRow([]interface{}{
				tl.Id, tl.Names.Ru, tl.Status.String, tl.Type.FullString, tl.Torrents.Series.String,
				tr.Hash[0:9], utils.GetMBytesFromBytes(tr.TotalSize), tr.Seeders, tr.Leechers,
			})

		}
	}

	tb.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
		{Number: 4, AutoMerge: true},
		{Number: 5, AutoMerge: true},
	})
	tb.Style().Options.SeparateRows = true

	return buf, e
}
