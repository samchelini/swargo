package netlink

import (
	"golang.org/x/sys/unix"
	"log"
)

type Message interface {
	Bytes() []byte
	Parse(msg []byte) error
}

type Netlink struct {
	fd int
}

func (nl *Netlink) SendMessage(m Message) error {
	return nil
}

// bind to netlink socket address
func (nl *Netlink) bind() error {
	return unix.Bind(nl.fd, &unix.SockaddrNetlink{Family: unix.AF_NETLINK})
}

// returns a netlink connection
func Dial(family int) (*Netlink, error) {
	var err error

	// create the socket
	conn := new(Netlink)
	conn.fd, err = unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, family)
	if err != nil {
		return nil, err
	}

	// connect to netlink
	err = conn.bind()
	return conn, err
}

func (nl *Netlink) GetFamilyId(name string) error {
	msg := NewGenericMessageBuilder().
    AddNetlinkHeader(unix.GENL_ID_CTRL, Do).
		AddGenericHeader(unix.CTRL_CMD_GETFAMILY).
    AddAttributeFromString(unix.CTRL_ATTR_FAMILY_NAME, name).
		Build()

	// send to network
	log.Println("sending message...")
	err := unix.Send(nl.fd, msg.Bytes(), 0)
	if err != nil {
		return err
	}

	// create receive buffer
  bufSize := unix.Getpagesize()
  if bufSize < 8192 {
    bufSize = 8192
  }
	log.Printf("buffer size: %d\n", bufSize)
	rbuf := make([]byte, bufSize)

	// receive response
	log.Println("receiving message...")
	n, _, err := unix.Recvfrom(nl.fd, rbuf, 0)
	if err != nil {
		return err
	}
	log.Printf("bytes received: %d\n", n)

	return nil
}

func (nl *Netlink) GetFd() int {
	return nl.fd
}

