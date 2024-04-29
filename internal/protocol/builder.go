package protocol

import "strconv"

type Builder []byte

func NewBuilder() Builder {
    return make([]byte, 0)
}

func (b Builder) AddArray(length int) Builder {
    b = append(b, ARRAY_BYTE)
    b = b.addLength(length)
    return b.addEnd()
}

func (b Builder) AddBulkString(str string) Builder {
    b = append(b, BULK_STRING_BYTE)
    b = b.addLength(len(str))
    b = b.addEnd()
    b = b.addString(str)
    return b.addEnd()
}

func (b Builder) AddSimpleString(str string) Builder {
    b = append(b, SIMPLE_STRING_BYTE)
    b = b.addString(str)
    return b.addEnd()
}

func (b Builder) addString(str string) Builder {
    for _, ch := range str {
        b = append(b, byte(ch))
    }
    return b
}

func (b Builder) addLength(length int) Builder {
    lenString := strconv.Itoa(length)
    for _, ch := range lenString {
        b = append(b, byte(ch))
    }
    return b
}

func (b Builder) addEnd() Builder {
    b = append(b, '\r')
    b = append(b, '\n')
    return b
}
