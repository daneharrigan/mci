package main

import (
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/daneharrigan/mci/config"
	"github.com/daneharrigan/mci/db"
	"github.com/daneharrigan/mci/scanner"
)

var (
	cfg = config.New()
	wg  sync.WaitGroup
)

func init() {
	log.SetPrefix("ns=mci-scanner ")
	log.SetFlags(log.Ltime | log.LUTC)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	ch := make(chan scanner.Result)
	defer close(ch)

	for i := 0; i < runtime.NumCPU(); i++ {
		go receiver(ch)
	}

	s := scanner.New(cfg.URL, cfg.PublicKey, cfg.PrivateKey)
	for r := range s.Scan() {
		wg.Add(1)
		ch <- r
	}

	wg.Wait()

	if s.Err() != nil {
		log.Printf("fn=Scan error=%q", s.Err())
		os.Exit(1)
	}
}

func receiver(ch <-chan scanner.Result) {
	for r := range ch {
		handle(r)
	}
}

func handle(r scanner.Result) {
	defer wg.Done()
	comic := new(db.Comic)
	if err := comic.FindBy("name", r.Title); err == nil {
		return
	} else if !db.IsNotFound(err) {
		log.Printf("fn=FindBy error=%q", err)
		return
	}

	series := new(db.Series)
	if err := series.FindBy("name", r.Series.Name); err != nil {
		if !db.IsNotFound(err) {
			log.Printf("fn=FindBy error=%q", err)
			return
		}

		series.Name = r.Series.Name
		if err := series.Create(); err != nil {
			log.Printf("fn=Create error=%q", err)
			return
		}
	}

	thumbnail := r.Thumbnail.Path + "." + r.Thumbnail.Extension
	releasedAt, err := scanner.FindUnlimitedDate(r.Dates)
	if err != nil {
		log.Printf("fn=FindUnlimitedDate error=%q", err)
		return
	}

	comic.Name = r.Title
	comic.SeriesID = series.ID
	comic.Thumbnail = thumbnail
	comic.ReleasedAt = releasedAt

	if err := comic.Create(); err != nil {
		log.Printf("fn=Create error=%q", err)
	}
}
