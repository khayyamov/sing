package socks

import (
	"net"

	"github.com/khayyamov/sing/common"
	"github.com/khayyamov/sing/common/buf"
	"github.com/khayyamov/sing/common/bufio"
	M "github.com/khayyamov/sing/common/metadata"
	N "github.com/khayyamov/sing/common/network"
)

var _ N.VectorisedPacketWriter = (*VectorisedAssociatePacketConn)(nil)

type VectorisedAssociatePacketConn struct {
	AssociatePacketConn
	N.VectorisedPacketWriter
}

func NewVectorisedAssociatePacketConn(conn net.PacketConn, writer N.VectorisedPacketWriter, remoteAddr M.Socksaddr, underlying net.Conn) *VectorisedAssociatePacketConn {
	return &VectorisedAssociatePacketConn{
		AssociatePacketConn{
			NetPacketConn: bufio.NewPacketConn(conn),
			remoteAddr:    remoteAddr,
			underlying:    underlying,
		},
		writer,
	}
}

func NewVectorisedAssociateConn(conn net.Conn, writer N.VectorisedWriter, remoteAddr M.Socksaddr, underlying net.Conn) *VectorisedAssociatePacketConn {
	return &VectorisedAssociatePacketConn{
		AssociatePacketConn{
			NetPacketConn: bufio.NewUnbindPacketConn(conn),
			remoteAddr:    remoteAddr,
			underlying:    underlying,
		},
		&bufio.UnbindVectorisedPacketWriter{VectorisedWriter: writer},
	}
}

func (v *VectorisedAssociatePacketConn) WriteVectorisedPacket(buffers []*buf.Buffer, destination M.Socksaddr) error {
	header := buf.NewSize(3 + M.SocksaddrSerializer.AddrPortLen(destination))
	defer header.Release()
	common.Must(header.WriteZeroN(3))
	err := M.SocksaddrSerializer.WriteAddrPort(header, destination)
	if err != nil {
		return err
	}
	return v.VectorisedPacketWriter.WriteVectorisedPacket(append([]*buf.Buffer{header}, buffers...), destination)
}

func (c *VectorisedAssociatePacketConn) FrontHeadroom() int {
	return 0
}
