package clientUtil

import (
	"context"
	"net"
	"net/smtp"
	"time"
)

// DialTimeout. call Go net.DialTimeout underlying but with the retry feature based on the passed in reqInfo parameter.
func DialTimeout(network, address string, reqInfo *RequestInfo) (net.Conn, error) {
	conn, _, err := dialerCommon(nil, nil, network, address, "", reqInfo, 0, false)
	if conn != nil {
		return conn.(net.Conn), err
	} else {
		return nil, err
	}
}

// Dial. call Go net.Dial underlying but with the retry feature based on the passed in reqInfo parameter.
func Dial(dialer *net.Dialer, network, address string, reqInfo *RequestInfo) (net.Conn, error) {
	dialer.Timeout = time.Duration(reqInfo.TimeoutSec) * time.Second
	conn, _, err := dialerCommon(dialer, nil, network, address, "", reqInfo, 1, false)
	if conn != nil {
		return conn.(net.Conn), err
	} else {
		return nil, err
	}
}

// DialContext. call Go net.DialContext underlying but with the retry feature based on the passed in reqInfo parameter.
func DialContext(dialer *net.Dialer, network, address string, reqInfo *RequestInfo) (net.Conn, error) {
	conn, _, err := dialerCommon(dialer, nil, network, address, "", reqInfo, 2, true)
	if conn != nil {
		return conn.(net.Conn), err
	} else {
		return nil, err
	}
}

// DialIP. call Go net.DialContext underlying but with the retry feature based on the passed in reqInfo parameter.
func DialIP(network, address string, reqInfo *RequestInfo) (*net.IPConn, error) {
	conn, _, err := dialerCommon(nil, nil, network, address, "", reqInfo, 0, false)
	if conn != nil {
		return conn.(*net.IPConn), err
	} else {
		return nil, err
	}
}

// DialTCP. call Go net.DialContext underlying but with the retry feature based on the passed in reqInfo parameter.
func DialTCP(network, address string, reqInfo *RequestInfo) (*net.TCPConn, error) {
	conn, _, err := dialerCommon(nil, nil, network, address, "", reqInfo, 0, false)
	if conn != nil {
		return conn.(*net.TCPConn), err
	} else {
		return nil, err
	}
}

// DialUDP. call Go net.DialContext underlying but with the retry feature based on the passed in reqInfo parameter.
func DialUDP(network, address string, reqInfo *RequestInfo) (*net.UDPConn, error) {
	conn, _, err := dialerCommon(nil, nil, network, address, "", reqInfo, 0, false)
	if conn != nil {
		return conn.(*net.UDPConn), err
	} else {
		return nil, err
	}
}

// DialUnix. call Go net.DialContext underlying but with the retry feature based on the passed in reqInfo parameter.
func DialUnix(network, address string, reqInfo *RequestInfo) (*net.UnixConn, error) {
	conn, _, err := dialerCommon(nil, nil, network, address, "", reqInfo, 0, false)
	if conn != nil {
		return conn.(*net.UnixConn), err
	} else {
		return nil, err
	}
}

func dialerCommon(dialer *net.Dialer, resolver *net.Resolver, network, address, name string, reqInfo *RequestInfo, mode int, needContext bool) (interface{}, interface{}, error) {
	var (
		conn     net.Conn
		err      error
		retryCnt = 1
		ctx      context.Context
		cancel   context.CancelFunc
		names    []string
		cname    string
		ip       []net.IPAddr
		mx       []*net.MX
		ns       []*net.NS
		port     int
		addrs    []*net.SRV
	)
	if needContext {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(reqInfo.TimeoutSec)*time.Second)
		defer cancel()
	}
	if mode == 0 {
		conn, err = net.DialTimeout(network, address, time.Duration(reqInfo.TimeoutSec)*time.Second)
	} else if mode == 1 {
		conn, err = dialer.Dial(network, address)
	} else if mode == 2 {
		conn, err = dialer.DialContext(ctx, network, address)
	} else if mode == 3 {
		names, err = resolver.LookupAddr(ctx, address)
	} else if mode == 4 {
		cname, err = resolver.LookupCNAME(ctx, address)
	} else if mode == 5 {
		names, err = resolver.LookupHost(ctx, address)
	} else if mode == 6 {
		ip, err = resolver.LookupIPAddr(ctx, address)
	} else if mode == 7 {
		mx, err = resolver.LookupMX(ctx, address)
	} else if mode == 8 {
		ns, err = resolver.LookupNS(ctx, address)
	} else if mode == 9 {
		port, err = resolver.LookupPort(ctx, network, address)
	} else if mode == 10 {
		cname, addrs, err = resolver.LookupSRV(ctx, network, address, name)
	} else if mode == 11 {
		names, err = resolver.LookupTXT(ctx, address)
	}
	for err != nil && retryCnt < reqInfo.RetryTimes {
		retryCnt++
		<-time.After(time.Duration(reqInfo.WaitBeforeRetrySec) * time.Second)
		if mode == 0 {
			conn, err = net.DialTimeout(network, address, time.Duration(reqInfo.TimeoutSec)*time.Second)
		} else if mode == 1 {
			conn, err = dialer.Dial(network, address)
		} else if mode == 2 {
			conn, err = dialer.DialContext(ctx, network, address)
		} else if mode == 3 {
			names, err = resolver.LookupAddr(ctx, address)
		} else if mode == 4 {
			cname, err = resolver.LookupCNAME(ctx, address)
		} else if mode == 5 {
			names, err = resolver.LookupHost(ctx, address)
		} else if mode == 6 {
			ip, err = resolver.LookupIPAddr(ctx, address)
		} else if mode == 7 {
			mx, err = resolver.LookupMX(ctx, address)
		} else if mode == 8 {
			ns, err = resolver.LookupNS(ctx, address)
		} else if mode == 9 {
			port, err = resolver.LookupPort(ctx, network, address)
		} else if mode == 10 {
			cname, addrs, err = resolver.LookupSRV(ctx, network, address, name)
		} else if mode == 11 {
			names, err = resolver.LookupTXT(ctx, address)
		}
	}
	if err == nil {
		switch mode {
		case 0, 1, 2:
			return conn, nil, err
		case 3, 5, 11:
			return names, nil, err
		case 4:
			return cname, nil, err
		case 6:
			return ip, nil, err
		case 7:
			return mx, nil, err
		case 8:
			return ns, nil, err
		case 9:
			return port, nil, err
		case 10:
			return cname, addrs, err
		default:
			return nil, nil, err
		}
	} else {
		return nil, nil, err
	}
}

// NewSmtpClient. call Go net.Dial follow by smtp.NewClient underlying but with the retry feature based on the passed in reqInfo parameter.
func NewSmtpClient(network, address string, reqInfo *RequestInfo) (*smtp.Client, error) {
	var (
		conn   interface{}
		err    error
		client *smtp.Client
	)
	conn, _, err = dialerCommon(nil, nil, network, address, "", reqInfo, 0, false)
	if err != nil {
		return nil, err
	}
	client, err = smtp.NewClient(conn.(net.Conn), address)
	if err != nil {
		return nil, err
	}
	return client, err
}

// LookupAddr. call Go resolver.LookupAddr underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupAddr(resolver *net.Resolver, addr string, reqInfo *RequestInfo) (names []string, err error) {
	r, _, err := dialerCommon(nil, resolver, "", addr, "", reqInfo, 3, true)
	if r != nil {
		return r.([]string), err
	} else {
		return nil, err
	}
}

// LookupCNAME. call Go resolver.LookupCNAME underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupCNAME(resolver *net.Resolver, host string, reqInfo *RequestInfo) (cname string, err error) {
	r, _, err := dialerCommon(nil, resolver, "", host, "", reqInfo, 4, true)
	if r != nil {
		return r.(string), err
	} else {
		return "", err
	}
}

// LookupHost. call Go resolver.LookupHost underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupHost(resolver *net.Resolver, host string, reqInfo *RequestInfo) (addrs []string, err error) {
	r, _, err := dialerCommon(nil, resolver, "", host, "", reqInfo, 5, true)
	if r != nil {
		return r.([]string), err
	} else {
		return nil, err
	}
}

// LookupIPAddr. call Go resolver.LookupIPAddr underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupIPAddr(resolver *net.Resolver, host string, reqInfo *RequestInfo) ([]net.IPAddr, error) {
	r, _, err := dialerCommon(nil, resolver, "", host, "", reqInfo, 6, true)
	if r != nil {
		return r.([]net.IPAddr), err
	} else {
		return nil, err
	}
}

// LookupMX. call Go resolver.LookupMX underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupMX(resolver *net.Resolver, name string, reqInfo *RequestInfo) ([]*net.MX, error) {
	r, _, err := dialerCommon(nil, resolver, "", name, "", reqInfo, 7, true)
	if r != nil {
		return r.([]*net.MX), err
	} else {
		return nil, err
	}
}

// LookupNS. call Go resolver.LookupNS underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupNS(resolver *net.Resolver, name string, reqInfo *RequestInfo) ([]*net.NS, error) {
	r, _, err := dialerCommon(nil, resolver, "", name, "", reqInfo, 8, true)
	if r != nil {
		return r.([]*net.NS), err
	} else {
		return nil, err
	}

}

// LookupPort. call Go resolver.LookupPort underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupPort(resolver *net.Resolver, network, service string, reqInfo *RequestInfo) (port int, err error) {
	r, _, err := dialerCommon(nil, resolver, network, service, "", reqInfo, 9, true)
	if r != nil {
		return r.(int), err
	} else {
		return -1, err
	}
}

// LookupSRV. call Go resolver.LookupSRV underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupSRV(resolver *net.Resolver, service, proto, name string, reqInfo *RequestInfo) (cname string, addrs []*net.SRV, err error) {
	r, s, err := dialerCommon(nil, resolver, service, proto, name, reqInfo, 10, true)
	if r != nil {
		return r.(string), s.([]*net.SRV), err
	} else {
		return "", nil, err
	}
}

// LookupTXT. call Go resolver.LookupTXT underlying but with the retry feature based on the passed in reqInfo parameter.
func LookupTXT(resolver *net.Resolver, name string, reqInfo *RequestInfo) ([]string, error) {
	r, _, err := dialerCommon(nil, resolver, "", name, "", reqInfo, 11, true)
	if r != nil {
		return r.([]string), err
	} else {
		return nil, err
	}
}
