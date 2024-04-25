package main

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	rpc "github.com/libp2p/go-libp2p-gorpc"
)

func TestService_StartMessaging(t *testing.T) {
	type fields struct {
		rpcServer *rpc.Server
		rpcClient *rpc.Client
		host      host.Host
		protocol  protocol.ID
		counter   int
	}
	srcHost, _ := NewHost(context.Background(), 0, 0)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "happy case",
			fields: fields{
				host:     srcHost,
				protocol: protocol.TestingID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				rpcServer: tt.fields.rpcServer,
				rpcClient: tt.fields.rpcClient,
				host:      tt.fields.host,
				protocol:  tt.fields.protocol,
				counter:   tt.fields.counter,
			}
			errors := s.StartMessaging(tt.args.ctx)
			for _, err := range errors {
				if err != nil {
					t.Errorf("Connect() error = %v, wantErr %v", err, nil)
					return
				}
			}
		})
	}
}
