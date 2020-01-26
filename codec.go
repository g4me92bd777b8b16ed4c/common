package game

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/g4me92bd777b8b16ed4c/common/chatenc"
)

// https://golang.org/ref/spec#Size_and_alignment_guarantees
//
// type                                 size in bytes
//
// byte, uint8, int8                     1
// uint16, int16                         2
// uint32, int32, float32                4
// uint64, int64, float64, complex64     8
// complex128                           16

const MaxPacketSize = 1024

func Endian() binary.ByteOrder {
	return binary.LittleEndian
}

func GenUUID() uint64 {
	// fast not secure
	return rand.Uint64()
}

type Login struct {
	ID uint64
}

type Ping struct {
	ID   uint64
	Time time.Time
}

var ErrLength = errors.New("type length is wrong")
var ErrZero = errors.New("id is zero")

func (l *Login) Encode(b []byte) (n int, err error) {
	//entrybuf := make([]byte, game.MsgSize(game.MsgPrefixLogin)) // 1 + 8 + enc
	if l.ID == 0 {
		panic("ouch")
		return 0, ErrZero
	}
	n += 8
	endian.PutUint64(b[0:], l.ID)
	n += copy(b[8:], chatenc.Load("g4me92bd777b8b16ed4c.1").Encrypt([]byte("login please")))
	return n, nil

}
func (l *Login) Decode(b []byte) (err error) {
	// if len(b) != 8 {
	// 	return ErrLength
	// }
	l.ID = Endian().Uint64(b[0:8])

	if "login please" != string(chatenc.Load("g4me92bd777b8b16ed4c.1").Decrypt(b[8:])) {
		log.Printf("Bad login decod from %d: %02x", l.ID, b)
		return errors.New("could not decrypt login")
	}

	return nil
}

// func (p *Ping) String() string {
// 	return fmt.Sprintf("[ping] frm=%s at=%s", p.ID, p.Time)
// }

func (m MessageType) Byte() byte {
	return byte(m)
}

func (l *Login) Type() byte { return byte(MsgPrefixLogin) }
func (l *Ping) Type() byte  { return byte(MsgPrefixPing) }
func (l *Ping) Encode(b []byte) (n int, err error) {
	//entrybuf := make([]byte, game.MsgSize(game.MsgPrefixLogin)) // 1 + 8 + enc

	if l.ID == 0 {
		l.ID = math.MaxUint64
	}
	// endian.PutUint64(b[0:], l.ID)
	// n += 8

	buf := new(bytes.Buffer)
	err = gob.NewEncoder(buf).Encode(l)
	if err != nil {
		log.Println("Encoding Ping:", err)
		return 0, err
	}
	n += copy(b[:], buf.Bytes())
	return n, nil

}

func (p Ping) String() string {
	return fmt.Sprintf("ping from %d at %s", p.ID, p.Time)
}
func (l *Ping) Decode(b []byte) (err error) {
	if len(b) < 10 {
		return ErrLength
	}
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(l)
	if err != nil && err != io.EOF {
		log.Println("Decoding ping:", err)
		return err
	}
	log.Println("PING FROM", l.ID, "delay:", time.Since(l.Time))

	// sl := b[8:]
	// // if "login please" != string(chatenc.Load("g4me92bd777b8b16ed4c.1").Decrypt(sl[:len(sl)-2])) {
	// // 	return fmt.Errorf("could not decrypt ping:  %x , %q", sl, string(chatenc.Load("g4me92bd777b8b16ed4c.1").Decrypt(sl)))
	// // }

	return nil
}

type Pong = Ping

// codec.Encodable:
//
// type Encodable interface {
// 	Encode(b []byte) (int, error) // should return an error if buffer is too small
// 	Decode(b []byte) error        // should return first error encountered along the way
// }

func TYPE(t byte) string {
	return fmt.Sprintf("type[%d]", t)
}

func (m *Message) Encode(b []byte) (n int, err error) {

	buf := new(bytes.Buffer)
	err = gob.NewEncoder(buf).Encode(m)
	if err != nil {
		log.Println("Encoding Ping:", err)
		return 0, err
	}
	n += copy(b[:], buf.Bytes())
	return n, nil

}

func (p Message) Type() byte {
	return byte(MsgPrefixDPad)
}

func (p Message) String() string {
	return fmt.Sprintf("message={%d %d %d}", p.Dpad, p.Key, p.Keymod)
}
func (l *Message) Decode(b []byte) (err error) {
	if len(b) < 10 {
		return ErrLength
	}
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(l)
	if err != nil && err != io.EOF {
		log.Println("Decoding:", err)
		return err
	}
	return nil
}
