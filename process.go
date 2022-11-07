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

const template = `download:
{
	"title":"%(info.title)s", 
	"thumbnail":"%(info.thumbnail)s", 
	"eta":%(progress.eta)s, 
	"resolution":"%(info.resolution)s", 
	"percentage":"%(progress._percent_str)s",
	"speed":%(progress.speed)s, 
	"size":"%(progress._total_bytes_str)s"
}`

type ProgressTemplate struct {
	Resolution string  `json:"resolution"`
	Percentage string  `json:"percentage"`
	Thumbnail  string  `json:"thumbnail"`
	Speed      float32 `json:"speed"`
	Size       string  `json:"size"`
	Eta        int     `json:"eta"`
}

type Process struct {
	mem      *MemoryDB
	url      string
	params   []string
	pid      int
	progress Progress
}

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
	p.mem.Set(p.pid, p)

	go func() {
		defer cmd.Wait()
		defer r.Close()
		defer p.Kill()
		for scan.Scan() {
			stdout := ProgressTemplate{}
			err := json.Unmarshal([]byte(scan.Text()), &stdout)
			if err == nil {
				p.mem.Update(p.pid, Progress{
					Percentage: stdout.Percentage,
					Thumbnail:  stdout.Thumbnail,
					Speed:      stdout.Speed,
					Size:       stdout.Size,
					PID:        p.pid,
					URL:        p.url,
				})
			}
		}
	}()
}

func (p *Process) Kill() error {
	cmd := exec.Command("kill", []string{string(rune(p.pid))}...)
	cmd.Start()
	err := cmd.Wait()
	p.mem.Delete(p.pid)
	log.Printf("Killed process %d\n", p.pid)
	return err
}
