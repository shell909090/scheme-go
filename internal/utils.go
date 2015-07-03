package internal

import (
	"errors"

	"github.com/op/go-logging"
)

var (
	ErrArguments = errors.New("wrong arguments")
)

var (
	log = logging.MustGetLogger("internal")
)
