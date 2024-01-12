package templates

import (
	"fmt"
	"time"
)

func toTimestamp(epoch uint) string {
	return fmt.Sprint(time.Unix(int64(epoch),0))
}