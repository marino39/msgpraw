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
	idx  int
}

func (r *MsgpReader) Read() (Type, int, []byte, error) {
	if r.idx >= len(r.Buff) {
		return Type(0), 0, nil, EOF
	}

	msgpType := Type(r.Buff[r.idx])
	r.idx++

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
		r.idx++
		return msgpType, 0, r.Buff[r.idx-1 : r.idx], nil
	} else if msgpType == Int16 || msgpType == Uint16 {
		r.idx += 2
		return msgpType, 0, r.Buff[r.idx-2 : r.idx], nil
	} else if msgpType == Int32 || msgpType == Uint32 || msgpType == Float32 {
		r.idx += 4
		return msgpType, 0, r.Buff[r.idx-4 : r.idx], nil
	} else if msgpType == Int64 || msgpType == Uint64 || msgpType == Float64 {
		r.idx += 8
		return msgpType, 0, r.Buff[r.idx-8 : r.idx], nil
	} else if msgpType == Bin8 || msgpType == Ext8 || msgpType == Str8 {
		length := int(r.Buff[r.idx])
		r.idx++
		return msgpType, 0, r.Buff[r.idx : r.idx+length], nil
	} else if msgpType == Bin16 || msgpType == Ext16 || msgpType == Str16 {
		length := int(r.Buff[r.idx])<<8 | int(r.Buff[r.idx+1])
		r.idx += 2
		return msgpType, 0, r.Buff[r.idx : r.idx+length], nil
	} else if msgpType == Bin32 || msgpType == Ext32 || msgpType == Str32 {
		length := int(r.Buff[r.idx])<<24 | int(r.Buff[r.idx+1])<<16 | int(r.Buff[r.idx+2])<<8 | int(r.Buff[r.idx+3])
		r.idx += 4
		return msgpType, 0, r.Buff[r.idx : r.idx+length], nil
	} else if msgpType == FixExt1 {
		r.idx++
		return msgpType, 0, r.Buff[r.idx : r.idx+1], nil
	} else if msgpType == FixExt2 {
		r.idx += 2
		return msgpType, 0, r.Buff[r.idx-2 : r.idx], nil
	} else if msgpType == FixExt4 {
		r.idx += 4
		return msgpType, 0, r.Buff[r.idx-4 : r.idx], nil
	} else if msgpType == FixExt8 {
		r.idx += 8
		return msgpType, 0, r.Buff[r.idx-8 : r.idx], nil
	} else if msgpType == FixExt16 {
		r.idx += 16
		return msgpType, 0, r.Buff[r.idx-16 : r.idx], nil
	} else if msgpType >= FixArray && msgpType <= FixArrayMax {
		length := int(msgpType) - int(FixArray)
		return msgpType, length, r.Buff[r.idx:], nil
	} else if msgpType == Array16 {
		length := int(r.Buff[r.idx])<<8 | int(r.Buff[r.idx+1])
		r.idx += 2
		return msgpType, length, r.Buff[r.idx:], nil
	} else if msgpType == Array32 {
		length := int(r.Buff[r.idx])<<24 | int(r.Buff[r.idx+1])<<16 | int(r.Buff[r.idx+2])<<8 | int(r.Buff[r.idx+3])
		r.idx += 4
		return msgpType, length, r.Buff[r.idx:], nil
	} else if msgpType >= FixMap && msgpType <= FixMapMax {
		length := int(msgpType) - int(FixMap)
		return msgpType, length, r.Buff[r.idx:], nil
	} else if msgpType == Map16 {
		length := int(r.Buff[r.idx])<<8 | int(r.Buff[r.idx+1])
		r.idx += 2
		return msgpType, length, r.Buff[r.idx:], nil
	} else if msgpType == Map32 {
		length := int(r.Buff[r.idx])<<24 | int(r.Buff[r.idx+1])<<16 | int(r.Buff[r.idx+2])<<8 | int(r.Buff[r.idx+3])
		r.idx += 4
		return msgpType, length, r.Buff[r.idx:], nil
	} else {
		return msgpType, 0, nil, fmt.Errorf("unknown msgp type: %x", msgpType)
	}
}
