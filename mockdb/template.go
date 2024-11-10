package mockdb

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/niming-dev/ddd-demo/go-common/files"
)

// NewTemplate 新建模板实例
func NewTemplate(path string, suffix ...string) (*template.Template, error) {
	fs, err := files.GetAllFile(path, suffix...)
	if err != nil {
		return nil, err
	}
	tpl := template.New("template")
	tpl = tpl.Funcs(GetTplFuncMap(tpl))
	return parseFiles(tpl, path, fs...)
}

func getTplName(tplPath, fileName string) string {
	if strings.Index(fileName, tplPath) != 0 {
		return fileName
	}
	tpl := strings.Replace(fileName, tplPath, "", 1)
	if tpl[:1] == "/" {
		tpl = tpl[1:]
	}
	return tpl
}

func mergeFunc(src, dst template.FuncMap) template.FuncMap {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func readFile(tplPath, file string) (name string, b []byte, err error) {
	name = getTplName(tplPath, file)
	b, err = os.ReadFile(file)
	return
}

func parseFiles(t *template.Template, tplPath string, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}

	for _, filename := range filenames {
		name, bytes, err := readFile(tplPath, filename)
		if err != nil {
			return nil, err
		}

		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}

		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(string(bytes))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// GetTplFuncMap 获取模板函数
func GetTplFuncMap(t *template.Template) template.FuncMap {
	add := template.FuncMap{
		"join":    strings.Join,
		"include": include(t),
	}
	return mergeFunc(sprig.TxtFuncMap(), add)
}

func include(t *template.Template) func(name string, data interface{}) (string, error) {
	return func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		err := t.ExecuteTemplate(&buf, name, data)
		return buf.String(), err
	}
}
