package net

import (
	"errors"
	"net"
)

func Serve(network, addr string) error {
	listener, err := net.Listen(network, addr)
	if err != nil {
		// 比较常见的就是端口被占用
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if er := handleConn(conn); er != nil {
				_ = conn.Close()
			}

		}()
	}
}

func handleConn(conn net.Conn) error {
	for {
		bs := make([]byte, 8)
		n, err := conn.Read(bs)
		if err != nil {
			return err
		}
		if n != 8 {
			return errors.New("micro: 没读够数据")
		}

		res := handlMsg(bs)
		n, err = conn.Write(res)
		if n != len(res) {
			return errors.New("micro: 没写完数据")
		}
		
	}
}

func handlMsg(req []byte) []byte {
	res := make([]byte, 2*len(req))
	copy(res[:len(req)], req)
	copy(res[len(req):], req)
	return res
}
