package codec

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"time"
)

var MaxSize = 1024

// Codec to encode/decode network packets uniformly across clients and servers
type Codec struct {
	Version byte
	endian  binary.ByteOrder

	reader io.Reader // for Decode
	writer io.Writer // for Encode
	closer io.Closer // for Close :)

	buf  []byte
	wbuf []byte
}

// NewCodec starts everything.. conn can be nil
func NewCodec(endian binary.ByteOrder, conn net.Conn) *Codec {
	rand.Seed(time.Now().UTC().UnixNano())
	return &Codec{
		Version: 1,
		endian:  endian,

		// reader:  bufio.NewReader(conn),
		// writer:  bufio.NewWriter(conn),
		reader: conn,
		writer: conn,
		closer: conn,
	}
}

// New copies a codec with a new connection
func (c *Codec) New(conn net.Conn) *Codec {
	rand.Seed(time.Now().UTC().UnixNano())
	return &Codec{
		Version: 1,
		endian:  c.endian,

		reader: conn,
		writer: conn,
		closer: conn,
	}
}

func (c *Codec) Close() error {
	return c.closer.Close()
}

// Encodable is an easy to implement interface that can be sent over the wire
type Encodable interface {
	Encode(b []byte) (int, error) // should return an error if buffer is too small
	Decode(b []byte) error        // should return first error encountered along the way
	Type() byte                   // returns 0-255, sent as prefix
}

var ErrSizeMismatch = errors.New("type size bad")
var ErrShort = errors.New("short error")
var ErrSizeBig = errors.New("big error")
var Debug = false

func (c *Codec) Read(b []byte) (typ byte, reqid uint64, n int, err error) {
	if c.buf == nil {
		c.buf = make([]byte, MaxSize)
	}
	n, err = c.reader.Read(c.buf[:])
	if err != nil && n == 0 {
		log.Println("err != nil:", err)
		return 0, 0, n, err
	}

	if n == MaxSize {
		return 0, 0, n, ErrSizeBig
	}

	// skip 1 byte + skip 8 byte id
	if i := copy(b, c.buf[9:n]); i != n-9 {
		if Debug {
			log.Printf("%d != %d", i, n-9)
		}
		return c.buf[0], 0, n, ErrShort
	}

	if Debug {
		log.Printf("codec %02d read %d bytes + 8 id bytes +  1 type byte (%d)", c.buf[0], c.buf[0])
		log.Printf("codec %02d buf %02x", c.buf[0], c.buf[0:n])
	}
	// return first byte
	return c.buf[0], c.DecodeUint64(c.buf[1:9]), n - 9, nil // actually n-9, but caller needs to count the connection bytes
}
func (c *Codec) Write(v Encodable) (n int, err error) {
	if c.wbuf == nil {
		c.wbuf = make([]byte, MaxSize)
	}
	c.wbuf[0] = v.Type()
	copy(c.wbuf[1:], c.EncodeUint64(rand.Uint64()))
	n, err = v.Encode(c.wbuf[9:])
	if err != nil {
		return 0, err
	}

	if Debug {
		log.Printf("codec %02d wrote %d bytes + 8 id bytes + 1 type byte (%d)", c.wbuf[0], n, c.wbuf[0])
		log.Printf("codec %02d buf %02x", c.wbuf[0], c.wbuf[0:n])
	}

	if n > MaxSize {
		return n, ErrSizeBig
	}
	// write type
	n, err = c.writer.Write(c.wbuf[:n+9])
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
