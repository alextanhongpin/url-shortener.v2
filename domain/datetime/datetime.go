package datetime

import "time"

type DateTime struct {
	time.Time
}

func New() *DateTime {
	return &DateTime{time.Now()}
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return d.Time.MarshalJSON()
}
