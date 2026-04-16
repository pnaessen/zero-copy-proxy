package proxy

import (
	"net"
)

func HandleConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

}
