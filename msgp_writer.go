package msgpraw

import (
	"encoding/binary"
	"errors"
	"math"
)

var (
	ErrPosFixIntRange = errors.New("msgpraw: value out of PosFixInt range (0..127)")
	ErrNegFixIntRange = errors.New("msgpraw: value out of NegFixInt range (-32..-1)")
	ErrFixStrRange    = errors.New("msgpraw: string length out of FixStr range (0..31)")
	ErrFixArrayRange  = errors.New("msgpraw: array length out of FixArray range (0..15)")
	ErrFixMapRange    = errors.New("msgpraw: map length out of FixMap range (0..15)")
	ErrFixExtSize     = errors.New("msgpraw: data length does not match FixExt format size")
	ErrStr8Range      = errors.New("msgpraw: string length out of Str8 range (0..255)")
	ErrStr16Range     = errors.New("msgpraw: string length out of Str16 range (0..65535)")
	ErrStr32Range     = errors.New("msgpraw: string length out of Str32 range (0..2^32-1)")
	ErrBin8Range      = errors.New("msgpraw: data length out of Bin8 range (0..255)")
	ErrBin16Range     = errors.New("msgpraw: data length out of Bin16 range (0..65535)")
	ErrBin32Range     = errors.New("msgpraw: data length out of Bin32 range (0..2^32-1)")
	ErrArray16Range   = errors.New("msgpraw: array length out of Array16 range (0..65535)")
	ErrArray32Range   = errors.New("msgpraw: array length out of Array32 range (0..2^32-1)")
	ErrMap16Range     = errors.New("msgpraw: map length out of Map16 range (0..65535)")
	ErrMap32Range     = errors.New("msgpraw: map length out of Map32 range (0..2^32-1)")
	ErrExt8Range      = errors.New("msgpraw: data length out of Ext8 range (0..255)")
	ErrExt16Range     = errors.New("msgpraw: data length out of Ext16 range (0..65535)")
	ErrExt32Range     = errors.New("msgpraw: data length out of Ext32 range (0..2^32-1)")
	ErrStringTooLong  = errors.New("msgpraw: string too long for any Str format")
	ErrBytesTooLong   = errors.New("msgpraw: byte slice too long for any Bin format")
	ErrArrayTooLong   = errors.New("msgpraw: array too long for any Array format")
	ErrMapTooLong     = errors.New("msgpraw: map too long for any Map format")
	ErrExtTooLong     = errors.New("msgpraw: ext data too long for any Ext format")
)

const (
	maxFixStr   = 31
	maxFixArray = 15
	maxFixMap   = 15
	maxUint8    = 0xff
	maxUint16   = 0xffff
	maxUint32   = 0xffffffff
)

type IMsgpWriter interface {
	// Scalars
	WritePosFixInt(uint8) error
	WriteNegFixInt(int8) error
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
	WriteNil() error
	WriteBool(bool) error

	// Auto-sized variable-length writers (pick the smallest format that fits).
	WriteString(string) error
	WriteBytes([]byte) error
	WriteArray(int) error
	WriteMap(int) error
	WriteExt(extType int8, data []byte) error

	// Explicit fixed-format variants. Each errors if the input does not fit
	// the named format.
	WriteFixStr(string) error
	WriteStr8(string) error
	WriteStr16(string) error
	WriteStr32(string) error

	WriteBin8([]byte) error
	WriteBin16([]byte) error
	WriteBin32([]byte) error

	WriteFixArray(int) error
	WriteArray16(int) error
	WriteArray32(int) error

	WriteFixMap(int) error
	WriteMap16(int) error
	WriteMap32(int) error

	WriteFixExt1(extType int8, data []byte) error
	WriteFixExt2(extType int8, data []byte) error
	WriteFixExt4(extType int8, data []byte) error
	WriteFixExt8(extType int8, data []byte) error
	WriteFixExt16(extType int8, data []byte) error
	WriteExt8(extType int8, data []byte) error
	WriteExt16(extType int8, data []byte) error
	WriteExt32(extType int8, data []byte) error
}

// MsgpWriter appends msgpack-encoded values to Buff. Buff is exposed so
// callers can pre-size it: w := &MsgpWriter{Buff: make([]byte, 0, n)}.
type MsgpWriter struct {
	Buff []byte
}

// --- scalars ----------------------------------------------------------------

func (w *MsgpWriter) WritePosFixInt(i uint8) error {
	if i > 0x7f {
		return ErrPosFixIntRange
	}
	w.Buff = append(w.Buff, i)
	return nil
}

func (w *MsgpWriter) WriteNegFixInt(i int8) error {
	if i < -32 || i >= 0 {
		return ErrNegFixIntRange
	}
	w.Buff = append(w.Buff, byte(i))
	return nil
}

func (w *MsgpWriter) WriteInt(i int) error    { return w.WriteInt64(int64(i)) }
func (w *MsgpWriter) WriteUint(u uint) error  { return w.WriteUint64(uint64(u)) }

func (w *MsgpWriter) WriteInt8(i int8) error {
	w.Buff = append(w.Buff, byte(Int8), byte(i))
	return nil
}

func (w *MsgpWriter) WriteInt16(i int16) error {
	w.Buff = append(w.Buff, byte(Int16))
	w.Buff = binary.BigEndian.AppendUint16(w.Buff, uint16(i))
	return nil
}

func (w *MsgpWriter) WriteInt32(i int32) error {
	w.Buff = append(w.Buff, byte(Int32))
	w.Buff = binary.BigEndian.AppendUint32(w.Buff, uint32(i))
	return nil
}

func (w *MsgpWriter) WriteInt64(i int64) error {
	w.Buff = append(w.Buff, byte(Int64))
	w.Buff = binary.BigEndian.AppendUint64(w.Buff, uint64(i))
	return nil
}

func (w *MsgpWriter) WriteUint8(u uint8) error {
	w.Buff = append(w.Buff, byte(Uint8), u)
	return nil
}

func (w *MsgpWriter) WriteUint16(u uint16) error {
	w.Buff = append(w.Buff, byte(Uint16))
	w.Buff = binary.BigEndian.AppendUint16(w.Buff, u)
	return nil
}

func (w *MsgpWriter) WriteUint32(u uint32) error {
	w.Buff = append(w.Buff, byte(Uint32))
	w.Buff = binary.BigEndian.AppendUint32(w.Buff, u)
	return nil
}

func (w *MsgpWriter) WriteUint64(u uint64) error {
	w.Buff = append(w.Buff, byte(Uint64))
	w.Buff = binary.BigEndian.AppendUint64(w.Buff, u)
	return nil
}

func (w *MsgpWriter) WriteFloat32(f float32) error {
	w.Buff = append(w.Buff, byte(Float32))
	w.Buff = binary.BigEndian.AppendUint32(w.Buff, math.Float32bits(f))
	return nil
}

func (w *MsgpWriter) WriteFloat64(f float64) error {
	w.Buff = append(w.Buff, byte(Float64))
	w.Buff = binary.BigEndian.AppendUint64(w.Buff, math.Float64bits(f))
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

// --- strings (auto + explicit) ---------------------------------------------

func (w *MsgpWriter) WriteString(s string) error {
	n := len(s)
	switch {
	case n <= maxFixStr:
		return w.WriteFixStr(s)
	case n <= maxUint8:
		return w.WriteStr8(s)
	case n <= maxUint16:
		return w.WriteStr16(s)
	case uint64(n) <= maxUint32:
		return w.WriteStr32(s)
	default:
		return ErrStringTooLong
	}
}

func (w *MsgpWriter) WriteFixStr(s string) error {
	if len(s) > maxFixStr {
		return ErrFixStrRange
	}
	w.Buff = append(w.Buff, byte(FixStr)|byte(len(s)))
	w.Buff = append(w.Buff, s...)
	return nil
}

func (w *MsgpWriter) WriteStr8(s string) error {
	if len(s) > maxUint8 {
		return ErrStr8Range
	}
	w.Buff = append(w.Buff, byte(Str8), byte(len(s)))
	w.Buff = append(w.Buff, s...)
	return nil
}

func (w *MsgpWriter) WriteStr16(s string) error {
	if len(s) > maxUint16 {
		return ErrStr16Range
	}
	w.Buff = append(w.Buff, byte(Str16))
	w.Buff = binary.BigEndian.AppendUint16(w.Buff, uint16(len(s)))
	w.Buff = append(w.Buff, s...)
	return nil
}

func (w *MsgpWriter) WriteStr32(s string) error {
	if uint64(len(s)) > maxUint32 {
		return ErrStr32Range
	}
	w.Buff = append(w.Buff, byte(Str32))
	w.Buff = binary.BigEndian.AppendUint32(w.Buff, uint32(len(s)))
	w.Buff = append(w.Buff, s...)
	return nil
}

// --- bytes (auto + explicit) ------------------------------------------------

func (w *MsgpWriter) WriteBytes(b []byte) error {
	n := len(b)
	switch {
	case n <= maxUint8:
		return w.WriteBin8(b)
	case n <= maxUint16:
		return w.WriteBin16(b)
	case uint64(n) <= maxUint32:
		return w.WriteBin32(b)
	default:
		return ErrBytesTooLong
	}
}

func (w *MsgpWriter) WriteBin8(b []byte) error {
	if len(b) > maxUint8 {
		return ErrBin8Range
	}
	w.Buff = append(w.Buff, byte(Bin8), byte(len(b)))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteBin16(b []byte) error {
	if len(b) > maxUint16 {
		return ErrBin16Range
	}
	w.Buff = append(w.Buff, byte(Bin16))
	w.Buff = binary.BigEndian.AppendUint16(w.Buff, uint16(len(b)))
	w.Buff = append(w.Buff, b...)
	return nil
}

func (w *MsgpWriter) WriteBin32(b []byte) error {
	if uint64(len(b)) > maxUint32 {
		return ErrBin32Range
	}
	w.Buff = append(w.Buff, byte(Bin32))
	w.Buff = binary.BigEndian.AppendUint32(w.Buff, uint32(len(b)))
	w.Buff = append(w.Buff, b...)
	return nil
}

// --- arrays (auto + explicit) -----------------------------------------------

func (w *MsgpWriter) WriteArray(n int) error {
	switch {
	case n < 0:
		return ErrArray32Range
	case n <= maxFixArray:
		return w.WriteFixArray(n)
	case n <= maxUint16:
		return w.WriteArray16(n)
	case uint64(n) <= maxUint32:
		return w.WriteArray32(n)
	default:
		return ErrArrayTooLong
	}
}

func (w *MsgpWriter) WriteFixArray(n int) error {
	if n < 0 || n > maxFixArray {
		return ErrFixArrayRange
	}
	w.Buff = append(w.Buff, byte(FixArray)|byte(n))
	return nil
}

func (w *MsgpWriter) WriteArray16(n int) error {
	if n < 0 || n > maxUint16 {
		return ErrArray16Range
	}
	w.Buff = append(w.Buff, byte(Array16))
	w.Buff = binary.BigEndian.AppendUint16(w.Buff, uint16(n))
	return nil
}

func (w *MsgpWriter) WriteArray32(n int) error {
	if n < 0 || uint64(n) > maxUint32 {
		return ErrArray32Range
	}
	w.Buff = append(w.Buff, byte(Array32))
	w.Buff = binary.BigEndian.AppendUint32(w.Buff, uint32(n))
	return nil
}

// --- maps (auto + explicit) -------------------------------------------------

func (w *MsgpWriter) WriteMap(n int) error {
	switch {
	case n < 0:
		return ErrMap32Range
	case n <= maxFixMap:
		return w.WriteFixMap(n)
	case n <= maxUint16:
		return w.WriteMap16(n)
	case uint64(n) <= maxUint32:
		return w.WriteMap32(n)
	default:
		return ErrMapTooLong
	}
}

func (w *MsgpWriter) WriteFixMap(n int) error {
	if n < 0 || n > maxFixMap {
		return ErrFixMapRange
	}
	w.Buff = append(w.Buff, byte(FixMap)|byte(n))
	return nil
}

func (w *MsgpWriter) WriteMap16(n int) error {
	if n < 0 || n > maxUint16 {
		return ErrMap16Range
	}
	w.Buff = append(w.Buff, byte(Map16))
	w.Buff = binary.BigEndian.AppendUint16(w.Buff, uint16(n))
	return nil
}

func (w *MsgpWriter) WriteMap32(n int) error {
	if n < 0 || uint64(n) > maxUint32 {
		return ErrMap32Range
	}
	w.Buff = append(w.Buff, byte(Map32))
	w.Buff = binary.BigEndian.AppendUint32(w.Buff, uint32(n))
	return nil
}

// --- ext (auto + explicit) --------------------------------------------------

func (w *MsgpWriter) WriteExt(extType int8, data []byte) error {
	switch len(data) {
	case 1:
		return w.WriteFixExt1(extType, data)
	case 2:
		return w.WriteFixExt2(extType, data)
	case 4:
		return w.WriteFixExt4(extType, data)
	case 8:
		return w.WriteFixExt8(extType, data)
	case 16:
		return w.WriteFixExt16(extType, data)
	}
	n := len(data)
	switch {
	case n <= maxUint8:
		return w.WriteExt8(extType, data)
	case n <= maxUint16:
		return w.WriteExt16(extType, data)
	case uint64(n) <= maxUint32:
		return w.WriteExt32(extType, data)
	default:
		return ErrExtTooLong
	}
}

func (w *MsgpWriter) WriteFixExt1(extType int8, data []byte) error {
	if len(data) != 1 {
		return ErrFixExtSize
	}
	w.Buff = append(w.Buff, byte(FixExt1), byte(extType))
	w.Buff = append(w.Buff, data...)
	return nil
}

func (w *MsgpWriter) WriteFixExt2(extType int8, data []byte) error {
	if len(data) != 2 {
		return ErrFixExtSize
	}
	w.Buff = append(w.Buff, byte(FixExt2), byte(extType))
	w.Buff = append(w.Buff, data...)
	return nil
}

func (w *MsgpWriter) WriteFixExt4(extType int8, data []byte) error {
	if len(data) != 4 {
		return ErrFixExtSize
	}
	w.Buff = append(w.Buff, byte(FixExt4), byte(extType))
	w.Buff = append(w.Buff, data...)
	return nil
}

func (w *MsgpWriter) WriteFixExt8(extType int8, data []byte) error {
	if len(data) != 8 {
		return ErrFixExtSize
	}
	w.Buff = append(w.Buff, byte(FixExt8), byte(extType))
	w.Buff = append(w.Buff, data...)
	return nil
}

func (w *MsgpWriter) WriteFixExt16(extType int8, data []byte) error {
	if len(data) != 16 {
		return ErrFixExtSize
	}
	w.Buff = append(w.Buff, byte(FixExt16), byte(extType))
	w.Buff = append(w.Buff, data...)
	return nil
}

func (w *MsgpWriter) WriteExt8(extType int8, data []byte) error {
	if len(data) > maxUint8 {
		return ErrExt8Range
	}
	w.Buff = append(w.Buff, byte(Ext8), byte(len(data)), byte(extType))
	w.Buff = append(w.Buff, data...)
	return nil
}

func (w *MsgpWriter) WriteExt16(extType int8, data []byte) error {
	if len(data) > maxUint16 {
		return ErrExt16Range
	}
	w.Buff = append(w.Buff, byte(Ext16))
	w.Buff = binary.BigEndian.AppendUint16(w.Buff, uint16(len(data)))
	w.Buff = append(w.Buff, byte(extType))
	w.Buff = append(w.Buff, data...)
	return nil
}

func (w *MsgpWriter) WriteExt32(extType int8, data []byte) error {
	if uint64(len(data)) > maxUint32 {
		return ErrExt32Range
	}
	w.Buff = append(w.Buff, byte(Ext32))
	w.Buff = binary.BigEndian.AppendUint32(w.Buff, uint32(len(data)))
	w.Buff = append(w.Buff, byte(extType))
	w.Buff = append(w.Buff, data...)
	return nil
}
