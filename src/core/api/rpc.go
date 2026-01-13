package api

import "net/rpc"

type RPCServer struct {
	// This is the real implementation
	Impl UserStructureData
}

type RPCClient struct {
	client *rpc.Client
}
