/*
Package server is used to start a new ABCI server.

It contains two server implementation:
 * gRPC server
 * socket server

*/
package server

import (
	"fmt"

	"github.com/daotl/go-log/v2"
	ssrv "github.com/daotl/guts/service/suture"

	"github.com/daotl/go-acei/types"
)

func NewServer(protoAddr, transport string, app types.Application, logger log.StandardLogger,
) (ssrv.Service, error) {
	var s ssrv.Service
	var err error
	switch transport {
	case "socket":
		s, err = NewSocketServer(protoAddr, app, logger)
	case "grpc":
		s, err = NewGRPCServer(protoAddr, types.NewGRPCApplication(app), logger)
	default:
		err = fmt.Errorf("unknown server type %s", transport)
	}
	return s, err
}
