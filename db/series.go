package db

func (m *Series) FindBy(name, value string) error {
	return findBy(m, name, value)
}

func (m *Series) Create() error {
	return create(m)
}

func (m *Series) Update() error {
	return update(m)
}

func (m *Series) Destroy() error {
	return destroy(m)
}

// for Record interface
func (m *Series) Values() []interface{} {
	return []interface{}{&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt}
}

func (m *Series) Columns() []string {
	return []string{"id", "name", "created_at", "updated_at"}
}

func (m *Series) TableName() string {
	return "series"
}

func (m *Series) IdentifierName() string {
	return "id"
}

func (m *Series) IdentifierValue() string {
	return m.ID
}

func (m *Series) Touch() {
	touch(m)
}
