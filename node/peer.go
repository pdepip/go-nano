package node

import (
    "fmt"
    "net"
)

type Peer struct {
    IP   net.IP
    Port uint16
}


// Gets Resolved, connectable, UDP Address of Peer
func (p *Peer) Addr() (addr *net.UDPAddr, err error) {
    addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.IP.String(), p.Port))
    return
}
