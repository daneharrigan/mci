package db

func (m *Comic) FindBy(name, value string) error {
	return findBy(m, name, value)
}

func (m *Comic) Create() error {
	return create(m)
}

func (m *Comic) Update() error {
	return update(m)
}

func (m *Comic) Destroy() error {
	return destroy(m)
}

func (m *Comic) Series() *Series {
	series := new(Series)
	series.FindBy("id", m.SeriesID)
	return series
}

// for Record interface
func (m *Comic) Values() []interface{} {
	return []interface{}{&m.ID, &m.SeriesID, &m.Name, &m.Thumbnail, &m.URL, &m.ReleasedAt, &m.CreatedAt, &m.UpdatedAt}
}

func (m *Comic) Columns() []string {
	return []string{"id", "series_id", "name", "thumbnail", "url", "released_at", "created_at", "updated_at"}
}

func (m *Comic) TableName() string {
	return "comics"
}

func (m *Comic) IdentifierName() string {
	return "id"
}

func (m *Comic) IdentifierValue() string {
	return m.ID
}

func (m *Comic) Touch() {
	touch(m)
}
