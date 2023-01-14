package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"bytes"
	"fmt"
	"io"
)

//HeaderSize is the length of the size, type and tag headers in bytes.
const HeaderSize uint64 = 4 + 1 + 2

//An encoder writes itself as 9p representation to a io.Writer
type encoder interface {
	encode(w io.Writer) (int64, error)
}

//A decoder reads a 9p representation of itself from a io.Reader
type decoder interface {
	decode(r io.Reader) (int64, error)
}

//Content is implemented by the message types containing the real information.
type Content interface {
	encoder
	decoder
}

/*Message is the general message type.

It wraps implementations of interface Content.
*/
type Message struct {
	typ     messageType
	Tag     uint16
	Content Content
}

// get the matching messageType for a Content
func contentType(c Content) (messageType, error) {
	switch c.(type) {
	case *RVersion:
		return rversion, nil
	case *RAuth:
		return rauth, nil
	case *RAttach:
		return rattach, nil
	case *RError:
		return rerror, nil
	case *RFlush:
		return rflush, nil
	case *RWalk:
		return rwalk, nil
	case *ROpen:
		return ropen, nil
	case *RCreate:
		return rcreate, nil
	case *RRead:
		return rread, nil
	case *RWrite:
		return rwrite, nil
	case *RClunk:
		return rclunk, nil
	case *RRemove:
		return rremove, nil
	case *RStat:
		return rstat, nil
	case *RWstat:
		return rwstat, nil
	case *TVersion:
		return tversion, nil
	case *TAuth:
		return tauth, nil
	case *TAttach:
		return tattach, nil
	case *TFlush:
		return tflush, nil
	case *TWalk:
		return twalk, nil
	case *TOpen:
		return topen, nil
	case *TCreate:
		return tcreate, nil
	case *TRead:
		return tread, nil
	case *TWrite:
		return twrite, nil
	case *TClunk:
		return tclunk, nil
	case *TRemove:
		return tremove, nil
	case *TStat:
		return tstat, nil
	case *TWstat:
		return twstat, nil
	default:
		return 0, fmt.Errorf("unknown content type: %T", c)
	}
}

/*New creates a new Message with a given tag and Content.*/
func New(tag uint16, content Content) (Message, error) {
	typ, err := contentType(content)
	if err != nil {
		return Message{}, err
	}
	return Message{typ: typ, Tag: tag, Content: content}, nil
}

//Response returns a new Message prepared for using as reply.
func (m Message) Response() Message {
	return Message{Tag: m.Tag}
}

/*Decode reads a message from a io.Reader.

It returns the count of read bytes and an error.*/
func (m *Message) Decode(r io.Reader) (int64, error) {
	size, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	typ, n2, err := basic.Uint8Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	tag, n3, err := basic.Uint16Decode(r)
	if err != nil {
		return n1 + n2 + n3, err
	}

	switch messageType(typ) {
	case tversion:
		m.Content = &TVersion{}
	case tauth:
		m.Content = &TAuth{}
	case tattach:
		m.Content = &TAttach{}
	case tflush:
		m.Content = &TFlush{}
	case twalk:
		m.Content = &TWalk{}
	case topen:
		m.Content = &TOpen{}
	case tcreate:
		m.Content = &TCreate{}
	case tread:
		m.Content = &TRead{}
	case twrite:
		m.Content = &TWrite{}
	case tclunk:
		m.Content = &TClunk{}
	case tremove:
		m.Content = &TRemove{}
	case tstat:
		m.Content = &TStat{}
	case twstat:
		m.Content = &TWstat{}
	default:
		return n1 + n2 + n3, fmt.Errorf("unknown message type: %o", typ)
	}

	n4, err := m.Content.decode(r)
	if err != nil {
		return n1 + n2 + n3 + n4, err
	}

	if uint32(n1+n2+n3+n4) != size {
		return n1 + n2 + n3 + n4, fmt.Errorf("short read: %v read, %v expected", n1+n2+n3+n4, size)
	}

	m.typ = messageType(typ)
	m.Tag = tag

	return n1 + n2 + n3 + n4, nil
}

/*Encode writes a Message to an io.Writer.

It returns the count of bytes written and an error.*/
func (m *Message) Encode(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	typ, err := contentType(m.Content)
	if err != nil {
		return 0, err
	}

	_, err = basic.Uint8Encode(&buf, uint8(typ))
	if err != nil {
		return 0, err
	}

	_, err = basic.Uint16Encode(&buf, m.Tag)
	if err != nil {
		return 0, err
	}

	_, err = m.Content.encode(&buf)
	if err != nil {
		return 0, err
	}

	size := uint32(4 + buf.Len())

	n1, err := basic.Uint32Encode(w, size)
	if err != nil {
		return n1, err
	}

	n2, err := buf.WriteTo(w)
	if err != nil {
		return n1 + n2, err
	}

	if uint32(n1+n2) != size {
		return n1 + n2, fmt.Errorf("short write: %v written %v expected", n1+n2, size)
	}

	return n1 + n2, nil
}
