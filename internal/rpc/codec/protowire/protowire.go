// Package protowire defines the protobuf wire codec. Importing this package will
// register the codec.
package protowire

import (
	"google.golang.org/grpc/encoding"
)

// Name is the name registered for the proto compressor.
const Name = "protowire"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type codec struct{}

type Frame struct {
	Payload []byte
}

func (codec) Marshal(v interface{}) ([]byte, error) {
	out, ok := v.(*Frame)
	if !ok {
		panic("unmarshal not frame")
	}
	return out.Payload, nil

}

func (codec) Unmarshal(data []byte, v interface{}) error {
	dst, ok := v.(*Frame)
	if !ok {
		panic("unmarshal not frame")
	}
	dst.Payload = data
	return nil
}

func (codec) Name() string {
	return Name
}
