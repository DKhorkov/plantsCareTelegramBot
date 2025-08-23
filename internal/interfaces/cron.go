package interfaces

type Callback = func() error

//go:generate mockgen -source=cron.go -destination=../../mocks/cron/cron.go -package=mockcron
type Cron interface {
	Run() error
	Stop() error
}

//go:generate mockgen -source=cron.go -destination=../../mocks/cron/cron.go -package=mockcron
type Preparer interface {
	GetCallback() Callback
}
