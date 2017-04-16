package scanner_test

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/daneharrigan/mci/config"
	"github.com/daneharrigan/mci/scanner"
)

func TestScan(t *testing.T) {
	m := &TestServer{
		i: 0,
		t: t,
		f: []string{"comics-100-000.json", "comics-100-100.json"},
	}
	srv := httptest.NewServer(m)
	defer srv.Close()

	f, err := os.Open("../data/comics-results.json")
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()
	var want []scanner.Result
	if err := json.NewDecoder(f).Decode(&want); err != nil {
		t.Fatal(err)
	}

	cfg := config.New()
	s := scanner.New(srv.URL, cfg.PublicKey, cfg.PrivateKey)
	var got []scanner.Result
	for r := range s.Scan() {
		if !contains(got, r) {
			got = append(got, r)
		}
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %d results; wanted %d results", len(got), len(want))
	}
}

type TestServer struct {
	t *testing.T
	f []string
	i int
}

func (s *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(s.f) <= s.i {
		w.WriteHeader(500)
		return
	}

	fileName := s.f[s.i]
	f, err := os.Open(filepath.Join("..", "data", fileName))
	if err != nil {
		s.t.Log(err)
		w.WriteHeader(500)
		return
	}

	defer f.Close()
	if _, err := bufio.NewReader(f).WriteTo(w); err != nil {
		s.t.Log(err)
	}

	s.i++
}

func contains(got []scanner.Result, r scanner.Result) bool {
	for _, g := range got {
		if g.Title == r.Title {
			return true
		}
	}

	return false
}
