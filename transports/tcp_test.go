package transports

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/SierraSoftworks/chat/protocol"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTCPTransport(t *testing.T) {
	Convey("TCPTransport", t, func(c C) {
		wg := sync.WaitGroup{}

		svr, err := net.Listen("tcp", "127.0.0.1:0")
		So(err, ShouldBeNil)

		var svrConn net.Conn

		wg.Add(1)
		go func() {
			conn, err := svr.Accept()
			c.So(err, ShouldBeNil)
			c.So(conn, ShouldNotBeNil)
			svrConn = conn
			wg.Done()
		}()

		clConn, err := net.Dial(svr.Addr().Network(), svr.Addr().String())
		So(err, ShouldBeNil)
		So(clConn, ShouldNotBeNil)

		wg.Wait()
		defer clConn.Close()
		defer svrConn.Close()
		defer svr.Close()

		clTr := NewTcpTransport(clConn)
		So(clTr, ShouldNotBeNil)

		svrTr := NewTcpTransport(svrConn)
		So(svrTr, ShouldNotBeNil)

		Convey("Send/Receive", func() {
			cmd := protocol.NewPing()

			So(svrTr.Send(cmd), ShouldBeNil)

			cmd2, err := clTr.Receive(100 * time.Millisecond)
			So(err, ShouldBeNil)
			So(cmd2, ShouldResemble, cmd)
		})
	})
}
