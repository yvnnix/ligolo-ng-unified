// Ligolo-ng
// Copyright (C) 2025 Nicolas Chatelain (nicocha30)

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package proxy

import (
	"context"
	"github.com/hashicorp/yamux"
	"github.com/nicocha30/ligolo-ng/pkg/proxy/netstack"
	"github.com/sirupsen/logrus"
)

const (
	MaxConnectionHandler = 4096
)

type LigoloTunnel struct {
	nstack *netstack.NetStack
}

func NewLigoloTunnel(stackSettings netstack.StackSettings) (*LigoloTunnel, error) {
	// Create a new stack, but without connPool.
	// The connPool will be created when using the *start* command
	nstack, err := netstack.NewStack(stackSettings, nil)
	if err != nil {
		return nil, err
	}
	return &LigoloTunnel{nstack: nstack}, nil
}

func (t *LigoloTunnel) HandleSession(session *yamux.Session, ctx context.Context) {

	// Create a new, empty, connpool to store connections/packets
	connPool := netstack.NewConnPool(MaxConnectionHandler)
	t.nstack.SetConnPool(&connPool)

	// Cleanup pool if channel is closed
	defer connPool.Close()

	for {
		select {
		case <-ctx.Done():
			t.Close()
			return
		case <-connPool.CloseChan: // pool closed, we can't process packets!
			logrus.Infof("Connection pool closed")
			t.Close()
			return
		case tunnelPacket := <-connPool.Pool: // Process connections/packets
			go netstack.HandlePacket(t.nstack.GetStack(), tunnelPacket, session)
		}
	}
}

func (t *LigoloTunnel) GetStack() *netstack.NetStack {
	return t.nstack
}

func (t *LigoloTunnel) Close() {
	t.nstack.Close()
}
