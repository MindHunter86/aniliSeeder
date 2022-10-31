package deluge

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	delugeclient "github.com/MindHunter86/go-libdeluge"
)

func (m *Client) GetTorrentsStatus() (e error) {
	var trrs map[string]*delugeclient.TorrentStatus
	if trrs, e = m.deluge.TorrentsStatus(delugeclient.StateUnspecified, nil); e != nil {
		return
	}

	je := json.NewEncoder(os.Stdout)
	je.SetIndent("", "  ")

	if e = je.Encode(trrs); e != nil {
		return
	}

	return
}

func (m *Client) GetTorrents() (map[string]*delugeclient.TorrentStatus, error) {
	return m.deluge.TorrentsStatus(delugeclient.StateUnspecified, nil)
}

func (m *Client) GetTorrentsHashes() (thashes []string, e error) {
	var trrs = make(map[string]*delugeclient.TorrentStatus)
	if trrs, e = m.GetTorrents(); e != nil {
		return
	}

	for hash := range trrs {
		thashes = append(thashes, hash)
	}

	gLog.Debug().Int("hashes_length", len(thashes)).Msg("the torrents hashes has been collected")
	return
}

// TODO:
// weak score formula (VKSCORE):
//
// uploaded / seed time * 100
// --------------------------
// 			size
//
// formula is valid for for torrents with ratio >= 1
// if score < 25 - give weak flag for torrent
// if torrent has 3 weak flags - drop
//
// ratio notice:
// if seed time > N days and ratio < 1 = drop torrent

func (m *Client) GetScoreFromInput(upld, seed, size float64) float64 {
	return m.roundGivenScore(upld/seed*100/size, 3)
}

func (*Client) roundGivenScore(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (m *Client) GetWeakTorrents() ([]*delugeclient.TorrentStatus, error) {
	trrs, e := m.GetTorrents()
	if e != nil {
		return nil, e
	}

	var weakTrrs []*delugeclient.TorrentStatus

	for hash, torrent := range trrs {
		if torrent.SeedingTime < 86400 {
			continue
		}
		if torrent.Ratio > 0.5 {
			continue
		}

		gLog.Info().Str("hash", hash).Str("torrent_name", torrent.Name).Msg("marking as weak ...")
		weakTrrs = append(weakTrrs, torrent)
	}

	return weakTrrs, e
}

type Torrent struct {
	Hash          string
	ActiveTime    int64
	Ratio         float32
	IsFinished    bool
	IsSeed        bool
	Name          string
	NumPeers      int64
	NumPieces     int64
	NumSeeds      int64
	PieceLength   int64
	SeedingTime   int64
	State         string
	TotalDone     int64
	TotalPeers    int64
	TotalSeeds    int64
	TotalSize     int64
	TotalUploaded int64
	TrackerStatus string

	Files []*TorrentFile
	// Peers          []delugeclient.Peer
	// FilePriorities []int64
	// FileProgress   []float32
}
type TorrentFile struct {
	Index  int64
	Size   int64
	Offset int64
	Path   string
}

func (m *Client) newTorrentFromStatus(hash string, t *delugeclient.TorrentStatus) *Torrent {
	return &Torrent{
		Hash:          hash,
		ActiveTime:    t.ActiveTime,
		Ratio:         t.Ratio,
		IsFinished:    t.IsFinished,
		IsSeed:        t.IsSeed,
		Name:          t.Name,
		NumPeers:      t.NumPeers,
		NumPieces:     t.NumPieces,
		NumSeeds:      t.NumSeeds,
		PieceLength:   t.PieceLength,
		SeedingTime:   t.SeedingTime,
		State:         t.State,
		TotalPeers:    t.TotalPeers,
		TotalSeeds:    t.TotalSeeds,
		TotalDone:     t.TotalDone,
		TotalUploaded: t.TotalUploaded,
		TotalSize:     t.TotalSize,
		TrackerStatus: t.TrackerStatus,

		Files: m.newTorrentFilesFromStatus(&t.Files),
	}
}

func (*Client) newTorrentFilesFromStatus(t *[]delugeclient.File) []*TorrentFile {
	var tfiles []*TorrentFile

	for _, file := range *t {
		tfiles = append(tfiles, &TorrentFile{
			Index:  file.Index,
			Size:   file.Size,
			Offset: file.Offset,
			Path:   file.Path,
		})
	}

	return tfiles
}

func (m *Client) GetTorrentsV2() (_ []*Torrent, e error) {
	var trrs map[string]*delugeclient.TorrentStatus
	if trrs, e = m.GetTorrents(); e != nil {
		return
	}

	var trrs2 []*Torrent
	for h, t := range trrs {
		trrs2 = append(trrs2, m.newTorrentFromStatus(h, t))
	}

	return trrs2, e
}

func (*Client) SaveTorrentFile(fname string, buf io.Reader) (_ int64, e error) {
	path := filepath.Join(gCli.String("deluge-torrents-path"), fname)

	if _, e = os.Stat(path); e != nil {
		if !os.IsNotExist(e) {
			gLog.Debug().Err(e).Str("path", path).Msg("given path is already exists; drop request")
			return
		}

		gLog.Debug().Err(e).Str("path", path).Msg("given path was not found; continue ...")
	}

	var fd *os.File
	if fd, e = os.Create(path); e != nil {
		return
	}
	defer fd.Close()

	return io.Copy(fd, buf)
}

func (m *Client) RemoveTorrent(hash string, withData bool) (bool, error) {
	return m.deluge.RemoveTorrent(hash, withData)
}

func (m *Client) TorrentStatus(hash string) (_ *Torrent, e error) {
	var tstatus *delugeclient.TorrentStatus
	if tstatus, e = m.deluge.TorrentStatus(hash); e != nil {
		return
	}

	return m.newTorrentFromStatus(hash, tstatus), e
}

func (m *Client) ForceReannounce(hashes ...string) (e error) {
	return m.deluge.ForceReannounce(hashes)
}

func (m *Torrent) GetVKScore() (_ float64) {
	seedtime := time.Duration(m.SeedingTime) * time.Second
	seeddays := seedtime.Hours() / float64(24)
	return m.roundGivenScore(float64(m.TotalUploaded)/seeddays*100/float64(m.TotalSize), 3)
}

func (*Torrent) roundGivenScore(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (m *Torrent) GetTrackerStatus() TrackerStatus {
	switch m.TrackerStatus {
	case "Announce OK":
		return TrackerStatusOK
	case "Announce Sent":
		return TrackerStatusSent
	case "Error: Connection timed out":
		return TrackerStatusConnTimedOut
	case "Error: timed out":
		return TrackerStatusTimedOut
	case "Error: Торрент не зарегистрирован":
		return TrackerStatusNotRegistered
	default:
		return TrackerStatusUnknown
	}
}

func (m *Torrent) GetTrackerRawError() string {
	return m.TrackerStatus
}

func (m *Torrent) GetName() string {
	// https://github.com/MindHunter86/aniliSeeder/issues/74
	name := strings.ReplaceAll(m.Name, "_", " ")

	name, _, _ = strings.Cut(name, "- AniLibria.TV")
	return strings.TrimSpace(name)
}

func (m *Torrent) GetShortHash() string {
	return m.Hash[0:9]
}

func (m *Torrent) GetQuality() string {
	// https://github.com/MindHunter86/aniliSeeder/issues/74
	name := strings.ReplaceAll(m.Name, "_", " ")

	// strings.Trim("][") is not worked here; and I don't know why...
	_, rawquality, _ := strings.Cut(strings.Trim(name, "]"), "[")
	return strings.TrimSpace(rawquality)
}
