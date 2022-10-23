package app

import (
	"bytes"
	"io"
	"time"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type rpcCommand uint8

const (
	cmdRpcUndefined rpcCommand = iota

	cmdsGetTorrents

	cmdWorkersList

	cmdLoadAniUpdates
	cmdLoadAniChanges
	cmdLoadAniSchedule

	cmdDeployAniUpdates
	cmdDryDeployAniUpdates
	// cmdDryDeployAniChanges
	// cmdDryDeployAniSchedule
)

type cmds struct{}

func newCmds() *cmds { return &cmds{} }

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

func (*cmds) getMasterTorrents() (_ io.ReadWriter, e error) {

	tb := table.NewWriter()
	defer tb.Render()

	var buf = bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{"Worker", "Hash", "Name", "TotalSize", "Ratio", "Uploaded", "Seedtime", "Announce", "VKScore"})

	for id, wrk := range gSwarm.GetConnectedWorkers() {
		for _, trr := range wrk.ActiveTorrents {
			seedTime := time.Duration(trr.SeedingTime) * time.Second
			tb.AppendRow([]interface{}{
				id[0:8], trr.GetShortHash(), trr.GetName(), trr.TotalSize / 1024 / 1024, trr.Ratio,
				trr.TotalUploaded / 1024 / 1024, seedTime.String(), trr.GetTrackerStatus(), trr.GetVKScore(),
			})

		}
	}

	tb.SetRowPainter(func(raw table.Row) text.Colors {
		if raw[8].(float64) >= float64(gCli.Int("torrents-vkscore-line")) && raw[7] == "OK" {
			return text.Colors{text.FgGreen}
		}
		if raw[4].(float32) < 1 {
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
				tl.Id, tl.Names.Ru, tl.Status.String, tl.Type.String, tl.Torrents.Series.String,
				tr.GetShortHash(), tr.TotalSize / 1024 / 1024, tr.Seeders, tr.Leechers,
			})

		}
	}

	tb.SortBy([]table.SortBy{
		{Name: "Name", Mode: table.Asc},
	})

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
func (*cmds) loadAniChanges() (_ io.ReadWriter, e error) {
	tb := table.NewWriter()
	defer tb.Render()

	buf := bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{
		"ID", "Name", "Status", "Type", "Series", "Torrent", "Size", "Seeders", "Leechers",
	})

	var titles []*anilibria.Title
	if titles, e = gAniApi.GetLastChanges(); e != nil {
		return
	}

	for _, tl := range titles {
		for _, tr := range tl.Torrents.List {
			tb.AppendRow([]interface{}{
				tl.Id, tl.Names.Ru, tl.Status.String, tl.Type.String, tl.Torrents.Series.String,
				tr.GetShortHash(), tr.TotalSize / 1024 / 1024, tr.Seeders, tr.Leechers,
			})

		}
	}

	tb.SortBy([]table.SortBy{
		{Name: "Name", Mode: table.Asc},
	})

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

func (*cmds) loadAniSchedule() (_ io.ReadWriter, e error) {
	tb := table.NewWriter()
	defer tb.Render()

	buf := bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{
		"Weekday", "ID", "Name", "Status", "Type", "Series", "Torrent", "Size", "Seeders", "Leechers",
	})

	var schedule []*anilibria.TitleSchedule
	if schedule, e = gAniApi.GetTitlesSchedule(); e != nil {
		return
	}

	for _, day := range schedule {
		for _, tl := range day.List {
			for _, tr := range tl.Torrents.List {
				tb.AppendRow([]interface{}{
					day.Day, tl.Id, tl.Names.Ru, tl.Status.String, tl.Type.String, tl.Torrents.Series.String,
					tr.GetShortHash(), tr.TotalSize / 1024 / 1024, tr.Seeders, tr.Leechers,
				})
			}
		}
	}

	tb.SortBy([]table.SortBy{
		{Name: "Weekday", Mode: table.Asc},
	})

	tb.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
		{Number: 4, AutoMerge: true},
		{Number: 5, AutoMerge: true},
		{Number: 6, AutoMerge: true},
	})
	tb.Style().Options.SeparateRows = true

	return buf, e
}

func (*cmds) dryDeployAniUpdates() (_ io.ReadWriter, e error) {
	tb := table.NewWriter()
	defer tb.Render()

	buf := bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{
		"Worker", "Torrent", "Size", "Seeders", "Leechers", "Uploaded",
	})

	dpl := newDeploy()
	var deployTitles = make(map[string][]anilibria.TitleTorrent)

	if deployTitles, e = dpl.dryRun(); e != nil {
		return
	}

	for wid, trrs := range deployTitles {
		for _, trr := range trrs {
			tb.AppendRow([]interface{}{
				wid[0:8], trr.GetShortHash(), trr.TotalSize / 1024 / 1024, trr.Seeders, trr.Leechers,
				time.Unix(int64(trr.UploadedTimestamp), 0).String(),
			})
		}
	}

	tb.SortBy([]table.SortBy{
		{Name: "Worker", Mode: table.Asc},
	})

	tb.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	tb.Style().Options.SeparateRows = true

	return buf, e
}

func (*cmds) deployAniUpdates() (_ io.ReadWriter, e error) {
	tb := table.NewWriter()
	defer tb.Render()

	buf := bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{
		"Worker", "Torrent", "Size", "Seeders", "Leechers", "Uploaded",
	})

	dpl := newDeploy()
	var deployTitles = make(map[string][]anilibria.TitleTorrent)

	if deployTitles, e = dpl.run(); e != nil {
		return
	}

	for wid, trrs := range deployTitles {
		for _, trr := range trrs {
			tb.AppendRow([]interface{}{
				wid[0:8], trr.GetShortHash(), trr.TotalSize / 1024 / 1024, trr.Seeders, trr.Leechers,
				time.Unix(int64(trr.UploadedTimestamp), 0).String(),
			})
		}
	}

	tb.SortBy([]table.SortBy{
		{Name: "Worker", Mode: table.Asc},
	})

	tb.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	tb.Style().Options.SeparateRows = true

	return buf, e
}
