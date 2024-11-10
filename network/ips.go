package network

import "net"

func GetIps(name string) ([]string, error) {
	ret := []string{}
	intrs, err := net.Interfaces()
	if nil != err {
		return nil, err
	}
	for _, intr := range intrs {
		if len(name) > 0 && name != intr.Name {
			continue
		}

		addrs, err := intr.Addrs()
		if nil != err {
			return nil, err
		}
		for _, addr := range addrs {
			ipNet, isValidIpNet := addr.(*net.IPNet)
			if isValidIpNet {
				if ipNet.IP.To4() != nil {
					ret = append(ret, ipNet.IP.String())
				}
			}
		}
	}
	return ret, nil
}
