package pkg

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"goytdlp.rpc/m/pkg/rx"
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

// Process descriptor
type Process struct {
	id       string
	url      string
	params   []string
	Info     DownloadInfo
	Progress DownloadProgress
	mem      *MemoryDB
	proc     *os.Process
}

// Starts spawns/forks a new yt-dlp process and parse its stdout.
// The process is spawned to outputting a custom progress text that
// Resembles a JSON Object in order to Unmarshal it later.
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

	p.id = p.mem.Set(p)
	p.proc = cmd.Process

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
		p.mem.Update(p.id, DownloadInfo{
			URL:        p.url,
			Title:      info.Title,
			Thumbnail:  info.Thumbnail,
			Resolution: info.Resolution,
		})
	}()

	// --------------- progress block --------------- //
	// spawn a goroutine that does the dirty job of parsing the stdout
	eventChan := make(chan string)

	// fill the channel with as many stdout line as yt-dlp produces (producer)
	go func() {
		defer cmd.Wait()
		defer r.Close()
		defer p.Complete()
		for scan.Scan() {
			eventChan <- scan.Text()
		}
	}()

	// debounce the unmarshal operation by 500ms (consumer)
	go rx.Debounce(time.Millisecond*500, eventChan, func(text string) {
		stdout := ProgressTemplate{}
		err := json.Unmarshal([]byte(text), &stdout)
		if err == nil {
			p.mem.UpdateProgress(p.id, DownloadProgress{
				Percentage: stdout.Percentage,
				Speed:      stdout.Speed,
				ETA:        stdout.Eta,
			})
		}
	})
	// ------------- end progress block ------------- //
}

// Keep process in the memoryDB but marks it as complete
// Convention: All completed processes has progress -1
// and speed 0 bps.
func (p *Process) Complete() {
	p.mem.UpdateProgress(p.id, DownloadProgress{
		Percentage: "-1",
		Speed:      0,
		ETA:        0,
	})
}

// Kill a process and remove it from the memory
func (p *Process) Kill() error {
	err := p.proc.Kill()
	p.mem.Delete(p.id)
	log.Printf("Killed process %s\n", p.id)
	return err
}
