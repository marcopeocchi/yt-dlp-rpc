package main

import (
	"log"
)

type Service int

type NoArgs struct{}
type ProcArgs struct {
	PID    int
	URL    string
	Params []string
}

func (t *Service) Exec(args ProcArgs, result *string) error {
	log.Printf("Spawning new service for %s\n", args.URL)
	p := Process{mem: &db, url: args.URL, params: args.Params}
	p.Start()
	*result = args.URL
	return nil
}

func (t *Service) Result(args ProcArgs, progress *Progress) error {
	*progress = db.Get(args.PID).progress
	return nil
}

func (t *Service) Pending(args NoArgs, pending *Pending) error {
	*pending = Pending(db.Keys())
	return nil
}

func (t *Service) Running(args NoArgs, running *Running) error {
	*running = db.All()
	return nil
}

func (t *Service) Kill(args int, killed *string) error {
	proc := db.Get(args)
	var err error
	if proc != nil {
		err = proc.Kill()
	}
	return err
}
