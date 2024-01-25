package msgpraw

import (
	"encoding/binary"
	"math"
)

type IMsgpWriter interface {
	WriteInt(int) error
	WriteInt8(int8) error
	WriteInt16(int16) error
	WriteInt32(int32) error
	WriteInt64(int64) error
	WriteUint(uint) error
	WriteUint8(uint8) error
	WriteUint16(uint16) error
	WriteUint32(uint32) error
	WriteUint64(uint64) error
	WriteFloat32(float32) error
	WriteFloat64(float64) error
	WriteString(string) error
	WriteBytes([]byte) error
	WriteNil() error
	WriteBool(bool) error
	WriteArray(int) error
	WriteMap(int) error
}

type MsgpWriter struct {
	Buffer []byte
	idx    int
}

func (w *MsgpWriter) WriteInt(i int) error {
	return w.WriteInt64(int64(i))
}

func (w *MsgpWriter) WriteInt8(i int8) error {
	w.Buffer = append(w.Buffer, byte(Int8))
	w.Buffer = append(w.Buffer, byte(i))
	return nil
}

func (w *MsgpWriter) WriteInt16(i int16) error {
	w.Buffer = append(w.Buffer, byte(Int16))
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(i))
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteInt32(i int32) error {
	w.Buffer = append(w.Buffer, byte(Int32))
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteInt64(i int64) error {
	w.Buffer = append(w.Buffer, byte(Int64))
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteUint(u uint) error {
	return w.WriteUint64(uint64(u))
}

func (w *MsgpWriter) WriteUint8(u uint8) error {
	w.Buffer = append(w.Buffer, byte(Uint8))
	w.Buffer = append(w.Buffer, u)
	return nil
}

func (w *MsgpWriter) WriteUint16(u uint16) error {
	w.Buffer = append(w.Buffer, byte(Uint16))
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, u)
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteUint32(u uint32) error {
	w.Buffer = append(w.Buffer, byte(Uint32))
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u)
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteUint64(u uint64) error {
	w.Buffer = append(w.Buffer, byte(Uint64))
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteFloat32(f float32) error {
	w.Buffer = append(w.Buffer, byte(Float32))
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(f))
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteFloat64(f float64) error {
	w.Buffer = append(w.Buffer, byte(Float64))
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteString(s string) error {
	w.Buffer = append(w.Buffer, byte(Str8))
	w.Buffer = append(w.Buffer, byte(len(s)))
	w.Buffer = append(w.Buffer, s...)
	return nil
}

func (w *MsgpWriter) WriteBytes(b []byte) error {
	w.Buffer = append(w.Buffer, byte(Bin8))
	w.Buffer = append(w.Buffer, byte(len(b)))
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteNil() error {
	w.Buffer = append(w.Buffer, byte(Nil))
	return nil
}

func (w *MsgpWriter) WriteBool(b bool) error {
	if b {
		w.Buffer = append(w.Buffer, byte(True))
	} else {
		w.Buffer = append(w.Buffer, byte(False))
	}
	return nil
}

func (w *MsgpWriter) WriteArray(length int) error {
	w.Buffer = append(w.Buffer, byte(Array16))
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(length))
	w.Buffer = append(w.Buffer, b...)
	return nil
}

func (w *MsgpWriter) WriteMap(length int) error {
	w.Buffer = append(w.Buffer, byte(Map16))
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(length))
	w.Buffer = append(w.Buffer, b...)
	return nil
}
