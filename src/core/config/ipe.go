package config

import (
	"errors"
	"strings"

	"github.com/eaglexiang/go-settings"
)

const (
	ipeSplitSig = ","
	defaultPort = "8080"
)

// IPPorts 每个ipPorts由一个IP与多个Port组成，ports的长度至少为1
type IPPorts struct {
	IP    string
	Ports []string
}

func (ip *IPPorts) setDefaultPort(port string) {
	if len(ip.Ports) == 0 || ip.Ports[0] == "" {
		ip.Ports = []string{port}
	}
}

func (ip *IPPorts) addPort(port string) {
	for _, p := range ip.Ports {
		if p == port {
			return
		}
	}

	ip.Ports = append(ip.Ports, port)
}

func (ip IPPorts) toString() (result string) {
	for _, port := range ip.Ports {
		ipe := ip.IP + ":" + port

		if result != "" {
			result += ipeSplitSig
		}
		result += ipe
	}
	return
}

func parseIPPortsSlice(src string) []*IPPorts {
	// map[ipPorts.ip] ipPorts
	ipPortsMap := make(map[string]*IPPorts)

	ipes := strings.Split(src, ipeSplitSig)
	for _, ipe := range ipes {
		ip, port, err := getIPPort(ipe)
		if err != nil {
			panic(err)
		}

		if ipports, ok := ipPortsMap[ip]; ok {
			ipports.addPort(port)
		} else {
			ipPortsMap[ip] = &IPPorts{
				IP:    ip,
				Ports: []string{port},
			}
		}
	}

	// 设置默认端口
	for _, ip := range ipPortsMap {
		ip.setDefaultPort(defaultPort)
	}

	ipPortsSlice := []*IPPorts{}
	for _, ip := range ipPortsMap {
		ipPortsSlice = append(ipPortsSlice, ip)
	}

	return ipPortsSlice
}

func getIPPort(ipe string) (ip, port string, err error) {
	if strings.HasPrefix(ipe, "[") {
		ip, port, err = getIPPortFromIPv6IPE(ipe)
	} else {
		ip, port, err = getIPPortFromIPv4IPE(ipe)
	}
	return
}

func getIPPortFromIPv4IPE(ipe string) (ip, port string, err error) {
	ipeSlice := strings.Split(ipe, ":")

	ip = ipeSlice[0]

	if len(ipeSlice) > 1 {
		port = ipeSlice[1]
	}

	return
}

func getIPPortFromIPv6IPE(ipe string) (ip, port string, err error) {
	ipeSlice := strings.Split(ipe, ":")

	if len(ipeSlice) == 6 {
		if !strings.HasSuffix(ipe, "]") {
			err = errors.New("ipv6 need ]")
		} else {
			ip = ipe
		}
	} else if len(ipeSlice) == 7 {
		ip = strings.Join(ipeSlice[:6], ":")
		port = ipeSlice[6]
	} else {
		err = errors.New("invalid ipe")
	}

	return
}

// initListens ipes的示例：192.168.0.1:8080,192.168.0.1:8081
func initListens() {
	ListenIPEs = parseIPPortsSlice(settings.Get("listen"))
}

func initRelays() {
	RelayIPEs = parseIPPortsSlice(settings.Get("relay"))
}

func RelayIPE() string {
	relayIPPorts := RelayIPEs[0]
	relayIPE := relayIPPorts.IP + ":" + relayIPPorts.Ports[0]
	return relayIPE
}
