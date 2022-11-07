package main

type Running []Progress
type Pending []int
type Progress struct {
	Percentage string  `json:"percentage"`
	Thumbnail  string  `json:"thumbnail"`
	Speed      float32 `json:"speed"`
	Size       string  `json:"size"`
	ETA        int     `json:"eta"`
	URL        string  `json:"url"`
	PID        int     `json:"pid"`
}
