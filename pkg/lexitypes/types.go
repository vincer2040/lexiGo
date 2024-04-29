package lexitypes

type LexiDataType int

const (
	Null LexiDataType = iota
	String
	Int
	Double
	Error
	Array
)

type LexiType struct {
	DataType LexiDataType
	Data     LexiData
}

type LexiData interface {
	lexiData()
}

type LexiString string

func (s LexiString) lexiData() {}

type LexiInt int64

func (i LexiInt) lexiData() {}

type LexiDouble float64

func (f LexiDouble) lexiData() {}

type LexiArray []LexiType

func (a LexiArray) lexiData() {}
