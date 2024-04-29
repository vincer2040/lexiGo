package protocol

import (
	"bufio"
	"bytes"
	"github.com/vincer2040/lexiGo/pkg/lexitypes"
	"testing"
)

func TestSimpleStrings(t *testing.T) {
	tests := []struct {
		input string
		exp   string
	}{
		{
			input: "+OK\r\n",
			exp:   "OK",
		},
		{
			input: "+PONG\r\n",
			exp:   "PONG",
		},
	}
	for _, test := range tests {
		rd := bufio.NewReader(bytes.NewReader([]byte(test.input)))
		reader := NewReader(rd)
		res, err := reader.ReadReply()
		if err != nil {
			t.Fatalf("expected: %s, got error %+v\n", test.exp, err)
		}
		if res.DataType != lexitypes.String {
			t.Fatalf("expected String, got %+v\n", res.DataType)
		}
		s := string(res.Data.(lexitypes.LexiString))
		if s != test.exp {
			t.Fatalf("expected: %s, got %s\n", test.exp, s)
		}
	}
}

func TestBulkStrings(t *testing.T) {
	tests := []struct {
		input string
		exp   string
	}{
		{
			input: "$3\r\nfoo\r\n",
			exp:   "foo",
		},
		{
			input: "$10\r\nfoo\r\nfoo\r\n\r\n",
			exp:   "foo\r\nfoo\r\n",
		},
	}
	for _, test := range tests {
		rd := bufio.NewReader(bytes.NewReader([]byte(test.input)))
		reader := NewReader(rd)
		res, err := reader.ReadReply()
		if err != nil {
			t.Fatalf("expected: %s, got error %+v\n", test.exp, err)
		}
		if res.DataType != lexitypes.String {
			t.Fatalf("expected String, got %+v\n", res.DataType)
		}
		s := string(res.Data.(lexitypes.LexiString))
		if s != test.exp {
			t.Fatalf("expected: %s, got %s\n", test.exp, s)
		}
	}
}

func TestIntegers(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		{
			input: ":123\r\n",
			exp:   123,
		},
		{
			input: ":12345\r\n",
			exp:   12345,
		},
	}
	for _, test := range tests {
		rd := bufio.NewReader(bytes.NewReader([]byte(test.input)))
		reader := NewReader(rd)
		res, err := reader.ReadReply()
		if err != nil {
			t.Fatalf("expected: %d, got error %+v\n", test.exp, err)
		}
		if res.DataType != lexitypes.Int {
			t.Fatalf("expected String, got %+v\n", res.DataType)
		}
		s := int64(res.Data.(lexitypes.LexiInt))
		if s != test.exp {
			t.Fatalf("expected: %d, got %d\n", test.exp, s)
		}
	}
}

func TestSimpleErrors(t *testing.T) {
	tests := []struct {
		input string
		exp   string
	}{
		{
			input: "-ERR\r\n",
			exp:   "ERR",
		},
		{
			input: "-ENOACCESS\r\n",
			exp:   "ENOACCESS",
		},
	}
	for _, test := range tests {
		rd := bufio.NewReader(bytes.NewReader([]byte(test.input)))
		reader := NewReader(rd)
		res, err := reader.ReadReply()
		if err != nil {
			t.Fatalf("expected: %s, got error %+v\n", test.exp, err)
		}
		if res.DataType != lexitypes.Error {
			t.Fatalf("expected Error, got %+v\n", res.DataType)
		}
		s := string(res.Data.(lexitypes.LexiString))
		if s != test.exp {
			t.Fatalf("expected: %s, got %s\n", test.exp, s)
		}
	}
}

func TestBulkErrors(t *testing.T) {
	tests := []struct {
		input string
		exp   string
	}{
		{
			input: "!3\r\nfoo\r\n",
			exp:   "foo",
		},
		{
			input: "!10\r\nfoo\r\nfoo\r\n\r\n",
			exp:   "foo\r\nfoo\r\n",
		},
	}
	for _, test := range tests {
		rd := bufio.NewReader(bytes.NewReader([]byte(test.input)))
		reader := NewReader(rd)
		res, err := reader.ReadReply()
		if err != nil {
			t.Fatalf("expected: %s, got error %+v\n", test.exp, err)
		}
		if res.DataType != lexitypes.Error {
			t.Fatalf("expected error, got %+v\n", res.DataType)
		}
		s := string(res.Data.(lexitypes.LexiString))
		if s != test.exp {
			t.Fatalf("expected: %s, got %s\n", test.exp, s)
		}
	}

}
