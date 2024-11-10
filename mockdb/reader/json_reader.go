package reader

import (
	"database/sql/driver"
	"os"

	jsoniter "github.com/json-iterator/go"
)

type JsonReader struct {
	filename string
}

func NewJsonReader(filename string) *JsonReader {
	return &JsonReader{filename: filename}
}

func (csv JsonReader) Read() (values []map[string]driver.Value, err error) {
	f, err := os.Open(csv.filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	r := jsoniter.NewDecoder(f)
	err = r.Decode(&values)
	if err != nil {
		return nil, err
	}

	return values, nil
}
