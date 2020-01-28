package common

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/g4me92bd777b8b16ed4c/common/chatenc"
	"github.com/g4me92bd777b8b16ed4c/common/types"
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

// MaxPacketSize anything bigger gets booted,
// or if server sends over 1024 in one packet there will be a read error on the clients
// For multipart binary transfers we use http to make things easier
const MaxPacketSize = 1024

var Debug = false

func Endian() binary.ByteOrder {
	return binary.LittleEndian
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
func GenUUID() uint64 {
	// fast not secure
	return rand.Uint64()
}

type Login struct {
	ID       uint64
	Password [32]byte `gob:",omitempty"`
}

type Ping struct {
	ID   uint64
	Time time.Time
}

var ErrLength = errors.New("type length is wrong")
var ErrZero = errors.New("id is zero")

func (l Login) Type() types.Type { return types.Login }

func (l Login) EncodeTo(w io.Writer) (n int, err error) {
	//entrybuf := make([]byte, game.MsgSize(game.MsgPrefixLogin)) // 1 + 8 + enc
	if l.ID == 0 {
		panic("ouch")
		// return 0, ErrZero
	}
	b := make([]byte, 0, 1024)
	n, err = l.Encode(b)
	if err != nil {
		return n, err
	}
	return w.Write(b[:n])

}
func (l Login) Encode(b []byte) (n int, err error) {
	//entrybuf := make([]byte, game.MsgSize(game.MsgPrefixLogin)) // 1 + 8 + enc
	if l.ID == 0 {
		panic("ouch")
		// return 0, ErrZero
	}
	endian.PutUint64(b[0:], l.ID)
	n += 8 // above uint64
	n += copy(b[n:], l.Password[:])
	n += copy(b[n:], chatenc.Load("g4me92bd777b8b16ed4c.1").Encrypt([]byte("login please")))
	return n, nil

}
func (l *Login) Decode(b []byte) (err error) {
	// if len(b) != 8 {
	// 	return ErrLength
	// }
	l.ID = Endian().Uint64(b[0:8])
	copy(l.Password[:], b[8:8+32])
	if "login please" != string(chatenc.Load("g4me92bd777b8b16ed4c.1").Decrypt(b[8+32:])) {
		log.Printf("Bad login decod from %d: %02x", l.ID, b)
		return errors.New("could not decrypt login")
	}

	return nil
}

func (l Ping) Type() types.Type { return types.Ping }
func (l Ping) Encode(b []byte) (n int, err error) {
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

// func (l *Ping) Decode(b []byte) (err error) {
// 	if len(b) < 10 {
// 		return ErrLength
// 	}
// 	err = gob.NewDecoder(bytes.NewReader(b)).Decode(l)
// 	if err != nil && err != io.EOF {
// 		if Debug {
// 			log.Println("Decoding ping:", err)
// 		}
// 		return err
// 	}
// 	log.Println("PING FROM", l.ID, "delay:", time.Since(l.Time))

// 	// sl := b[8:]
// 	// // if "login please" != string(chatenc.Load("g4me92bd777b8b16ed4c.1").Decrypt(sl[:len(sl)-2])) {
// 	// // 	return fmt.Errorf("could not decrypt ping:  %x , %q", sl, string(chatenc.Load("g4me92bd777b8b16ed4c.1").Decrypt(sl)))
// 	// // }

// 	return nil
// }

type Pong Ping

func (p Pong) Type() types.Type {
	return types.Pong
}

// codec.Encodable:
//
// type Encodable interface {
// 	Encode(b []byte) (int, error) // should return an error if buffer is too small
// 	Decode(b []byte) error        // should return first error encountered along the way
// }

func TYPE(t types.Type) string {
	return fmt.Sprintf("type[%d]", t)
}

// func (m Message) Encode(b []byte) (n int, err error) {

// 	buf := new(bytes.Buffer)
// 	err = gob.NewEncoder(buf).Encode(m)
// 	if err != nil {
// 		log.Println("Encoding Ping:", err)
// 		return 0, err
// 	}
// 	n += copy(b[:], buf.Bytes())
// 	return n, nil

// }

func (p Message) Type() types.Type {
	return types.DPad
}

func (p Message) String() string {
	s := ""
	if p.Count > 1 {
		s = fmt.Sprintf(" x%d", p.Count)
	}
	return fmt.Sprintf("message={%d %d %d}%s", p.Dpad, p.Key, p.Keymod, s)
}

// func (l *Message) Decode(b []byte) (err error) {
// 	if len(b) < 10 {
// 		return ErrLength
// 	}
// 	err = gob.NewDecoder(bytes.NewReader(b)).Decode(l)
// 	if err != nil && err != io.EOF {
// 		log.Println("Decoding:", err)
// 		return err
// 	}
// 	return nil
// }

type PlayerLogoff struct {
	UID uint64
}

// func (p PlayerLogoff) Encode(b []byte) (n int, err error) {
// 	buf := new(bytes.Buffer)
// 	err = gob.NewEncoder(buf).Encode(p)
// 	if err != nil {
// 		log.Println("Encoding Ping:", err)
// 		return 0, err
// 	}
// 	n += copy(b[:], buf.Bytes())
// 	return n, nil
// }
// func (p *PlayerLogoff) Decode(b []byte) (err error) {
// 	if len(b) < 10 {
// 		return ErrLength
// 	}
// 	err = gob.NewDecoder(bytes.NewReader(b)).Decode(p)
// 	if err != nil && err != io.EOF {
// 		log.Println("Decoding:", err)
// 		return err
// 	}
// 	return nil
// }
func (p PlayerLogoff) Type() types.Type {
	return types.PlayerLogoff
}

type PlayerMessage struct {
	From    uint64
	To      uint64
	Message string
}

// func (p PlayerMessage) Encode(b []byte) (n int, err error) {
// 	buf := new(bytes.Buffer)
// 	err = gob.NewEncoder(buf).Encode(p)
// 	if err != nil {
// 		log.Println("Encoding PlayerMessage:", err)
// 		return 0, err
// 	}
// 	n += copy(b[:], buf.Bytes())
// 	return n, nil
// }

// func (p *PlayerMessage) Decode(b []byte) (err error) {
// 	if len(b) < 10 {
// 		return ErrLength
// 	}
// 	err = gob.NewDecoder(bytes.NewReader(b)).Decode(p)
// 	if err != nil && err != io.EOF {
// 		log.Println("Decoding PlayerMessage:", err)
// 		return err
// 	}
// 	return nil
// }

func (p PlayerMessage) Type() types.Type {
	return types.PlayerMessage
}

type ServerUpdate struct {
	T          string
	UpdateTick time.Duration
	Restart    bool
}
type ClientUpdate struct {
	T          string
	UpdateTick time.Duration
}

// type C struct{}

// func (c *C) Decode(b []byte) error {
// 	return gob.NewDecoder(bytes.NewReader(b)).Decode(c)
// }
// func (c *C) Encode(b []byte) (int, error) {
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(c); err != nil {
// 		return 0, err
// 	}
// 	n := copy(b, buf.Bytes())
// 	return n, nil
// }

// func (c *ServerUpdate) Decode(b []byte) error {
// 	return gob.NewDecoder(bytes.NewReader(b)).Decode(c)
// }

// func (c ServerUpdate) Encode(b []byte) (int, error) {
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(c); err != nil {
// 		return 0, err
// 	}
// 	n := copy(b, buf.Bytes())
// 	return n, nil
// }

// func (c *ClientUpdate) Decode(b []byte) error {
// 	return gob.NewDecoder(bytes.NewReader(b)).Decode(c)
// }

// func (c ClientUpdate) Encode(b []byte) (int, error) {
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(c); err != nil {
// 		return 0, err
// 	}
// 	n := copy(b, buf.Bytes())
// 	return n, nil
// }

func (s ServerUpdate) Type() types.Type {
	return types.Update
}
func (s ClientUpdate) Type() types.Type {
	return types.Update
}

type PlayerAction struct {
	ID     uint64
	At     [2]int64
	DPad   byte
	Action uint16
	HP     float64
}

func (p PlayerAction) Pos() [2]float64 {
	return [2]float64{float64(p.At[0]), float64(p.At[1])}
}

// func (p PlayerAction) Encode(b []byte) (n int, err error) {
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(p); err != nil {
// 		return 0, err
// 	}
// 	n = copy(b, buf.Bytes())
// 	return n, nil
// }

func (p PlayerAction) Type() types.Type {
	return types.PlayerAction
}

// func (p *Player) Decode(b []byte) error {
// 	return gob.NewDecoder(bytes.NewReader(b)).Decode(p)
// }
