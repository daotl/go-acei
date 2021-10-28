# go-acei 

Go implementation of [Application Consensus Engine Interface (ACEI)](https://github.com/daotl/acei).

## Generate Go code from [ACEI Protocol Buffers definitions](https://github.com/daotl/acei/tree/master/proto).

In [acei/proto](https://github.com/daotl/acei) directory, run:
```shell
protoc --gogofaster_out=. --go-grpc_out=. \
  -I=${GOPATH}/pkg/mod/google.golang.org/protobuf@v1.27.1/types/known/emptypb \ 
  -I=${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.2 -I=. ./daotl/acei/*.proto
```

## License

Apache 2.0

Copyright for portions of this fork are held by Tendermint as part of the original
[Tendermint Core](https://github.com/tendermint/tendermint) project. All other
copyright for this fork are held by DAOT Labs. All rights reserved.
