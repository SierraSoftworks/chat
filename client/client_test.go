package client

import (
	"errors"
	"testing"
	"time"

	"github.com/SierraSoftworks/chat/protocol"
	"github.com/SierraSoftworks/chat/transports"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	Convey("Client", t, func() {
		tr1, tr2 := transports.NewTestTransports()
		So(tr1, ShouldNotBeNil)
		So(tr2, ShouldNotBeNil)
		defer tr1.Disconnect()
		defer tr2.Disconnect()

		cl := NewClient("test", tr1)
		So(cl, ShouldNotBeNil)
		defer cl.Disconnect()

		So(cl.Name, ShouldEqual, "test")

		Convey("SendMessage()", func() {
			So(cl.SendMessage("This is a test"), ShouldBeNil)

			cmd, err := tr2.Receive(0)
			So(err, ShouldBeNil)
			So(cmd, ShouldNotBeNil)
			So(cmd.Operand, ShouldEqual, "msg")
			So(cmd.Arguments, ShouldResemble, []string{"test", "This is a test"})
		})

		Convey("Receiving", func() {
			Convey(".ping", func() {
				So(tr2.Send(protocol.NewPing()), ShouldBeNil)

				cmd, err := tr2.Receive(500 * time.Millisecond)
				So(err, ShouldBeNil)
				So(cmd, ShouldNotBeNil)
				So(cmd.Operand, ShouldEqual, "ping")
			})

			Convey(".msg", func() {
				So(tr2.Send(protocol.NewMessage("test", "This is a test message")), ShouldBeNil)

				select {
				case msg := <-cl.Messages():
					So(msg.Operand, ShouldEqual, "msg")
					So(msg.Sender(), ShouldEqual, "test")
					So(msg.Message(), ShouldEqual, "This is a test message")
				case <-time.After(time.Second):
					So(errors.New("timeout"), ShouldBeNil)
				}
			})
		})
	})
}
