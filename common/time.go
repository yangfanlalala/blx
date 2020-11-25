package common

import (
	"fmt"
	"time"
)

const (
	ISO8601Format     = "2006-01-02T15:04:05Z"
	ISO8601DateFormat = "2006-01-02"
)

func FormatISO8601Date(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(ISO8601Format)
}

type ISO8601Date time.Time

func (d ISO8601Date) MarshalJSON() ([]byte, error) {
	format := fmt.Sprintf("\"%s\"", time.Time(d).Format(ISO8601DateFormat))
	return []byte(format), nil
}

func (d *ISO8601Date) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return
	}
	t, err := time.Parse(`"`+ISO8601DateFormat+`"`, string(data))
	*d = ISO8601Date(t)
	return
}

func (d *ISO8601Date) String() string {
	return time.Time(*d).Format(ISO8601DateFormat)
}

type ISO8601Datetime time.Time

func (d ISO8601Datetime) MarshalJSON() ([]byte, error) {
	format := fmt.Sprintf("\"%s\"", time.Time(d).Format(ISO8601Format))
	return []byte(format), nil
}

func (d *ISO8601Datetime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return
	}
	t, err := time.Parse(`"`+ISO8601Format+`"`, string(data))
	*d = ISO8601Datetime(t)
	return
}
