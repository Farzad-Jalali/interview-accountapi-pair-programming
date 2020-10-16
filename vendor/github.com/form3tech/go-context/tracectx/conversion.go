package tracectx

import "github.com/form3tech/go-context/tracectx/internal"

type parser func(val string) (interface{}, error)
type stringer func(val interface{}) (string, error)

// BidirectionalConversion defines two functions; one for parsing from string to typed value, and stringer from typed value to string
type BidirectionalConversion struct {
	Parser   parser
	Stringer stringer
}

var (
	StringConverter = BidirectionalConversion{
		Parser:   internal.StringParser,
		Stringer: internal.Stringer,
	}
	TimeConverter = BidirectionalConversion{
		Parser:   internal.TimeParser,
		Stringer: internal.TimeStringer,
	}
	DurationConverter = BidirectionalConversion{
		Parser:   internal.DurationParser,
		Stringer: internal.DurationStringer,
	}
	IntConverter = BidirectionalConversion{
		Parser:   internal.IntParser,
		Stringer: internal.IntStringer,
	}
	Int32Converter = BidirectionalConversion{
		Parser:   internal.Int32Parser,
		Stringer: internal.Int32Stringer,
	}
	Int64Converter = BidirectionalConversion{
		Parser:   internal.Int64Parser,
		Stringer: internal.Int64Stringer,
	}
)
