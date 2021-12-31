package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/turnage/graw/reddit"
)

type savereddit struct {
	Name       string
	Attrs      Attrs
	CreatedUTC uint64
}

type Attrs reddit.Post

func (r Attrs) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(`type assertion to []byte failed`)
	}

	return json.Unmarshal(b, &r)
}

func (s savereddit) Insert() {
	_, err := dbpool.Exec(context.Background(),
		"INSERT INTO savereddit (name, attrs, created) VALUES($1, $2, $3)",
		s.Name, s.Attrs, s.CreatedUTC)
	if err != nil {
		Log.Debug().Err(err).Send()
	}
}

func (s *savereddit) Select() {
	if err := dbpool.QueryRow(context.Background(), "SELECT * FROM savereddit LIMIT 1").Scan(&s.Name, &s.Attrs, &s.CreatedUTC); err != nil {
		Log.Debug().Err(err).Send()
	}
}
