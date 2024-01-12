package templates

import (
	"time"
)

func toTimestamp(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}
