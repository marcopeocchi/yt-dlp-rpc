package pkg

type DownloadProgress struct {
	Percentage string  `json:"percentage"`
	Speed      float32 `json:"speed"`
	ETA        int     `json:"eta"`
}

type DownloadInfo struct {
	URL        string `json:"url"`
	Title      string `json:"title"`
	Thumbnail  string `json:"thumbnail"`
	Resolution string `json:"resolution"`
	Size       string `json:"size"`
}

// struct representing the response sent to the client
// as JSON-RPC result field
type ProcessResponse struct {
	Id       string           `json:"id"`
	Progress DownloadProgress `json:"progress"`
	Info     DownloadInfo     `json:"info"`
}