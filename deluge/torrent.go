package deluge

import (
	"encoding/json"
	"math"
	"os"

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

func (m *Client) GetTorrentsHashes() ([]string, error) {
	var e error
	var trrs map[string]*delugeclient.TorrentStatus

	if trrs, e = m.deluge.TorrentsStatus(delugeclient.StateUnspecified, nil); e != nil {
		return nil, e
	}

	var hashes []string
	for hash := range trrs {
		hashes = append(hashes, hash)
	}

	gLog.Debug().Int("hashes_length", len(hashes)).Msg("the torrnets hashes has been collected")
	return hashes, e
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

		gLog.Info().Str("hash", hash).Str("torrnet_name", torrent.Name).Msg("marking as weak ...")
		weakTrrs = append(weakTrrs, torrent)
	}

	return weakTrrs, e
}
