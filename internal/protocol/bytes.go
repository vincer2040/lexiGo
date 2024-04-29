package protocol

type typeByte = byte

const (
	ARRAY_BYTE         typeByte = '*'
	SIMPLE_STRING_BYTE typeByte = '+'
	BULK_STRING_BYTE   typeByte = '$'
	INT_BYTE           typeByte = ':'
	DBL_BYTE           typeByte = ','
	SIMPLE_ERROR_BYTE  typeByte = '-'
	BULK_ERROR_BYTE    typeByte = '!'
	NULL_BYTE          typeByte = '_'
)
