package msgpraw

type Type byte

const (
	PosFixInt    Type = 0x00
	PosFixIntMax Type = 0x7f
	FixMap       Type = 0x80
	FixMapMax    Type = 0x8f
	FixArray     Type = 0x90
	FixArrayMax  Type = 0x9f
	FixStr       Type = 0xa0
	FixStrMax    Type = 0xbf
	Nil          Type = 0xc0
	False        Type = 0xc2
	True         Type = 0xc3
	Bin8         Type = 0xc4
	Bin16        Type = 0xc5
	Bin32        Type = 0xc6
	Ext8         Type = 0xc7
	Ext16        Type = 0xc8
	Ext32        Type = 0xc9
	Float32      Type = 0xca
	Float64      Type = 0xcb
	Uint8        Type = 0xcc
	Uint16       Type = 0xcd
	Uint32       Type = 0xce
	Uint64       Type = 0xcf
	Int8         Type = 0xd0
	Int16        Type = 0xd1
	Int32        Type = 0xd2
	Int64        Type = 0xd3
	FixExt1      Type = 0xd4
	FixExt2      Type = 0xd5
	FixExt4      Type = 0xd6
	FixExt8      Type = 0xd7
	FixExt16     Type = 0xd8
	Str8         Type = 0xd9
	Str16        Type = 0xda
	Str32        Type = 0xdb
	Array16      Type = 0xdc
	Array32      Type = 0xdd
	Map16        Type = 0xde
	Map32        Type = 0xdf
	NegFixInt    Type = 0xe0
	NegFixIntMax Type = 0xff
)
