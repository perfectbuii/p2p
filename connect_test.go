package main

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
)

func TestConnect(t *testing.T) {
	type args struct {
		ctx            context.Context
		host           host.Host
		bootstrapPeers []multiaddr.Multiaddr
		hostDtID       string
	}
	srcHost, _ := NewHost(context.Background(), 0, 0)
	dstHost, _ := NewHost(context.Background(), 0, 0)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "happy case",
			args: args{
				ctx:            context.Background(),
				host:           srcHost,
				bootstrapPeers: dstHost.Addrs(),
				hostDtID:       dstHost.ID().Pretty(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.args.ctx, tt.args.host, tt.args.bootstrapPeers, tt.args.hostDtID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
