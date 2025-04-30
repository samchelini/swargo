package netlink

import (
	"golang.org/x/sys/unix"
	"log"
  "encoding/binary"
)

const (
	minBufSize = 8192
)

type Message interface {
	Bytes() []byte
  String() string
}

type MessageBuilder interface {
  Parse(buf []byte) (Message, error)
}

type Netlink struct {
	fd int
}

// send a netlink message
func (nl *Netlink) SendMessage(msg Message) error {
	return unix.Send(nl.fd, msg.Bytes(), 0)
}

// receive a message
func (nl *Netlink) ReceiveMessage() ([]byte, error) {
	// create a receive buffer
	buf := make([]byte, nl.getBufSize())

	// receive message
	n, _, err := unix.Recvfrom(nl.fd, buf, 0)
	return buf[:n], err
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
  // build netlink message
  builder := NewGenericMessageBuilder()
	msg := builder.
		AddNetlinkHeader(unix.GENL_ID_CTRL, Do).
		AddGenericHeader(unix.CTRL_CMD_GETFAMILY).
		AddAttributeFromString(unix.CTRL_ATTR_FAMILY_NAME, name).
		Build()

	// send to netlink
	log.Printf("getting %s family id...\n", name)
	log.Println("sending message...")
	err := nl.SendMessage(msg)
	if err != nil {
		return err
	}

	// receive response
	log.Println("receiving message...")
	resp, err := nl.ReceiveMessage()
  if err != nil {
    return err
  }
  msgLen := binary.LittleEndian.Uint32(resp)
  log.Printf("bytes received: %d\tmessage length: %d\n", len(resp), msgLen)
  log.Printf("response: % X\n", resp)

  // parse response
  log.Println("parsing response...")
  builder.Parse(resp)
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
