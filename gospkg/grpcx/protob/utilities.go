package protob

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func ToProtoTime(t time.Time) *timestamp.Timestamp {
	if t.IsZero() {
		return nil
	}
	ts := &timestamp.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.UnixNano() - t.Unix()*int64(time.Second)),
	}

	return ts
}

func ToTime(ts *timestamp.Timestamp) time.Time {
	if ts == nil {
		return time.Unix(0, 0)
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos))
}
