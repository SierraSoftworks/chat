# Chat
**A simple chat protocol to practice programming**

![](https://user-images.githubusercontent.com/1760260/40576922-aee42886-60f5-11e8-8f53-2d3111217d71.gif)

Chat (en. "cat") is a simple chat protocol meant to be implemented as a means to learn network programming
and intended to be cross-platform compatible, easy to parse and simple to test.

## Protocol
Chat's protocol is composed of a command type, followed by zero or more string arguments. These commands are represented in text format, making it easy for a
person to visually debug issues or test using tools like `nc` or `telnet`.

Commands are stateless, which helps keep the client and server as simple as
possible while also easing debugging. There is also (currently) a single global
chat room to which all clients are connected, again, simplifying the server
implementation.

```
.msg sender "This is my message"
```

For a more comprehensive description of the protocol's grammar and available commands,
please look at the [PROTOCOL.md](PROTOCOL.md) file.

## Reference Implementation
There is a reference implementation of the client, server and protocol parser written
in Go. Most of these components are somewhat over-engineered with the objective of
making them easily testable and extensible for future expansions to the protocol and
feature-set of *Chat*. It is highly likely that these implementations could be
shortened down to a few dozen lines of code.

You can download pre-compiled versions of the *Chat* binaries on the
[releases](https://github.com/SierraSoftworks/chat/releases) page if you wish to
use them for testing your own clients or servers.