package rx

import "time"

/*
	Package rx contains:
	-	Definitions for common reactive programming functions/patterns
*/

// ReactiveX inspired debounce function.
//
// Debounce emits a string from the source channel only after a particular
// time span determined a Go Interval
// --A--B--CD--EFG-------|>
//	-t->                 |>
//	       -t->          |>   t is a timer tick
//	             -t->    |>
// --A-----C-----G-------|>
func Debounce(interval time.Duration, source chan string, cb func(emit string)) {
	var item string
	timer := time.NewTimer(interval)
	for {
		select {
		case item = <-source:
			timer.Reset(interval)
		case <-timer.C:
			if item != "" {
				cb(item)
			}
		}
	}
}
