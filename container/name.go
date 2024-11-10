package container

import (
	"io/ioutil"
	"os"
	"strings"
)

// GetId 获取容器id
// K8S: 取HOSTNAME，然后以 - 分隔取最后一段
// Docker: 取HOSTNAME
func GetId() string {
	hostname := GetHostName()

	nameParts := strings.Split(hostname, "-")
	partsLen := len(nameParts)
	if partsLen == 0 {
		return ""
	}

	return nameParts[partsLen-1]
}

// GetHostName 获取主机名称
func GetHostName() string {
	name := os.Getenv("HOSTNAME")
	if name != "" {
		return name
	}

	f, err := os.OpenFile("/etc/hostname", os.O_RDONLY, 0644)
	if err != nil {
		return ""
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}

	return string(bs)
}
