package main

type Progress struct {
	Percentage string  `json:"percentage"`
	Thumbnail  string  `json:"thumbnail"`
	Title      string  `json:"title"`
	Speed      float32 `json:"speed"`
	Size       string  `json:"size"`
	ETA        int     `json:"eta"`
	URL        string  `json:"url"`
	Id         string  `json:"id"`
}
