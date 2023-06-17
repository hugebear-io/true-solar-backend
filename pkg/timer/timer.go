package timer

import (
	"time"

	"github.com/hugebear-io/true-solar-backend/pkg/config"
)

func SetupTimezone() *time.Location {
	tz := config.Config.Timezone
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	time.Local = loc
	return loc
}
