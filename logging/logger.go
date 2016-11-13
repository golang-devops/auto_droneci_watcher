package logging

import (
	"github.com/go-zero-boilerplate/extended-apex-logger/logging"
)

//Logger just wraps the "external" logger as a "local" Logger
type Logger interface {
	logging.Logger
}
