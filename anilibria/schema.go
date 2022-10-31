package anilibria

type (
	apiError struct {
		Error *apiErrorDetails
	}
	apiErrorDetails struct {
		Code    int
		Message string
	}
	TitleSchedule struct {
		Day  int
		List []*Title
	}
	Title struct {
		Id         int
		Code       string
		Updated    uint64 // sometimes the anilibria project mark their update time as a NULL
		LastChange uint64 `json:"last_change"` // I dont know how to mark this fields as "if time.Parse fails - ignore"
		Names      *TitleNames
		Status     *TitleStatus
		Type       *TitleType
		Torrents   *TitleTorrents
	}
	TitleNames struct {
		Ru          string
		En          string
		Alternative string
	}
	TitleStatus struct {
		String string
		Code   int
	}
	TitleType struct {
		FullString string `json:"full_string"`
		Code       int
		String     string
		Series     interface{}
		Length     int
	}
	TitleTorrents struct {
		Series *TorrentSeries
		List   []*TitleTorrent
	}
	TitleTorrent struct {
		TorrentId         int `json:"torrent_id"`
		Series            *TorrentSeries
		Quality           *TorrentQuality
		Leechers          int
		Seeders           int
		Downloads         int
		TotalSize         int64 `json:"total_size"`
		Url               string
		UploadedTimestamp uint64 `json:"uploaded_timestamp"`
		Hash              string
		Metadata          *TorrentMetadata
		RawBase64File     interface{}
	}
	TorrentSeries struct {
		Firest int
		Last   int
		String string
	}
	TorrentQuality struct {
		String     string
		Type       string
		Resolution string
		Encoder    string
		LqAudio    interface{} `json:"lq_audio"`
	}
	TorrentMetadata struct {
		Hash             string
		Name             string
		Announce         []string
		CreatedTimestamp uint64          `json:"created_timestamp"`
		FilesList        []*MetadataFile `json:"files_list"`
	}
	MetadataFile struct {
		File   string
		Size   uint64
		Offset uint64
	}
)

// easyjson hacks:

//easyjson:json
type Titles []*Title
