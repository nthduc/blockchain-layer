package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	peers     map[NetAddr]*LocalTransport
	lock      sync.RWMutex
}

func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*LocalTransport),
	}
}
func (tr *LocalTransport) Addr() NetAddr {
	return tr.addr
}

func (tr *LocalTransport) Consume() <-chan RPC {
	return tr.consumeCh
}

func (tr *LocalTransport) Connect(trb Transport) error {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	tr.peers[trb.Addr()] = trb.(*LocalTransport)
	return nil
}

func (tr *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	if tr.addr == to {
		return fmt.Errorf("could not send message to yourself")
	}

	peer, ok := tr.peers[to]
	if !ok {
		return fmt.Errorf("%s: could mpt send message to %s", tr.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    tr.addr,
		Payload: payload,
	}

	return nil
}