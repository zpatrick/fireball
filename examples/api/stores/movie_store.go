package stores

import (
	"encoding/json"
	"github.com/zpatrick/fireball/examples/api/models"
	"github.com/zpatrick/go-sdata/container"
)

type MovieStore struct {
	container container.Container
	table     string
}

func NewMovieStore(container container.Container) *MovieStore {
	return &MovieStore{
		container: container,
		table:     "models.Movie",
	}
}

func (this *MovieStore) Init() error {
	return this.container.Init(this.table)
}

type MovieStoreInsert struct {
	*MovieStore
	data *models.Movie
}

func (this *MovieStore) Insert(data *models.Movie) *MovieStoreInsert {
	return &MovieStoreInsert{
		MovieStore: this,
		data:       data,
	}
}

func (this *MovieStoreInsert) Execute() error {
	bytes, err := json.Marshal(this.data)
	if err != nil {
		return err
	}

	return this.container.Insert(this.table, this.data.ID, bytes)
}

type MovieStoreSelect struct {
	*MovieStore
	query  string
	filter MovieFilter
	all    bool
}

func (this *MovieStore) Select(query string) *MovieStoreSelect {
	return &MovieStoreSelect{
		MovieStore: this,
		query:      query,
	}
}

func (this *MovieStore) SelectAll() *MovieStoreSelect {
	return &MovieStoreSelect{
		MovieStore: this,
		all:        true,
	}
}

type MovieFilter func(*models.Movie) bool

func (this *MovieStoreSelect) Where(filter MovieFilter) *MovieStoreSelect {
	this.filter = filter
	return this
}

func (this *MovieStoreSelect) Execute() ([]*models.Movie, error) {
	var query func() (map[string][]byte, error)

	if this.all {
		query = func() (map[string][]byte, error) { return this.container.SelectAll(this.table) }
	} else {
		query = func() (map[string][]byte, error) { return this.container.Select(this.table, this.query) }
	}

	data, err := query()
	if err != nil {
		return nil, err
	}

	results := []*models.Movie{}
	for _, d := range data {
		var value *models.Movie

		if err := json.Unmarshal(d, &value); err != nil {
			return nil, err
		}

		if this.filter == nil || this.filter(value) {
			results = append(results, value)
		}
	}

	return results, nil
}

type MovieStoreSelectFirst struct {
	*MovieStoreSelect
}

func (this *MovieStoreSelect) FirstOrNil() *MovieStoreSelectFirst {
	return &MovieStoreSelectFirst{
		MovieStoreSelect: this,
	}
}

func (this *MovieStoreSelectFirst) Execute() (*models.Movie, error) {
	results, err := this.MovieStoreSelect.Execute()
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results[0], nil
	}

	return nil, nil
}

type MovieStoreDelete struct {
	*MovieStore
	key string
}

func (this *MovieStore) Delete(key string) *MovieStoreDelete {
	return &MovieStoreDelete{
		MovieStore: this,
		key:        key,
	}
}

func (this *MovieStoreDelete) Execute() (bool, error) {
	return this.container.Delete(this.table, this.key)
}
