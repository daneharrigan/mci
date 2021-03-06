package main

import (
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/daneharrigan/mci/db"
	"github.com/daneharrigan/mci/notifier"
)

var wg sync.WaitGroup

func init() {
	log.SetPrefix("ns=mci-scanner ")
	log.SetFlags(log.Ltime | log.LUTC)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	ch := make(chan *db.User)
	defer close(ch)

	for i := 0; i < runtime.NumCPU(); i++ {
		go receiver(ch)
	}

	user := new(db.User)
	all, err := user.All()
	if err != nil {
		log.Printf("fn=All error=%q", err)
		os.Exit(1)
	}
	for user := range all {
		log.Printf("user=%q", user.Email)
		wg.Add(1)
		ch <- user
	}

	wg.Wait()
}

func receiver(ch <-chan *db.User) {
	for user := range ch {
		handle(user)
	}
}

func handle(user *db.User) {
	defer wg.Done()

	now := time.Now()
	earliestAt := findDate(now, time.Monday, -24)
	comics, err := user.Comics(earliestAt)
	if err != nil {
		log.Printf("fn=Comics error=%q", err)
		return
	}

	log.Printf("fn=Send user=%q comics=%d earliest_at=%v", user.Email, len(comics), earliestAt)
	if len(comics) == 0 {
		return
	}

	n := notifier.New(user, comics)
	if err := n.Send(); err != nil {
		log.Printf("fn=Send error=%q", err)
		return
	}
}

func findDate(t time.Time, day time.Weekday, offset time.Duration) time.Time {
	for {
		if t.Weekday() == day {
			break
		}

		t = t.Add(offset * time.Hour)
	}

	hrs := time.Duration(-t.Hour()) * time.Hour
	return t.Truncate(time.Hour).Add(hrs)
}
