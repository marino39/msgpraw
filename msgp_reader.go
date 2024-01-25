package msgpraw

import (
	"fmt"
	"io"
)

var (
	EOF = io.EOF
)

type IMsgpReader interface {
	Read() (Type, []byte, error)
}

type MsgpReader struct {
	Buff []byte
	Idx  int
}

func (r *MsgpReader) Read() (Type, int, []byte, error) {
	if r.Idx >= len(r.Buff) {
		return Type(0), 0, nil, EOF
	}

	msgpType := Type(r.Buff[r.Idx])
	r.Idx++

	if msgpType == Nil {
		return msgpType, 0, nil, nil
	} else if msgpType == True {
		return msgpType, 0, nil, nil
	} else if msgpType == False {
		return msgpType, 0, nil, nil
	} else if msgpType <= PosFixIntMax {
		return msgpType, 0, nil, nil
	} else if msgpType >= NegFixInt {
		return msgpType, 0, nil, nil
	} else if msgpType == Int8 || msgpType == Uint8 {
		r.Idx++
		return msgpType, 0, r.Buff[r.Idx-1 : r.Idx], nil
	} else if msgpType == Int16 || msgpType == Uint16 {
		r.Idx += 2
		return msgpType, 0, r.Buff[r.Idx-2 : r.Idx], nil
	} else if msgpType == Int32 || msgpType == Uint32 || msgpType == Float32 {
		r.Idx += 4
		return msgpType, 0, r.Buff[r.Idx-4 : r.Idx], nil
	} else if msgpType == Int64 || msgpType == Uint64 || msgpType == Float64 {
		r.Idx += 8
		return msgpType, 0, r.Buff[r.Idx-8 : r.Idx], nil
	} else if msgpType == Bin8 || msgpType == Ext8 || msgpType == Str8 {
		length := int(r.Buff[r.Idx])
		r.Idx++
		r.Idx += length
		return msgpType, 0, r.Buff[r.Idx-length : r.Idx], nil
	} else if msgpType == Bin16 || msgpType == Ext16 || msgpType == Str16 {
		length := int(r.Buff[r.Idx])<<8 | int(r.Buff[r.Idx+1])
		r.Idx += 2
		r.Idx += length
		return msgpType, 0, r.Buff[r.Idx-length : r.Idx], nil
	} else if msgpType == Bin32 || msgpType == Ext32 || msgpType == Str32 {
		length := int(r.Buff[r.Idx])<<24 | int(r.Buff[r.Idx+1])<<16 | int(r.Buff[r.Idx+2])<<8 | int(r.Buff[r.Idx+3])
		r.Idx += 4
		r.Idx += length
		return msgpType, 0, r.Buff[r.Idx-length : r.Idx], nil
	} else if msgpType == FixExt1 {
		r.Idx++
		return msgpType, 0, r.Buff[r.Idx : r.Idx+1], nil
	} else if msgpType == FixExt2 {
		r.Idx += 2
		return msgpType, 0, r.Buff[r.Idx-2 : r.Idx], nil
	} else if msgpType == FixExt4 {
		r.Idx += 4
		return msgpType, 0, r.Buff[r.Idx-4 : r.Idx], nil
	} else if msgpType == FixExt8 {
		r.Idx += 8
		return msgpType, 0, r.Buff[r.Idx-8 : r.Idx], nil
	} else if msgpType == FixExt16 {
		r.Idx += 16
		return msgpType, 0, r.Buff[r.Idx-16 : r.Idx], nil
	} else if msgpType >= FixArray && msgpType <= FixArrayMax {
		length := int(msgpType) - int(FixArray)
		return msgpType, length, r.Buff[r.Idx:], nil
	} else if msgpType == Array16 {
		length := int(r.Buff[r.Idx])<<8 | int(r.Buff[r.Idx+1])
		r.Idx += 2
		return msgpType, length, r.Buff[r.Idx:], nil
	} else if msgpType == Array32 {
		length := int(r.Buff[r.Idx])<<24 | int(r.Buff[r.Idx+1])<<16 | int(r.Buff[r.Idx+2])<<8 | int(r.Buff[r.Idx+3])
		r.Idx += 4
		return msgpType, length, r.Buff[r.Idx:], nil
	} else if msgpType >= FixMap && msgpType <= FixMapMax {
		length := int(msgpType) - int(FixMap)
		return msgpType, length, r.Buff[r.Idx:], nil
	} else if msgpType == Map16 {
		length := int(r.Buff[r.Idx])<<8 | int(r.Buff[r.Idx+1])
		r.Idx += 2
		return msgpType, length, r.Buff[r.Idx:], nil
	} else if msgpType == Map32 {
		length := int(r.Buff[r.Idx])<<24 | int(r.Buff[r.Idx+1])<<16 | int(r.Buff[r.Idx+2])<<8 | int(r.Buff[r.Idx+3])
		r.Idx += 4
		return msgpType, length, r.Buff[r.Idx:], nil
	} else {
		return msgpType, 0, nil, fmt.Errorf("unknown msgp type: %x", msgpType)
	}
}
