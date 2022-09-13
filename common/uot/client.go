package uot

import (
	"encoding/binary"
	"io"
	"net"
	"sync"

	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/buf"
	M "github.com/sagernet/sing/common/metadata"
)

type ClientConn struct {
	net.Conn
	readAccess  sync.Mutex
	writeAccess sync.Mutex
}

func NewClientConn(conn net.Conn) *ClientConn {
	return &ClientConn{
		Conn: conn,
	}
}

func (c *ClientConn) ReadPacket(buffer *buf.Buffer) (M.Socksaddr, error) {
	c.readAccess.Lock()
	defer c.readAccess.Unlock()

	destination, err := AddrParser.ReadAddrPort(c)
	if err != nil {
		return M.Socksaddr{}, err
	}
	var length uint16
	err = binary.Read(c, binary.BigEndian, &length)
	if err != nil {
		return M.Socksaddr{}, err
	}
	if buffer.FreeLen() < int(length) {
		return M.Socksaddr{}, io.ErrShortBuffer
	}
	return destination, common.Error(buffer.ReadFullFrom(c, int(length)))
}

func (c *ClientConn) WritePacket(buffer *buf.Buffer, destination M.Socksaddr) error {
	c.writeAccess.Lock()
	defer c.writeAccess.Unlock()

	defer buffer.Release()
	err := AddrParser.WriteAddrPort(c, destination)
	if err != nil {
		return err
	}
	err = binary.Write(c, binary.BigEndian, uint16(buffer.Len()))
	if err != nil {
		return err
	}
	return common.Error(c.Write(buffer.Bytes()))
}

func (c *ClientConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	c.readAccess.Lock()
	defer c.readAccess.Unlock()

	addrPort, err := AddrParser.ReadAddrPort(c)
	if err != nil {
		return 0, nil, err
	}
	var length uint16
	err = binary.Read(c, binary.BigEndian, &length)
	if err != nil {
		return 0, nil, err
	}
	if len(p) < int(length) {
		return 0, nil, io.ErrShortBuffer
	}
	n, err = io.ReadFull(c, p[:length])
	if err != nil {
		return 0, nil, err
	}
	addr = addrPort.UDPAddr()
	return
}

func (c *ClientConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	c.writeAccess.Lock()
	defer c.writeAccess.Unlock()

	err = AddrParser.WriteAddrPort(c, M.SocksaddrFromNet(addr))
	if err != nil {
		return
	}
	err = binary.Write(c, binary.BigEndian, uint16(len(p)))
	if err != nil {
		return
	}
	return c.Write(p)
}
