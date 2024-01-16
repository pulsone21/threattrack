package viewutils

import (
	"time"
)

func ToTimestamp(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}
