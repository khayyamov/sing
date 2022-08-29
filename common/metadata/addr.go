package metadata

import (
	"net"
	"net/netip"
	"strconv"
	"unsafe"
)

type Socksaddr struct {
	Addr netip.Addr
	Port uint16
	Fqdn string
}

func (ap Socksaddr) Network() string {
	return "socks"
}

func (ap Socksaddr) IsIP() bool {
	return ap.Addr.IsValid() && ap.Fqdn == ""
}

func (ap Socksaddr) IsIPv4() bool {
	return ap.Addr.Is4() || ap.Addr.Is4In6()
}

func (ap Socksaddr) IsIPv6() bool {
	return ap.Addr.Is6() && !ap.Addr.Is4In6()
}

func (ap Socksaddr) Unwrap() Socksaddr {
	if ap.Addr.Is4In6() {
		return Socksaddr{
			Addr: netip.AddrFrom4(ap.Addr.As4()),
			Port: ap.Port,
		}
	}
	return ap
}

func (ap Socksaddr) IsFqdn() bool {
	return !ap.Addr.IsValid() && ap.Fqdn != ""
}

func (ap Socksaddr) IsValid() bool {
	return ap.IsIP() || ap.IsFqdn()
}

func (ap Socksaddr) AddrString() string {
	if ap.Addr.IsValid() {
		return ap.Addr.String()
	} else {
		return ap.Fqdn
	}
}

func (ap Socksaddr) IPAddr() *net.IPAddr {
	return &net.IPAddr{
		IP: ap.Addr.AsSlice(),
	}
}

func (ap Socksaddr) TCPAddr() *net.TCPAddr {
	return &net.TCPAddr{
		IP:   ap.Addr.AsSlice(),
		Port: int(ap.Port),
	}
}

func (ap Socksaddr) UDPAddr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   ap.Addr.AsSlice(),
		Port: int(ap.Port),
	}
}

func (ap Socksaddr) AddrPort() netip.AddrPort {
	return *(*netip.AddrPort)(unsafe.Pointer(&ap))
}

func (ap Socksaddr) String() string {
	return net.JoinHostPort(ap.AddrString(), strconv.Itoa(int(ap.Port)))
}

func TCPAddr(ap netip.AddrPort) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   ap.Addr().AsSlice(),
		Port: int(ap.Port()),
	}
}

func UDPAddr(ap netip.AddrPort) *net.UDPAddr {
	return &net.UDPAddr{
		IP:   ap.Addr().AsSlice(),
		Port: int(ap.Port()),
	}
}

func AddrPortFrom(ip net.IP, port uint16) netip.AddrPort {
	return netip.AddrPortFrom(AddrFromIP(ip), port)
}

func SocksaddrFrom(addr netip.Addr, port uint16) Socksaddr {
	return SocksaddrFromNetIP(netip.AddrPortFrom(addr, port))
}

func SocksaddrFromNetIP(ap netip.AddrPort) Socksaddr {
	return Socksaddr{
		Addr: ap.Addr(),
		Port: ap.Port(),
	}
}

func SocksaddrFromNet(ap net.Addr) Socksaddr {
	if ap == nil {
		return Socksaddr{}
	}
	if socksAddr, ok := ap.(Socksaddr); ok {
		return socksAddr
	}
	addr := SocksaddrFromNetIP(AddrPortFromNet(ap))
	if addr.IsValid() {
		return addr
	}
	return ParseSocksaddr(ap.String())
}

func AddrFromNetAddr(netAddr net.Addr) netip.Addr {
	if addr := AddrPortFromNet(netAddr); addr.Addr().IsValid() {
		return addr.Addr()
	}
	switch addr := netAddr.(type) {
	case Socksaddr:
		return addr.Addr
	case *net.IPAddr:
		return AddrFromIP(addr.IP)
	case *net.IPNet:
		return AddrFromIP(addr.IP)
	default:
		return netip.Addr{}
	}
}

func AddrPortFromNet(netAddr net.Addr) netip.AddrPort {
	var ip net.IP
	var port uint16
	switch addr := netAddr.(type) {
	case Socksaddr:
		return addr.AddrPort()
	case *net.TCPAddr:
		ip = addr.IP
		port = uint16(addr.Port)
	case *net.UDPAddr:
		ip = addr.IP
		port = uint16(addr.Port)
	case *net.IPAddr:
		ip = addr.IP
	}
	return netip.AddrPortFrom(AddrFromIP(ip), port)
}

func AddrFromIP(ip net.IP) netip.Addr {
	addr, _ := netip.AddrFromSlice(ip)
	return addr
}

func ParseAddr(s string) netip.Addr {
	addr, _ := netip.ParseAddr(s)
	return addr
}

func ParseSocksaddr(address string) Socksaddr {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return Socksaddr{}
	}
	return ParseSocksaddrHostPortStr(host, port)
}

func ParseSocksaddrHostPort(host string, port uint16) Socksaddr {
	netAddr, err := netip.ParseAddr(host)
	if err != nil {
		return Socksaddr{
			Fqdn: host,
			Port: port,
		}
	} else {
		return Socksaddr{
			Addr: netAddr,
			Port: port,
		}
	}
}

func ParseSocksaddrHostPortStr(host string, portStr string) Socksaddr {
	port, _ := strconv.Atoi(portStr)
	netAddr, err := netip.ParseAddr(host)
	if err != nil {
		return Socksaddr{
			Fqdn: host,
			Port: uint16(port),
		}
	} else {
		return Socksaddr{
			Addr: netAddr,
			Port: uint16(port),
		}
	}
}
