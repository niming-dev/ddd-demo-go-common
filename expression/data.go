package expression

import (
	"fmt"
	"strconv"
	"strings"
)

const INVALID_KIND int16 = 0
const INT_KIND int16 = 1
const STRING_KIND int16 = 2
const BOOL_KIND int16 = 3
const STRUCT_KIND int16 = 4

var KIND_NAME []string = []string{
	"invalidKind",
	"int",
	"string",
	"bool",
	"struct",
}

func kindString(kind int16) string {
	if int(kind) >= len(KIND_NAME) {
		return KIND_NAME[0]
	}
	return KIND_NAME[kind]
}

const INVALID_DATA_STRING string = "[Invalid Data]"

var ErrDataTypeNotMatch = fmt.Errorf("data type not match")
var ErrZeroDiv = fmt.Errorf("zero div")

type Data struct {
	iVal     int64
	sVal     string
	bVal     bool
	mapVal   map[string]*Data
	iValKind int16 // 0 InvalidData, 1 int, 2 string, 3, bool, 4, struct
}

func (d *Data) Evaluate(ec ExecuteContext) (*Data, error) {
	return d, nil
}

func (d *Data) IsSameKind(o *Data) bool {
	return d.iValKind != INVALID_KIND && d.iValKind == o.iValKind
}

func (d *Data) IsInvalid() bool {
	return d.iValKind == INVALID_KIND
}

func (d *Data) Val() interface{} {
	switch d.iValKind {
	case INT_KIND:
		return d.iVal
	case STRING_KIND:
		return d.sVal
	case BOOL_KIND:
		return d.bVal
	case STRUCT_KIND:
		return d.mapVal
	default:
		return INVALID_DATA_STRING
	}
}

func (d *Data) IsInt() bool {
	return d.iValKind == INT_KIND
}
func (d *Data) Int() int64 {
	if d.iValKind != INT_KIND {
		return 0
	}
	return d.iVal
}

func (d *Data) IsBool() bool {
	return d.iValKind == BOOL_KIND
}
func (d *Data) Bool() bool {
	if d.iValKind != BOOL_KIND {
		return false
	}
	return d.bVal
}

func (d *Data) IsString() bool {
	return d.iValKind == STRING_KIND
}
func (d *Data) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: Data, type: %s, value: %v", deep*4, "", kindString(d.iValKind), d)
	return ret
}

func (d *Data) String() string {
	switch d.iValKind {
	case INT_KIND:
		return fmt.Sprintf("%v", d.iVal)
	case STRING_KIND:
		return d.sVal
	case BOOL_KIND:
		return fmt.Sprintf("%v", d.bVal)
	case INVALID_KIND:
		return INVALID_DATA_STRING
	default:
		return INVALID_DATA_STRING
	}
}

func (d *Data) IsStruct() bool {
	return d.iValKind == STRUCT_KIND
}

func (d *Data) Struct() map[string]*Data {
	if d.iValKind != STRUCT_KIND {
		return nil
	}
	return d.mapVal
}

func (d *Data) Neg() (*Data, error) {
	if d.iValKind != INT_KIND {
		return nil, ErrSyntaxError
	}
	return &Data{iValKind: INT_KIND, iVal: 0 - d.iVal}, nil
}

func (d *Data) ToInt() (*Data, error) {
	switch d.iValKind {
	case INT_KIND:
		return d, nil
	case BOOL_KIND:
		val := int64(0)
		if d.bVal {
			val = 1
		}
		return &Data{iValKind: INT_KIND, iVal: val}, nil
	case STRING_KIND:
		val, err := strconv.ParseInt(d.sVal, 10, 64)
		if nil != err {
			return nil, err
		}
		return &Data{iValKind: INT_KIND, iVal: val}, nil
	}
	return nil, ErrInvalidConvert
}

func (d *Data) ToBool() (*Data, error) {
	switch d.iValKind {
	case INT_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d.iVal != 0}, nil
	case BOOL_KIND:
		return d, nil
	case STRING_KIND:
		val, err := strconv.ParseBool(d.sVal)
		if nil != err {
			return nil, err
		}
		return &Data{iValKind: BOOL_KIND, bVal: val}, nil
	}
	return nil, ErrInvalidConvert
}

func NewStruct(m map[string]*Data) *Data {
	return &Data{
		iValKind: STRUCT_KIND,
		mapVal:   m,
	}
}

func NewInt(i int64) *Data {
	return &Data{iValKind: INT_KIND, iVal: i}
}

func NewString(s string) *Data {
	return &Data{iValKind: STRING_KIND, sVal: s}
}

func NewBool(b bool) *Data {
	return &Data{iValKind: BOOL_KIND, bVal: b}
}

func PrepareData(o1, o2 *Data) (*Data, *Data, error) {
	var d1, d2 *Data
	var err error
	if o1.iValKind == INT_KIND || o2.iValKind == INT_KIND {
		d1, err = o1.ToInt()
		if nil != err {
			return nil, nil, err
		}
		d2, err = o2.ToInt()
		if nil != err {
			return nil, nil, err
		}
		return d1, d2, nil
	} else if o1.iValKind == BOOL_KIND || o2.iValKind == BOOL_KIND {
		d1, err = o1.ToBool()
		if nil != err {
			return nil, nil, err
		}
		d2, err = o2.ToBool()
		if nil != err {
			return nil, nil, err
		}
		return d1, d2, nil
	} else {
		return o1, o2, nil
	}
}

func Add(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}

	if !d1.IsSameKind(d2) || d1.iValKind == BOOL_KIND {
		return nil, ErrDataTypeNotMatch
	}
	switch d1.iValKind {
	case INT_KIND:
		return &Data{iValKind: INT_KIND, iVal: d1.iVal + d2.iVal}, nil
	case STRING_KIND:
		return &Data{iValKind: STRING_KIND, sVal: d1.sVal + d2.sVal}, nil
	default:
		// never been here
		panic("Invalid ValKind")
	}
}

func Sub(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}

	if !d1.IsSameKind(d2) || d1.iValKind != INT_KIND {
		return nil, ErrDataTypeNotMatch
	}
	return &Data{iValKind: INT_KIND, iVal: d1.iVal - d2.iVal}, nil
}

func Mul(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}

	if !d1.IsSameKind(d2) || d1.iValKind != INT_KIND {
		return nil, ErrDataTypeNotMatch
	}
	return &Data{iValKind: INT_KIND, iVal: d1.iVal * d2.iVal}, nil
}

func Div(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}
	if !d1.IsSameKind(d2) || d1.iValKind != INT_KIND {
		return nil, ErrDataTypeNotMatch
	}
	if d2.iVal == 0 {
		return nil, ErrZeroDiv
	}
	return &Data{iValKind: INT_KIND, iVal: d1.iVal / d2.iVal}, nil
}

func Mod(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}
	if !d1.IsSameKind(d2) || d1.iValKind != INT_KIND {
		return nil, ErrDataTypeNotMatch
	}
	if d2.iVal == 0 {
		return nil, ErrZeroDiv
	}
	return &Data{iValKind: INT_KIND, iVal: d1.iVal % d2.iVal}, nil
}

func Greate(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}
	if !d1.IsSameKind(d2) || d1.iValKind == BOOL_KIND {
		return nil, ErrDataTypeNotMatch
	}
	switch d1.iValKind {
	case INT_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d1.iVal > d2.iVal}, nil
	case STRING_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: strings.Compare(d1.sVal, d2.sVal) == 1}, nil
	default:
		// never been here
		panic("Invalid kind")
	}
}

func Less(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}
	if !d1.IsSameKind(d2) || d1.iValKind == BOOL_KIND {
		return nil, ErrDataTypeNotMatch
	}
	switch d1.iValKind {
	case INT_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d1.iVal < d2.iVal}, nil
	case STRING_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: strings.Compare(d1.sVal, d2.sVal) == -1}, nil
	default:
		// never been here
		panic("Invalid kind")
	}
}

func GreateEqual(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}
	if !d1.IsSameKind(d2) || d1.iValKind == BOOL_KIND {
		return nil, ErrDataTypeNotMatch
	}
	switch d1.iValKind {
	case INT_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d1.iVal >= d2.iVal}, nil
	case STRING_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: strings.Compare(d1.sVal, d2.sVal) >= 0}, nil
	default:
		// never
		panic("Invalid kind")
	}
}

func LessEqual(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}
	if !d1.IsSameKind(d2) || d1.iValKind == BOOL_KIND {
		return nil, ErrDataTypeNotMatch
	}
	switch d1.iValKind {
	case INT_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d1.iVal <= d2.iVal}, nil
	case STRING_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: strings.Compare(d1.sVal, d2.sVal) <= 0}, nil
	default:
		// never
		panic("Invalid kind")
	}
}

func Equal(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}
	if !d1.IsSameKind(d2) {
		return nil, ErrDataTypeNotMatch
	}
	switch d1.iValKind {
	case INT_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d1.iVal == d2.iVal}, nil
	case STRING_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: strings.Compare(d1.sVal, d2.sVal) == 0}, nil
	case BOOL_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d1.bVal == d2.bVal}, nil
	default:
		panic("Invalid kind")
	}
}

func NotEqual(o1, o2 *Data) (*Data, error) {
	d1, d2, err := PrepareData(o1, o2)
	if nil != err {
		return nil, err
	}
	if !d1.IsSameKind(d2) {
		return nil, ErrDataTypeNotMatch
	}
	switch d1.iValKind {
	case INT_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d1.iVal != d2.iVal}, nil
	case STRING_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: strings.Compare(d1.sVal, d2.sVal) != 0}, nil
	case BOOL_KIND:
		return &Data{iValKind: BOOL_KIND, bVal: d1.bVal != d2.bVal}, nil
	default:
		panic("Invalid kind")
	}
}
