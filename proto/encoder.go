package proto

import (
	"google.golang.org/protobuf/proto"
)

type Encoder struct {
	buf []byte
	enc *proto.MarshalOptions
}

func NewEncoder() *Encoder {
	return &Encoder{
		enc: new(proto.MarshalOptions),
	}
}

func (b *Encoder) Reset() {
	b.buf = b.buf[:0]
}

func (b *Encoder) Bytes() []byte {
	return b.buf
}

func (b *Encoder) Marshal(m proto.Message) (err error) {
	b.buf, err = b.enc.MarshalAppend(b.buf[:0], m)
	return err
}
