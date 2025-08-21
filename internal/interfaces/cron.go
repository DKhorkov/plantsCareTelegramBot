package interfaces

import "time"

type Cron interface {
	Run(limit, offset int, interval time.Duration) error
	Stop() error
}
