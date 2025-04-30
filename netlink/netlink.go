package netlink

import (
	"golang.org/x/sys/unix"
	"log"
)

const (
	minBufSize = 8192
)

type Message interface {
	Bytes() []byte
	Parse(msg []byte) error
}

type Netlink struct {
	fd int
}

// send a netlink message
func (nl *Netlink) SendMessage(msg Message) error {
	return unix.Send(nl.fd, msg.Bytes(), 0)
}

// receive a message
func (nl *Netlink) ReceiveMessage(msg Message) error {
	// create a receive buffer
	buf := make([]byte, nl.getBufSize())

	// receive message
	_, _, err := unix.Recvfrom(nl.fd, buf, 0)
	if err != nil {
		return err
	}

	// parse message
	return msg.Parse(buf)
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
	log.Printf("getting %s family id...\n", name)
	log.Println("sending message...")
	err := nl.SendMessage(msg)
	if err != nil {
		return err
	}

	// receive response
	log.Println("receiving message...")
	resp := new(GenericMessage)
	err = nl.ReceiveMessage(resp)
	return err
}

// get netlink file descriptor
func (nl *Netlink) GetFd() int {
	return nl.fd
}

// return a buffer size
func (nl *Netlink) getBufSize() int {
	bufSize := unix.Getpagesize()
	if bufSize < minBufSize {
		bufSize = minBufSize
	}
	return bufSize
}
