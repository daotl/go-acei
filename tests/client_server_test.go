package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	abciclientent "github.com/daotl/go-acei/client"
	"github.com/daotl/go-acei/example/kvstore"
	abciserver "github.com/daotl/go-acei/server"
)

func TestClientServerNoAddrPrefix(t *testing.T) {
	addr := "localhost:26658"
	transport := "socket"
	app := kvstore.NewApplication()

	server, err := abciserver.NewServer(addr, transport, app, nil)
	assert.NoError(t, err, "expected no error on NewServer")
	readyCh, sResCh := server.Start(context.Background())
	assert.NoError(t, <-readyCh, "expected no error on server.Start")
	go func() { assert.NoError(t, <-sResCh, "expected no error on server stopping") }()

	client, err := abciclientent.NewClient(addr, transport, true, nil)
	assert.NoError(t, err, "expected no error on NewClient")
	readyCh, cResCh := client.Start(context.Background())
	assert.NoError(t, <-readyCh, "expected no error on client.Start")
	go func() { assert.NoError(t, <-cResCh, "expected no error on client stopping") }()
}
