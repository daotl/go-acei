package aceiclient

import (
	"fmt"

	"github.com/daotl/go-log/v2"
	gsync "github.com/daotl/guts/sync"

	"github.com/daotl/go-acei/types"
)

// LoggerCreator creates new loggers.
type LoggerCreator func() (log.StandardLogger, error)

// Creator creates new ACEI clients.
type Creator func() (Client, error)

// NewLocalCreator returns a Creator for the given app,
// which will be running locally.
func NewLocalCreator(app types.Application, loggerCreator LoggerCreator) Creator {
	mtx := new(gsync.Mutex)

	return func() (Client, error) {
		var logger log.StandardLogger = nil
		if loggerCreator != nil {
			var err error
			if logger, err = loggerCreator(); err != nil {
				return nil, err
			}
		}
		return NewLocalClient(mtx, app, logger)
	}
}

// NewRemoteCreator returns a Creator for the given address (e.g.
// "192.168.0.1") and transport (e.g. "tcp"). Set mustConnect to true if you
// want the client to connect before reporting success.
func NewRemoteCreator(addr, transport string, mustConnect bool, loggerCreator LoggerCreator) Creator {
	return func() (Client, error) {
		var logger log.StandardLogger = nil
		if loggerCreator != nil {
			var err error
			if logger, err = loggerCreator(); err != nil {
				return nil, err
			}
		}

		remoteApp, err := NewClient(addr, transport, mustConnect, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to proxy: %w", err)
		}

		return remoteApp, nil
	}
}
