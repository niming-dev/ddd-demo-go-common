package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/niming-dev/ddd-demo/go-common/expression"
)

type fetch struct{}

func (fetch) Name() string {
	return "fetch"
}

func GetStringValue(d *expression.Data, defaultValue string) string {
	if nil == d {
		return defaultValue
	}
	ret := d.String()
	if ret == expression.INVALID_DATA_STRING {
		return defaultValue
	}
	return ret
}

func DataToValues(d *expression.Data) (url.Values, error) {
	if !d.IsStruct() {
		return nil, expression.ErrDataTypeNotMatch
	}
	mapVals := d.Struct()
	vals := url.Values{}
	for k, v := range mapVals {
		vals[k] = []string{v.String()}
	}
	return vals, nil
}

func DataToJson(d *expression.Data) ([]byte, error) {
	if !d.IsStruct() {
		return nil, expression.ErrDataTypeNotMatch
	}
	mapVals := d.Struct()

	mapJson := map[string]interface{}{}
	for k, v := range mapVals {
		mapJson[k] = v.Val()
	}
	return json.Marshal(mapJson)
}

func (fetch) Call(ctx expression.ExecuteContext, args []*expression.Data) (*expression.Data, error) {
	if len(args) != 4 {
		return nil, expression.ErrMissArgument
	}
	method := GetStringValue(args[0], "GET")
	requestUrl := GetStringValue(args[1], "")
	encoding := GetStringValue(args[2], "form")

	vals, err := DataToValues(args[3])
	if nil != err {
		return nil, err
	}
	jsonBytes, err := DataToJson(args[3])
	if nil != err {
		return nil, err
	}

	var resp *http.Response

	switch method {
	case "GET":
		queryString := vals.Encode()
		resp, err = http.Get(requestUrl + "?" + queryString)
	case "POST":
		// application/x-www-form-urlencoded     application/json
		if encoding == "form" {
			resp, err = http.PostForm(requestUrl, vals)
		} else if encoding == "json" {
			resp, err = http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonBytes))
		} else {
			return nil, expression.ErrUnsupport
		}
	default:
		return nil, expression.ErrUnsupport
	}

	if nil != err {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got response status: %s", resp.Status)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return expression.NewString(string(respBytes)), nil
}

func init() {
	Register(&fetch{})
}
