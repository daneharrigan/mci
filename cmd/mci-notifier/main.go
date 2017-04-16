package main

import (
	"log"
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

	users := new(db.Users)
	for user := range users.All() {
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
	var startedAt time.Time
	endedAt := time.Now()
	comics, err := user.Comics(startedAt, endedAt)
	if err != nil {
		log.Printf("fn=Comics error=%q", err)
		return
	}

	if len(comics) == 0 {
		return
	}

	n := notifier.New(user, comics)
	log.Printf("fn=Send email=%q", user.Email)
	if err := n.Send(); err != nil {
		log.Printf("fn=Send error=%q", err)
		return
	}
}
