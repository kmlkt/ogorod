package ogorod

type ServiceState int8

const (
	Missing ServiceState = iota
	Idle
	Running
)

type Service interface {
	Key() string
	State() ServiceState
	Install() error
	Start() error
	Stop() error
}
