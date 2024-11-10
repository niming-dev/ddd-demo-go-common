package network

import (
	"runtime"
	"testing"
)

func Test_Ips(t *testing.T) {
	// get all ips
	ret, err := GetIps("")
	if nil != err {
		t.Fatal(err)
	}
	t.Log(ret)

	var nicName string
	switch runtime.GOOS {
	case `darwin`:
		nicName = "en0"
	default:
		nicName = "eth0"
	}

	ret, err = GetIps(nicName)
	if nil != err {
		t.Fatal(err)
	}
	t.Log(ret)
}
