package server

import (
	"testing"
	"time"

	"github.com/SierraSoftworks/chat/protocol"
	"github.com/SierraSoftworks/chat/transports"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	Convey("Client", t, func() {
		svr, err := newTcpServer("127.0.0.1:0")
		So(err, ShouldBeNil)
		defer svr.Shutdown()

		tr1, tr2 := transports.NewTestTransports()
		So(tr1, ShouldNotBeNil)
		So(tr2, ShouldNotBeNil)
		defer tr1.Disconnect()
		defer tr2.Disconnect()

		cl := newClient(svr, tr1)
		So(cl, ShouldNotBeNil)
		defer cl.Disconnect()

		Convey("Send()", func() {
			So(cl.Send(protocol.NewMessage("test", "this is a test")), ShouldBeNil)

			cmd, err := tr2.Receive(0)
			So(err, ShouldBeNil)
			So(cmd, ShouldNotBeNil)
			So(cmd.Operand, ShouldEqual, "msg")
			So(cmd.Arguments, ShouldResemble, []string{"test", "this is a test"})
		})

		Convey("Receiving", func() {
			Convey(".msg", func() {
				So(tr1.Send(protocol.NewMessage("test", "This is a test message")), ShouldBeNil)

				cmd, err := tr2.Receive(500 * time.Millisecond)
				So(err, ShouldBeNil)
				So(cmd, ShouldNotBeNil)
				So(cmd.Operand, ShouldEqual, "msg")
				So(cmd.Arguments, ShouldResemble, []string{"test", "This is a test message"})
			})
		})
	})
}
