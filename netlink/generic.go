package netlink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/unix"
)

const (
	Generic = unix.NETLINK_GENERIC
	Do      = unix.NLM_F_REQUEST | unix.NLM_F_ACK
)

type GenericMessage struct {
	hdr    *unix.NlMsghdr
	genHdr *unix.Genlmsghdr
	attrs  []*GenericAttribute
}

type GenericAttribute struct {
	unix.NlAttr
	data    []byte
	padding []byte
}

type GenericMessageBuilder struct {
	msg *GenericMessage
	err error
}

// return a generic netlink message builder
func NewGenericMessageBuilder() *GenericMessageBuilder {
	builder := new(GenericMessageBuilder)
	builder.msg = new(GenericMessage)
	builder.msg.attrs = make([]*GenericAttribute, 0)
	return builder
}

// add the netlink message header
func (b *GenericMessageBuilder) AddNetlinkHeader(family, flags uint16) *GenericMessageBuilder {
	b.msg.hdr = &unix.NlMsghdr{
		Type:  family,
		Flags: flags,
	}
	return b
}

// add the generic netlink header
func (b *GenericMessageBuilder) AddGenericHeader(cmd uint8) *GenericMessageBuilder {
	b.msg.genHdr = &unix.Genlmsghdr{Cmd: cmd, Version: 1}
	return b
}

// add an attribute from string
func (b *GenericMessageBuilder) AddAttributeFromString(nlaType uint16, str string) *GenericMessageBuilder {
	var err error
	attr := new(GenericAttribute)
	attr.Type = nlaType
	attr.data, err = unix.ByteSliceFromString(str)
	if err != nil {
		b.err = err
	}
	attr.pad()
	b.msg.attrs = append(b.msg.attrs, attr)
	return b
}

// finalize the lengths of the headers
func (b *GenericMessageBuilder) Build() *GenericMessage {
	b.msg.hdr.Len = unix.SizeofNlMsghdr + unix.SizeofNlAttr
	for i := 0; i < len(b.msg.attrs); i++ {
		b.msg.attrs[i].Len = unix.SizeofNlAttr + uint16(len(b.msg.attrs[i].data))
		b.msg.hdr.Len += uint32(b.msg.attrs[i].Len) + uint32(len(b.msg.attrs[i].padding))
	}
	return b.msg
}

// string representation of message
func (m *GenericMessage) String() string {
	var attrs []string
	for i := 0; i < len(m.attrs); i++ {
		attrs = append(attrs, fmt.Sprint(m.attrs[i]))
	}
	return fmt.Sprintf("Netlink Header: %v\nGeneric Header: %v\nAttributes: %v", *m.hdr, *m.genHdr, attrs)
}

// byte representation of message (little endian)
func (m *GenericMessage) Bytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, m.hdr)
	binary.Write(&buf, binary.LittleEndian, m.genHdr)
	for i := 0; i < len(m.attrs); i++ {
		binary.Write(&buf, binary.LittleEndian, m.attrs[i].Len)
		binary.Write(&buf, binary.LittleEndian, m.attrs[i].Type)
		binary.Write(&buf, binary.LittleEndian, m.attrs[i].data)
		binary.Write(&buf, binary.LittleEndian, m.attrs[i].padding)
	}
	return buf.Bytes()
}

// parse and return a GenericMessage
func (b *GenericMessageBuilder) Parse(msg []byte) (*GenericMessage, error) {
  b.msg = new(GenericMessage)
  b.msg.hdr = new(unix.NlMsghdr)
  b.msg.genHdr = new(unix.Genlmsghdr)
	b.msg.attrs = make([]*GenericAttribute, 0)
  reader := bytes.NewReader(msg)
  err := binary.Read(reader, binary.LittleEndian, b.msg.hdr)
  if err != nil {
    return nil, err
  }
  err = binary.Read(reader, binary.LittleEndian, b.msg.genHdr)
  if err != nil {
    return nil, err
  }

  for err == nil {
    attr := new(GenericAttribute)
    err = binary.Read(reader, binary.LittleEndian, &attr.Len)
    err = binary.Read(reader, binary.LittleEndian, &attr.Type)
    attr.data = make([]byte, attr.Len - unix.SizeofNlAttr)
    attr.padding = make([]byte, (unix.NLMSG_ALIGNTO - (len(attr.data) % unix.NLMSG_ALIGNTO)) % unix.NLMSG_ALIGNTO)
    fmt.Printf("padding len: %d\n", len(attr.padding))
    fmt.Printf("calc: %d\n\n", attr.Len % unix.NLMSG_ALIGNTO)
    err = binary.Read(reader, binary.LittleEndian, attr.data)
    if len(attr.padding) != 0 {
      err = binary.Read(reader, binary.LittleEndian, attr.padding)
    }
    b.msg.attrs = append(b.msg.attrs, attr)
  }
  
	return b.msg, err
}

// add padding to align data to 4 bytes (NLMSG_ALIGNTO)
func (attr *GenericAttribute) pad() {
	n := unix.NLMSG_ALIGNTO - (len(attr.data) % unix.NLMSG_ALIGNTO)
	if n > 0 {
		for i := 0; i < n; i++ {
			attr.padding = append(attr.padding, byte(0))
		}
	}
}
