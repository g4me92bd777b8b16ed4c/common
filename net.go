package common

import "net"

type ConnectOptions struct {
}

func Connect(userid uint64, server string, options *ConnectOptions) (*net.TCPConn, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	return conn.(*net.TCPConn), nil
}
