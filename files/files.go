package files

import (
	"io/ioutil"
	"path"
)

// GetAllFile 获取指定目录下，对应后缀名的所有文件，支持递归
func GetAllFile(filePath string, suffix ...string) (files []string, err error) {
	suffixMap := make(map[string]struct{}, len(suffix))
	for _, v := range suffix {
		suffixMap[v] = struct{}{}
	}

	rd, err := ioutil.ReadDir(filePath)
	if err != nil {
		return nil, err
	}

	for _, fi := range rd {
		// 文件，继续递归
		if fi.IsDir() {
			fs, err := GetAllFile(filePath+"/"+fi.Name(), suffix...)
			if err != nil {
				return nil, err
			}
			files = append(files, fs...)
			continue
		}

		// 以下处理的是文件

		// 如果没有传入后缀，则所有文件都符合要求
		if len(suffixMap) == 0 {
			fullName := filePath + "/" + fi.Name()
			files = append(files, fullName)
		}

		// 后缀符合
		if _, ok := suffixMap[path.Ext(fi.Name())]; ok {
			fullName := filePath + "/" + fi.Name()
			files = append(files, fullName)
		}
	}

	return files, nil
}

// GetAllDir 获取路径下所有的子目录，支持递归
func GetAllDir(path string) ([]string, error) {
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, fi := range rd {
		// 文件，继续递归
		if fi.IsDir() {
			dirs = append(dirs, path+"/"+fi.Name())

			ds, err := GetAllDir(path + "/" + fi.Name())
			if err != nil {
				return nil, err
			}
			dirs = append(dirs, ds...)
		}
	}

	return dirs, nil
}
