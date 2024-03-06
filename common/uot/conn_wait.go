package uot

import (
	"encoding/binary"

	"github.com/khayyamov/sing/common/buf"
	E "github.com/khayyamov/sing/common/exceptions"
	M "github.com/khayyamov/sing/common/metadata"
	N "github.com/khayyamov/sing/common/network"
)

func (c *Conn) InitializeReadWaiter(options N.ReadWaitOptions) (needCopy bool) {
	c.readWaitOptions = options
	return false
}

func (c *Conn) WaitReadPacket() (buffer *buf.Buffer, destination M.Socksaddr, err error) {
	if c.isConnect {
		destination = c.destination
	} else {
		destination, err = AddrParser.ReadAddrPort(c.Conn)
		if err != nil {
			return
		}
	}
	var length uint16
	err = binary.Read(c.Conn, binary.BigEndian, &length)
	if err != nil {
		return
	}
	buffer = c.readWaitOptions.NewPacketBuffer()
	_, err = buffer.ReadFullFrom(c.Conn, int(length))
	if err != nil {
		buffer.Release()
		return nil, M.Socksaddr{}, E.Cause(err, "UoT read")
	}
	c.readWaitOptions.PostReturn(buffer)
	return
}
