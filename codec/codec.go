// codec package for encoding/decoding messages over the wire
// this exists so that we use networking through easy-to-use interfaces,
// and that when this file changes
package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"sync"

	"github.com/g4me92bd777b8b16ed4c/common/types"
)

var LogTypeStringer func(b byte) string

// Debug ...
var Debug = os.Getenv("DEBUG") != ""

// DefaultMaxSize is what new codecs have
const DefaultMaxSize = 1024

// Codec to encode/decode network packets uniformly across clients and servers
// One codec per connection
// One codec should be writing to a single connection at a time
// One codec should be reading from a single connection at a time
// This should be easy enough to work with
type Codec struct {
	Version byte
	MaxSize int
	endian  binary.ByteOrder
	reader  io.Reader // for Decode
	writer  io.Writer // for Encode
	closer  io.Closer // for Close :)
	buf     []byte
	wbuf    []byte
	readmu  sync.Mutex
	writemu sync.Mutex
	encbuf  bytes.Buffer
}

// NewCodec starts everything..
// If server, conn can be nil, in which case we use New() for each conn
func NewCodec(endian binary.ByteOrder, conn io.ReadWriteCloser) *Codec {
	return &Codec{
		Version: 1,
		MaxSize: DefaultMaxSize,
		endian:  endian,
		reader:  conn,
		writer:  conn,
		closer:  conn,
	}
}

// New copies a codec but with a new connection. Server uses this for each new connection.
func (c *Codec) New(conn net.Conn) *Codec {
	return &Codec{
		Version: 1,
		endian:  c.endian,
		MaxSize: c.MaxSize,
		reader:  conn,
		writer:  conn,
		closer:  conn,
	}
}

// Close closes the network connection.
// If there is no conn, ErrConn will be returned.
func (c *Codec) Close() error {
	return c.closer.Close()
}

// Encodable is an easy to implement interface that can be sent over the wire
type Encodable interface {
	Encode(b []byte) (int, error) // should return an error if buffer is too small
	//Decode(b []byte) error        // should return first error encountered along the way
	Type() types.Type // returns 0-255, sent as prefix
	EncodeTo(w io.Writer) (int, error)
}

// ErrShort short or small
var ErrShort = errors.New("short error")

// ErrSizebig too big
var ErrSizeBig = errors.New("big error")

// ErrConn no connection
var ErrConn = errors.New("nil conn")

// Decode an object using gob
func Decode(buf []byte, vptr interface{}) (err error) {
	err = gob.NewDecoder(bytes.NewReader(buf)).Decode(vptr)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Encode an object using gob
func Encode(buf *bytes.Buffer, b []byte, vptr interface{}) (n int, err error) {
	buf.Reset()
	err = gob.NewEncoder(buf).Encode(vptr)
	n = copy(b, buf.Bytes())
	return
}

// Decode vptr from b
func (c *Codec) Decode(b []byte, vptr interface{}) (err error) {
	return Decode(b, vptr)
}

var ErrZeroType = errors.New("zero type")

// Read skips 2 byte and returns as typ, with request id req, and how many bytes to read (n) or an error.
func (c *Codec) Read(b []byte) (typ types.Type, reqid uint64, n int, err error) {
	c.readmu.Lock()
	defer c.readmu.Unlock()
	if c.buf == nil {
		c.buf = make([]byte, 1024)
	}
	fillBytes(c.buf)
	c.buf[0] = 0
	c.buf[1] = 0
	n, err = c.reader.Read(c.buf[:])
	if err != nil {
		return 0, 0, n, err
	}

	if n < 10 {
		return 0, 0, n, ErrShort
	}
	// skip 1 byte + skip 8 byte id
	typ = types.Type(c.endian.Uint16(c.buf[:2]))
	if typ == 0 {
		return 0, 0, n, ErrZeroType
	}

	if i := copy(b, c.buf[10:n]); i != n-10 {
		if Debug {
			log.Printf("%d != %d", i, n-10)
		}
		return typ, 0, n, ErrShort
	}

	reqid = c.DecodeUint64(c.buf[2:10])

	if n == c.MaxSize {
		//log.Println("Codec:", n, "max size")
		return typ, reqid, n, ErrSizeBig
	}
	if Debug {
		ts := typ.String()

		log.Printf("codec READ %s (%02d %02d) read %d bytes + 8 id bytes +  2 type byte (%d)(%d)", ts, c.buf[0], c.buf[1], n, c.buf[0], c.buf[1])
		log.Printf("codec READ %v %02d buf %02x", ts, c.buf[0:2], c.buf[2:n])
	}
	// return first byte
	return typ, reqid, n - 10, nil
}

type Typer interface {
	Type() types.Type
}

func fillBytes(b []byte) {
	for i := range b {
		b[i] = 0x00
	}
}
func (c *Codec) Write(v Typer) (n int, err error) {
	if v.Type() == 0 {
		panic("Cant encode type Zero")
	}
	c.writemu.Lock()
	defer c.writemu.Unlock()
	if c.wbuf == nil {
		c.wbuf = make([]byte, c.MaxSize)
	}
	c.endian.PutUint16(c.wbuf[0:2], v.Type().Uint16())
	copy(c.wbuf[2:], c.EncodeUint64(rand.Uint64()))

	if encoder, ok := v.(Encodable); ok {
		n, err = encoder.Encode(c.wbuf[10:])
	} else {
		n, err = Encode(&c.encbuf, c.wbuf[10:], v)
	}
	if err != nil {
		return n, err
	}
	if n+10 > c.MaxSize {
		log.Printf("Codec encode: maxsize=%d, size=%d (too big)", c.MaxSize, n)
		return n, ErrSizeBig
	}
	if Debug {
		log.Printf("codec WRITE %s (%02d %02d) encoded %d bytes + 8 id bytes +  2 type byte (%d)(%d)", v.Type().String(), c.wbuf[0], c.wbuf[1], n, c.wbuf[0], c.wbuf[1])
		log.Printf("codec WRITE %v %02d buf %02x", v.Type().String(), c.wbuf[0:2], c.wbuf[2:n])
	}

	// write type
	n, err = c.writer.Write(c.wbuf[:n+10])
	if err != nil {
		return n, err
	}
	return n, err
}

// encode/decode numbers

func (c *Codec) EncodeUint64(i uint64) []byte {
	b := [8]byte{}
	c.endian.PutUint64(b[:], i)
	return b[:]
}
func (c *Codec) DecodeUint64(b []byte) uint64 {
	return c.endian.Uint64(b)
}

func (c *Codec) EncodeFloat(i float64) []byte {
	return c.EncodeUint64(math.Float64bits(i))
}
func (c *Codec) DecodeFloat(b []byte) float64 {
	return math.Float64frombits(c.DecodeUint64(b))
}
