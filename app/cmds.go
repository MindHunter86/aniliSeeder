package app

import (
	"bytes"
	"io"
	"time"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
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

	cmdGetActiveSessions
	cmdDropAllActiveSessions

	cmdDryDeployFailedAnnounces
	cmdDeployFailedAnnounces
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
		var trrs []*deluge.Torrent
		if trrs, e = gSwarm.RequestTorrentsFromWorker(wrk.Id); e != nil {
			return
		}

		for _, trr := range trrs {
			seedTime := time.Duration(trr.SeedingTime) * time.Second
			tb.AppendRow([]interface{}{
				id[0:8], trr.GetShortHash(), trr.GetName(), trr.TotalSize / 1024 / 1024, trr.Ratio,
				trr.TotalUploaded / 1024 / 1024, seedTime.String(), trr.GetTrackerStatus(), trr.GetVKScore(),
			})

		}
	}

	tb.SetRowPainter(func(raw table.Row) text.Colors {
		if raw[8].(float64) >= float64(gCli.Int("cmd-vkscore-warn")) && raw[7] == "OK" {
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

func (*cmds) getActiveSessions() (_ io.ReadWriter, e error) {
	tb := table.NewWriter()
	defer tb.Render()

	buf := bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{
		"IP", "SID", "Time", "LifeTime", "Status",
	})

	var sessions *map[string][]string
	if sessions, e = gAniApi.GetActiveSessions(); e != nil {
		return
	}

	for sid, session := range *sessions {
		// tm, e := time.Parse(time.RFC3339, session[2])
		tm, e := time.Parse("2006-01-02 15:04", session[2])
		if e != nil {
			gLog.Warn().Err(e).Msg("got an error in active sessions cmd rendering")
		}

		tb.AppendRow([]interface{}{
			session[1], sid, tm.String(), time.Since(tm).String(), session[3],
		})
	}

	tb.SortBy([]table.SortBy{
		{Name: "IP", Mode: table.Asc},
		{Name: "Time", Mode: table.Asc},
	})

	tb.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	tb.Style().Options.SeparateRows = true

	return buf, e
}

func (*cmds) dropAllActiveSessions() (_ io.ReadWriter, e error) {
	var sessions *map[string][]string
	if sessions, e = gAniApi.GetActiveSessions(); e != nil {
		return
	}

	var sids []string
	for sid := range *sessions {
		sids = append(sids, sid)
	}

	gAniApi.DropActiveSessions(sids...)
	return bytes.NewBufferString("OK"), e
}

func (*cmds) deployFailedAnnounces(dryrun bool) (_ io.ReadWriter, e error) {
	tb := table.NewWriter()
	defer tb.Render()

	buf := bytes.NewBuffer(nil)
	tb.SetOutputMirror(buf)
	tb.AppendHeader(table.Row{
		"Worker", "Name", "Quality", "OldHash", "NewHash", "SizeChanges KB", // "Deploied" // TODO
	})

	var ftitles []*failedTitle
	if ftitles, e = newDeploy().deployFailedAnnounces(dryrun); e != nil {
		return
	}

	for _, ft := range ftitles {
		tb.AppendRow([]interface{}{
			ft.workerId[0:8], ft.oldTorrent.GetName(), ft.oldTorrent.GetQuality(),
			ft.oldTorrent.GetShortHash(), ft.aniTorrent.GetShortHash(), ft.sizeChanges / 1024,
		})
	}

	tb.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	tb.Style().Options.SeparateRows = true

	return buf, e
}
