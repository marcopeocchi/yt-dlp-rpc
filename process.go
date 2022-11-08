package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
)

var driver = os.Getenv("YT_DLP_PATH")

// "title":"%(info.title)s",
const template = `download:
{
	"thumbnail":"%(info.thumbnail)s", 
	"eta":%(progress.eta)s, 
	"resolution":"%(info.resolution)s", 
	"percentage":"%(progress._percent_str)s",
	"speed":%(progress.speed)s, 
	"size":"%(progress._total_bytes_str)s"
}`

type ProgressTemplate struct {
	Resolution string `json:"resolution"`
	Percentage string `json:"percentage"`
	Thumbnail  string `json:"thumbnail"`
	// Title      string  `json:"title"` TODO: effective way to convert unicode titles
	Speed float32 `json:"speed"`
	Size  string  `json:"size"`
	Eta   int     `json:"eta"`
}

// Process descriptor
type Process struct {
	id       string
	mem      *MemoryDB
	url      string
	params   []string
	pid      int
	progress Progress
}

// Starts spawns/forks a new yt-dlp process and parse its stdout.
// The process is spawned to outputting a custom progress text that
// Resembles a JSON Object in ordert to Unmarshall it later.
// This approach is anyhow not perfect. Unicode strings are not escaped
// so unmarshall a json is not that trivial.
// Each process is identified not by its PID but by a UUIDv2
func (p *Process) Start() {
	params := append([]string{
		p.url,
		"--newline",
		"--no-colors",
		"--progress-template", strings.ReplaceAll(template, "\n", ""),
	}, p.params...)
	params = append(params, "-o", "./%(title)s.%(ext)s")

	cmd := exec.Command(driver, params...)
	r, err := cmd.StdoutPipe()
	if err != nil {
		log.Panicln(err)
	}
	scan := bufio.NewScanner(r)

	err = cmd.Start()
	if err != nil {
		log.Panicln(err)
	}

	p.pid = cmd.Process.Pid
	p.id = p.mem.Set(p)

	// spawn a goroutine that does the dirty job of parsing the stdout
	go func() {
		defer cmd.Wait()
		defer r.Close()
		defer p.Kill()
		for scan.Scan() {
			stdout := ProgressTemplate{}
			err := json.Unmarshal([]byte(scan.Text()), &stdout)
			if err == nil {
				p.mem.Update(p.id, Progress{
					Percentage: stdout.Percentage,
					Thumbnail:  stdout.Thumbnail,
					// Title:      stdout.Title,
					Speed: stdout.Speed,
					Size:  stdout.Size,
					URL:   p.url,
					Id:    p.id,
				})
			}
		}
	}()
}

// Kill a process and remove it from the memory
func (p *Process) Kill() error {
	cmd := exec.Command("kill", []string{string(rune(p.pid))}...)
	cmd.Start()
	err := cmd.Wait()
	p.mem.Delete(p.id)
	log.Printf("Killed process %d\n", p.pid)
	return err
}
