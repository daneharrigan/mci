package db_test

import (
	"testing"
	"time"

	"github.com/daneharrigan/mci/db"
)

func TestAll(t *testing.T) {
	t.Log("records not found")
	user := new(db.User)
	comic := new(db.Comic)
	series := new(db.Series)

	errs := []error{
		user.FindBy("email", "boom"),
		series.FindBy("name", "boom"),
		comic.FindBy("name", "boom"),
	}

	for i, err := range errs {
		if !db.IsNotFound(err) {
			t.Errorf("FindBy (%d): %q", i, err)
		}
	}

	t.Log("create records")
	user = &db.User{Email: "me@example.com"}
	series = &db.Series{Name: "example"}
	errs = []error{
		user.Create(),
		series.Create(),
	}
	comic = &db.Comic{Name: "example", SeriesID: series.ID, ReleasedAt: time.Now()}
	errs = append(errs, comic.Create())

	for i, err := range errs {
		if err != nil {
			t.Errorf("Create (%d): %q", i, err)
		}
	}

	t.Log("find records")
	user = new(db.User)
	comic = new(db.Comic)
	series = new(db.Series)

	errs = []error{
		user.FindBy("email", "me@example.com"),
		series.FindBy("name", "example"),
		comic.FindBy("name", "example"),
	}

	for i, err := range errs {
		if err != nil {
			t.Errorf("FindBy (%d): %q", i, err)
		}
	}

	t.Log("verify timestamps")
	timestamps := []time.Time{
		user.CreatedAt,
		user.UpdatedAt,
		series.CreatedAt,
		series.UpdatedAt,
		comic.CreatedAt,
		comic.UpdatedAt,
		comic.ReleasedAt,
	}

	for i, ts := range timestamps {
		if ts.IsZero() {
			t.Errorf("timestamp (%d): %v", i, t)
		}
	}

	t.Log("verify ids")
	ids := []string{user.ID, series.ID, comic.ID}
	for i, id := range ids {
		if id == "" {
			t.Errorf("id (%d): %q", i, id)
		}
	}

	t.Log("verify values")
	values := [][]string{
		[]string{"me@example.com", user.Email},
		[]string{"example", series.Name},
		[]string{"example", comic.Name},
		[]string{series.ID, comic.SeriesID},
	}

	for i, v := range values {
		if v[0] != v[1] {
			t.Errorf("values (%d), got %q; wanted %q", i, v[1], v[0])
		}
	}

	t.Log("destroy records")
	errs = []error{
		user.Destroy(),
		comic.Destroy(),
		series.Destroy(),
	}

	for i, err := range errs {
		if err != nil {
			t.Errorf("Destroy (%d): %q", i, err)
		}
	}

	t.Log("verify destroyed")
	errs = []error{
		user.FindBy("email", "me@example.com"),
		series.FindBy("name", "example"),
		comic.FindBy("name", "example"),
	}

	for i, err := range errs {
		if !db.IsNotFound(err) {
			t.Errorf("FindBy (%d): %q", i, err)
		}
	}
}
