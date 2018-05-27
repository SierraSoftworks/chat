package transports

import (
	"testing"
	"time"

	"github.com/SierraSoftworks/chat/protocol"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChannelTransport(t *testing.T) {
	Convey("ChannelTransport", t, func() {
		tr := NewLoopbackTransport()
		So(tr, ShouldNotBeNil)

		tri, ok := tr.(*channelTransport)
		So(ok, ShouldBeTrue)
		So(tri, ShouldNotBeNil)

		Convey("Send()", func() {
			cmd := protocol.NewPing()

			So(tr.Send(cmd), ShouldBeNil)

			cmd2 := <-tri.ch
			So(cmd2, ShouldEqual, cmd)
		})

		Convey("Receive()", func() {
			cmd := protocol.NewPing()
			tri.ch <- cmd

			cmd2, err := tr.Receive(50 * time.Millisecond)
			So(err, ShouldBeNil)
			So(cmd2, ShouldNotBeNil)
			So(cmd2, ShouldEqual, cmd)
		})
	})
}
