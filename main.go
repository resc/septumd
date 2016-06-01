package main

import (
	"log"
	"time"

	"github.com/ReSc/septum"
)

var (
	config struct {
		dbPath      string
		servicePort int
		apiPath     string
		webPath     string
	}
)

func main() {

	configure()

	repo := newRepository()

	err := repo.Open(config.dbPath)
	log.Println("Opened database")
	defer func() {

		log.Println("Closing database")
		repo.Close()
	}()

	if err != nil {
		log.Fatalf("Failed to open event database '%s': %s", config.dbPath, err)
	}
	nexus := septum.NewNexus(septum.SystemEnvironment())
	now := time.Now()
	for i := 0; i < 100000; i += 1 {
		if i%100 == 0 {
			nextNow := time.Now()
			d := nextNow.Sub(now)
			now = nextNow
			log.Printf("%s %d\n", d, i)
		}
		e := septum.EventData{
			Timeline:  1,
			Timestamp: time.Now().UTC(),
			Kind:      "Test",
			Data:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		}
		var event *septum.Event
		if isTimer(e.Kind) {
			event = nexus.AddTimer(e.Timeline, e.Timestamp, e.Kind, e.Data)
			if err := repo.SaveTimer(event); err != nil {
				log.Fatal(err)
			}
		} else {
			event = nexus.AddEvent(e.Timeline, e.Timestamp, e.Kind, e.Data)
		}
		if err := repo.SaveEvent(event); err != nil {
			log.Fatal(err)
		}
	}

}

func isTimer(eventKind string) bool {
	return false
}
