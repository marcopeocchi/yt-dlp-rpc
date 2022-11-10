package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	driver = os.Getenv("YT_DLP_PATH")
)

const template = `download:
{
	"eta":%(progress.eta)s, 
	"percentage":"%(progress._percent_str)s",
	"speed":%(progress.speed)s
}`

type ProgressTemplate struct {
	Percentage string  `json:"percentage"`
	Speed      float32 `json:"speed"`
	Size       string  `json:"size"`
	Eta        int     `json:"eta"`
}

type DownloadInfo struct {
	Title      string `json:"title"`
	Thumbnail  string `json:"thumbnail"`
	Resolution string `json:"resolution"`
}

// Process descriptor
type Process struct {
	id       string
	pid      int
	url      string
	params   []string
	progress Progress
	mem      *MemoryDB
}

// Starts spawns/forks a new yt-dlp process and parse its stdout.
// The process is spawned to outputting a custom progress text that
// Resembles a YAML Object in order to Unmarshal it later.
// This approach is anyhow not perfect: quotes are not escaped properly.
// Each process is identified not by its PID but by a UUIDv2
func (p *Process) Start() {
	params := append([]string{
		strings.Split(p.url, "?list")[0], //no playlist
		"--newline",
		"--no-colors",
		"--progress-template", strings.ReplaceAll(template, "\n", ""),
	}, p.params...)
	params = append(params, "-o", "./downloads/%(title)s.%(ext)s")

	// ----------------- main block ----------------- //
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

	// ----------------- info block ----------------- //
	// spawn a goroutine that retrieves the info for the download
	go func() {
		cmd := exec.Command(driver, p.url, "-J")
		stdout, err := cmd.Output()
		if err != nil {
			log.Println("Cannot retrieve info for", p.url)
		}
		info := DownloadInfo{}
		json.Unmarshal(stdout, &info)
		p.mem.Update(p.id, Progress{
			Title:      info.Title,
			Thumbnail:  info.Thumbnail,
			Resolution: info.Resolution,
			Id:         p.id,
		})
	}()

	// --------------- end info block --------------- //

	// spawn a goroutine that does the dirty job of parsing the stdout
	eventChan := make(chan string)

	// fill the channel with as many stdout line as yt-dlp produces (producer)
	go func() {
		defer cmd.Wait()
		defer r.Close()
		defer p.Kill()
		for scan.Scan() {
			eventChan <- scan.Text()
		}
	}()

	// debounce the unmarshal operation by 500ms (consumer)
	go debounce(time.Millisecond*500, eventChan, func(text string) {
		stdout := ProgressTemplate{}
		err := json.Unmarshal([]byte(text), &stdout)
		if err == nil {
			p.mem.UpdateProgress(p.id, Progress{
				Percentage: stdout.Percentage,
				Speed:      stdout.Speed,
				ETA:        stdout.Eta,
				Id:         p.id,
			})
		}
	})
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
