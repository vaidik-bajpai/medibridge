package models

import (
	"strings"
	"time"
)

type DateOnly time.Time

const layoutDateOnly = "2006-01-02"

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	t, err := time.Parse(layoutDateOnly, s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d) // convert to time.Time to call Format
	return []byte(`"` + t.Format(layoutDateOnly) + `"`), nil
}
