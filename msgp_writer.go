package msgpraw

import (
	"encoding/binary"
	"fmt"
	"math"
)

type IMsgpWriter interface {
	WritePosFixInt(uint8) error
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
	Buff []byte
	idx  int
}

func (w *MsgpWriter) WriteInt(i int) error {
	return w.WriteInt64(int64(i))
}

func (w *MsgpWriter) WritePosFixInt(i uint8) error {
	if i > 0x7f {
		return fmt.Errorf("value is not a positive fixint")
	}
	w.Buff = append(w.Buff, i)
	return nil
}

func (w *MsgpWriter) WriteInt8(i int8) error {
	w.Buff = append(w.Buff, byte(Int8))
	w.Buff = append(w.Buff, byte(i))
	return nil
}

func (w *MsgpWriter) WriteInt16(i int16) error {
	w.Buff = append(w.Buff, byte(Int16))
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(i))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteInt32(i int32) error {
	w.Buff = append(w.Buff, byte(Int32))
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteInt64(i int64) error {
	w.Buff = append(w.Buff, byte(Int64))
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteUint(u uint) error {
	return w.WriteUint64(uint64(u))
}

func (w *MsgpWriter) WriteUint8(u uint8) error {
	w.Buff = append(w.Buff, byte(Uint8))
	w.Buff = append(w.Buff, u)
	return nil
}

func (w *MsgpWriter) WriteUint16(u uint16) error {
	w.Buff = append(w.Buff, byte(Uint16))
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, u)
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteUint32(u uint32) error {
	w.Buff = append(w.Buff, byte(Uint32))
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u)
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteUint64(u uint64) error {
	w.Buff = append(w.Buff, byte(Uint64))
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteFloat32(f float32) error {
	w.Buff = append(w.Buff, byte(Float32))
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(f))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteFloat64(f float64) error {
	w.Buff = append(w.Buff, byte(Float64))
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteString(s string) error {
	w.Buff = append(w.Buff, byte(Str8))
	w.Buff = append(w.Buff, byte(len(s)))
	w.Buff = append(w.Buff, s...)
	return nil
}

func (w *MsgpWriter) WriteBytes(b []byte) error {
	w.Buff = append(w.Buff, byte(Bin8))
	w.Buff = append(w.Buff, byte(len(b)))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteNil() error {
	w.Buff = append(w.Buff, byte(Nil))
	return nil
}

func (w *MsgpWriter) WriteBool(b bool) error {
	if b {
		w.Buff = append(w.Buff, byte(True))
	} else {
		w.Buff = append(w.Buff, byte(False))
	}
	return nil
}

func (w *MsgpWriter) WriteArray(length int) error {
	w.Buff = append(w.Buff, byte(Array16))
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(length))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteMap(length int) error {
	w.Buff = append(w.Buff, byte(Map16))
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(length))
	w.Buff = append(w.Buff, b...)
	return nil
}
