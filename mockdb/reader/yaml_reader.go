package reader

import (
	"database/sql/driver"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlReader struct {
	filename string
}

func NewYamlReader(filename string) *YamlReader {
	return &YamlReader{filename: filename}
}

func (y YamlReader) Read() (values []map[string]driver.Value, err error) {
	f, err := os.Open(y.filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	r := yaml.NewDecoder(f)
	err = r.Decode(&values)
	return values, err
}
