package scanner

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	DateType   = "unlimitedDate"
	DateLayout = "2006-01-02T15:04:05-0700"
)

type Date struct {
	Date string `json:"date"`
	Type string `json:"type"`
}

type Result struct {
	Title     string `json:"title"`
	Dates     []Date `json:"dates"`
	Thumbnail struct {
		Extension string `json:"extension"`
		Path      string `json:"path"`
	} `json:"thumbnail"`
	Series struct {
		Name string `json:"name"`
	} `json:"series"`
}

type Response struct {
	Code int `json:"code"`
	Data struct {
		Limit   int      `json:"limit"`
		Total   int      `json:"total"`
		Count   int      `json:"offset"`
		Results []Result `json:"results"`
	} `json:"data"`
}

type Scanner interface {
	Scan() <-chan Result
	Err() error
}

func New(url, publicKey, privateKey string) Scanner {
	s := &scanner{
		url:        url,
		err:        nil,
		ch:         make(chan Result),
		publicKey:  publicKey,
		privateKey: privateKey,
	}
	return s
}

type scanner struct {
	url        string
	err        error
	ch         chan Result
	publicKey  string
	privateKey string
}

func (s *scanner) Scan() <-chan Result {
	go s.fetch()
	return s.ch
}

func (s *scanner) Err() error {
	return s.err
}

func (s *scanner) fetch() {
	defer close(s.ch)
	params := map[string]string{
		"format":          "comic",
		"hasDigitalIssue": "true",
		"limit":           "100",
		"offset":          "0",
	}

	count, err := s.handle(s.get("/v1/public/comics", params))
	if err != nil {
		s.err = err
		return
	}

	if count == 0 {
		return
	}

	params["offset"] = "100"
	_, s.err = s.handle(s.get("/v1/public/comics", params))
	s.err = err
}

func (s *scanner) handle(r *Response, err error) (int, error) {
	if err != nil {
		return 0, err
	}

	if r.Code != 200 {
		return 0, NewHTTPError(r.Code)
	}

	for _, r := range r.Data.Results {
		t, err := FindUnlimitedDate(r.Dates)
		if err != nil || t.IsZero() {
			continue
		}

		s.ch <- r
	}

	count := r.Data.Total - r.Data.Count
	return count, nil
}

func (s *scanner) get(path string, params map[string]string) (*Response, error) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	h := md5.New()
	fmt.Fprintf(h, "%s%s%s", ts, s.privateKey, s.publicKey)

	params["hash"] = fmt.Sprintf("%x", h.Sum(nil))
	params["ts"] = ts
	params["apikey"] = s.publicKey

	u, err := url.Parse(s.url)
	if err != nil {
		return nil, err
	}

	q := u.Query()

	for k, v := range params {
		q.Set(k, v)
	}

	u.Path = path
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, NewHTTPError(res.StatusCode)
	}

	r := new(Response)
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return nil, err
	}

	_, err = io.Copy(ioutil.Discard, res.Body)
	return r, err
}

func FindUnlimitedDate(dates []Date) (time.Time, error) {
	for _, d := range dates {
		if d.Type != DateType {
			continue
		}

		return time.Parse(DateLayout, d.Date)
	}

	var t time.Time
	return t, nil
}

func NewHTTPError(code int) error {
	return errors.New(http.StatusText(code))
}
