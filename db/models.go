package db

import (
	"time"
)

type User struct {
	ID        string
	Email     string
	Provider  string
	Picture   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Series struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comic struct {
	ID         string
	SeriesID   string
	Name       string
	Thumbnail  string
	URL        string
	ReleasedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// user_series (
// 	id UUID NOT NULL PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
// 	user_id UUID REFERENCES users (id),
// 	series_id UUID REFERENCES series (id),
// 	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
// 	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
// );
