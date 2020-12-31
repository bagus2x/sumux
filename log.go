package sumux

import (
	"log"
	"time"
)

func simpleLog(method, path string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s: %s -> %v\n", method, path, time.Since(start))
	}
}
