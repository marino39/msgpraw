package msgpraw

import (
	"errors"
	"io"
)

var (
	EOF             = io.EOF
	ErrTruncated    = errors.New("msgpraw: truncated input")
	ErrUnknownType  = errors.New("msgpraw: unknown msgp type")
)

type IMsgpReader interface {
	// Read reads the next msgp value.
	//
	// Return values:
	//   Type   - the msgp tag byte
	//   int    - element count for FixArray/Array16/Array32/FixMap/Map16/Map32; otherwise 0
	//   []byte - raw payload (sub-slice of Buff, no copy):
	//              scalars         : the big-endian-encoded value bytes
	//              FixExt*/Ext*    : [extType, dataBytes...] - caller does data[0], data[1:]
	//              FixArray/Array* : remainder of buffer; caller calls Read repeatedly
	//              FixMap/Map*     : remainder of buffer; caller calls Read 2*count times
	//              Nil/True/False/PosFixInt/NegFixInt: nil
	//   error  - io.EOF when buffer exhausted, ErrTruncated when input ends mid-value,
	//            ErrUnknownType when the tag is not a valid msgp format
	Read() (Type, int, []byte, error)
	Skip() error
}

type MsgpReader struct {
	Buff []byte
	Idx  int
}

func (r *MsgpReader) need(n int) error {
	if r.Idx+n > len(r.Buff) {
		return ErrTruncated
	}
	return nil
}

// Read reads the next msgp value. See IMsgpReader for return value semantics.
func (r *MsgpReader) Read() (Type, int, []byte, error) {
	if r.Idx >= len(r.Buff) {
		return Type(0), 0, nil, EOF
	}

	msgpType := Type(r.Buff[r.Idx])
	r.Idx++

	switch {
	case msgpType == Nil, msgpType == True, msgpType == False:
		return msgpType, 0, nil, nil

	case msgpType <= PosFixIntMax:
		return msgpType, 0, nil, nil

	case msgpType >= NegFixInt:
		return msgpType, 0, nil, nil

	case msgpType == Int8, msgpType == Uint8:
		if err := r.need(1); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx++
		return msgpType, 0, r.Buff[r.Idx-1 : r.Idx], nil

	case msgpType == Int16, msgpType == Uint16:
		if err := r.need(2); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += 2
		return msgpType, 0, r.Buff[r.Idx-2 : r.Idx], nil

	case msgpType == Int32, msgpType == Uint32, msgpType == Float32:
		if err := r.need(4); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += 4
		return msgpType, 0, r.Buff[r.Idx-4 : r.Idx], nil

	case msgpType == Int64, msgpType == Uint64, msgpType == Float64:
		if err := r.need(8); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += 8
		return msgpType, 0, r.Buff[r.Idx-8 : r.Idx], nil

	case msgpType >= FixStr && msgpType <= FixStrMax:
		length := int(msgpType) - int(FixStr)
		if err := r.need(length); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += length
		return msgpType, 0, r.Buff[r.Idx-length : r.Idx], nil

	case msgpType == Bin8, msgpType == Str8:
		if err := r.need(1); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])
		r.Idx++
		if err := r.need(length); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += length
		return msgpType, 0, r.Buff[r.Idx-length : r.Idx], nil

	case msgpType == Bin16, msgpType == Str16:
		if err := r.need(2); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])<<8 | int(r.Buff[r.Idx+1])
		r.Idx += 2
		if err := r.need(length); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += length
		return msgpType, 0, r.Buff[r.Idx-length : r.Idx], nil

	case msgpType == Bin32, msgpType == Str32:
		if err := r.need(4); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])<<24 | int(r.Buff[r.Idx+1])<<16 | int(r.Buff[r.Idx+2])<<8 | int(r.Buff[r.Idx+3])
		r.Idx += 4
		if err := r.need(length); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += length
		return msgpType, 0, r.Buff[r.Idx-length : r.Idx], nil

	case msgpType == Ext8:
		if err := r.need(1); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])
		r.Idx++
		// payload = 1 type byte + length data bytes
		total := 1 + length
		if err := r.need(total); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += total
		return msgpType, 0, r.Buff[r.Idx-total : r.Idx], nil

	case msgpType == Ext16:
		if err := r.need(2); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])<<8 | int(r.Buff[r.Idx+1])
		r.Idx += 2
		total := 1 + length
		if err := r.need(total); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += total
		return msgpType, 0, r.Buff[r.Idx-total : r.Idx], nil

	case msgpType == Ext32:
		if err := r.need(4); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])<<24 | int(r.Buff[r.Idx+1])<<16 | int(r.Buff[r.Idx+2])<<8 | int(r.Buff[r.Idx+3])
		r.Idx += 4
		total := 1 + length
		if err := r.need(total); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += total
		return msgpType, 0, r.Buff[r.Idx-total : r.Idx], nil

	case msgpType == FixExt1:
		// 1 type byte + 1 data byte
		if err := r.need(2); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += 2
		return msgpType, 0, r.Buff[r.Idx-2 : r.Idx], nil

	case msgpType == FixExt2:
		// 1 type byte + 2 data bytes
		if err := r.need(3); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += 3
		return msgpType, 0, r.Buff[r.Idx-3 : r.Idx], nil

	case msgpType == FixExt4:
		if err := r.need(5); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += 5
		return msgpType, 0, r.Buff[r.Idx-5 : r.Idx], nil

	case msgpType == FixExt8:
		if err := r.need(9); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += 9
		return msgpType, 0, r.Buff[r.Idx-9 : r.Idx], nil

	case msgpType == FixExt16:
		if err := r.need(17); err != nil {
			return msgpType, 0, nil, err
		}
		r.Idx += 17
		return msgpType, 0, r.Buff[r.Idx-17 : r.Idx], nil

	case msgpType >= FixArray && msgpType <= FixArrayMax:
		length := int(msgpType) - int(FixArray)
		return msgpType, length, r.Buff[r.Idx:], nil

	case msgpType == Array16:
		if err := r.need(2); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])<<8 | int(r.Buff[r.Idx+1])
		r.Idx += 2
		return msgpType, length, r.Buff[r.Idx:], nil

	case msgpType == Array32:
		if err := r.need(4); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])<<24 | int(r.Buff[r.Idx+1])<<16 | int(r.Buff[r.Idx+2])<<8 | int(r.Buff[r.Idx+3])
		r.Idx += 4
		return msgpType, length, r.Buff[r.Idx:], nil

	case msgpType >= FixMap && msgpType <= FixMapMax:
		length := int(msgpType) - int(FixMap)
		return msgpType, length, r.Buff[r.Idx:], nil

	case msgpType == Map16:
		if err := r.need(2); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])<<8 | int(r.Buff[r.Idx+1])
		r.Idx += 2
		return msgpType, length, r.Buff[r.Idx:], nil

	case msgpType == Map32:
		if err := r.need(4); err != nil {
			return msgpType, 0, nil, err
		}
		length := int(r.Buff[r.Idx])<<24 | int(r.Buff[r.Idx+1])<<16 | int(r.Buff[r.Idx+2])<<8 | int(r.Buff[r.Idx+3])
		r.Idx += 4
		return msgpType, length, r.Buff[r.Idx:], nil

	default:
		return msgpType, 0, nil, ErrUnknownType
	}
}

// Skip skips the next msgp value.
func (r *MsgpReader) Skip() error {
	_, _, _, err := r.Read()
	return err
}
