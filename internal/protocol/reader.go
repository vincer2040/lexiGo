package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/vincer2040/lexiGo/pkg/lexitypes"
)

type Reader struct {
	rd *bufio.Reader
}

func NewReader(rd *bufio.Reader) Reader {
	return Reader{rd: rd}
}

func (r *Reader) ReadReply() (*lexitypes.LexiType, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case ARRAY_BYTE:
		return r.readArray(line)
	case SIMPLE_STRING_BYTE:
		return r.readSimpleString(line), nil
	case BULK_STRING_BYTE:
		return r.readBulkString(line)
	case INT_BYTE:
		return r.readInt(line)
	case DBL_BYTE:
		return r.readDouble(line)
	case SIMPLE_ERROR_BYTE:
		return r.readSimpleError(line), nil
	case BULK_ERROR_BYTE:
		return r.readBulkError(line)
	}
	return nil, fmt.Errorf("unknown type byte %c", line[0])
}

func (r *Reader) readArray(line []byte) (*lexitypes.LexiType, error) {
	length := r.readLen(line)
	arr := make([]lexitypes.LexiType, 0)
	for i := 0; i < length; i++ {
		r, err := r.ReadReply()
		if err != nil {
			return nil, err
		}
		arr = append(arr, *r)
	}
	return &lexitypes.LexiType{DataType: lexitypes.Array, Data: lexitypes.LexiArray(arr)}, nil
}

func (r *Reader) readSimpleString(line []byte) *lexitypes.LexiType {
	buf := bytes.NewBufferString("")
	i := 1
	for line[i] != '\r' {
		buf.WriteByte(line[i])
		i++
	}
	return &lexitypes.LexiType{
		DataType: lexitypes.String,
		Data:     lexitypes.LexiString(buf.String()),
	}
}

func (r *Reader) readBulkString(line []byte) (*lexitypes.LexiType, error) {
	length := r.readLen(line)
	buf := make([]byte, length+2)
	_, err := io.ReadFull(r.rd, buf)
	if err != nil {
		return nil, err
	}
	s := bytes.NewBufferString("")
	for _, ch := range buf[:length] {
		s.WriteByte(ch)
	}
	return &lexitypes.LexiType{
		DataType: lexitypes.String,
		Data:     lexitypes.LexiString(s.String()),
	}, nil
}

func (r *Reader) readInt(line []byte) (*lexitypes.LexiType, error) {
	s := bytes.NewBufferString("")
	i := 1
	for line[i] != '\r' {
		s.WriteByte(line[i])
		i++
	}
	res, err := strconv.ParseInt(s.String(), 10, 64)
	if err != nil {
		return nil, err
	}
	return &lexitypes.LexiType{
		DataType: lexitypes.Int,
		Data:     lexitypes.LexiInt(res),
	}, nil
}

func (r *Reader) readDouble(line []byte) (*lexitypes.LexiType, error) {
	s := bytes.NewBufferString("")
	i := 1
	for line[i] != '\r' {
		s.WriteByte(line[i])
		i++
	}
	res, err := strconv.ParseFloat(s.String(), 64)
	if err != nil {
		return nil, err
	}
	return &lexitypes.LexiType{
		DataType: lexitypes.Double,
		Data:     lexitypes.LexiDouble(res),
	}, nil
}

func (r *Reader) readSimpleError(line []byte) *lexitypes.LexiType {
	buf := bytes.NewBufferString("")
	i := 1
	for line[i] != '\r' {
		buf.WriteByte(line[i])
		i++
	}
	return &lexitypes.LexiType{
		DataType: lexitypes.Error,
		Data:     lexitypes.LexiString(buf.String()),
	}
}

func (r *Reader) readBulkError(line []byte) (*lexitypes.LexiType, error) {
	length := r.readLen(line)
	buf := make([]byte, length+2)
	_, err := io.ReadFull(r.rd, buf)
	if err != nil {
		return nil, err
	}
	s := bytes.NewBufferString("")
	for _, ch := range buf[:length] {
		s.WriteByte(ch)
	}
	return &lexitypes.LexiType{
		DataType: lexitypes.Error,
		Data:     lexitypes.LexiString(s.String()),
	}, nil
}

func (r *Reader) readNull(line []byte) *lexitypes.LexiType {
	return &lexitypes.LexiType{DataType: lexitypes.Null}
}

func (r *Reader) readLen(line []byte) int {
	res := 0
	i := 1
	for line[i] != '\r' {
		res = (res * 10) + int(line[i]-'0')
		i++
	}
	return res
}

func (r *Reader) readLine() ([]byte, error) {
	b, err := r.rd.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	if len(b) <= 2 || b[len(b)-1] != '\n' || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("invalid reply: %q", b)
	}
	return b, nil
}
