package main

import (
	"fmt"
	"github.com/sishen007/gocrontab/common"
	"net"
)

type T struct {
	Name string
	Age  int64
}

func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	fmt.Println()
	// 取第一个非IO非网卡IP
	for _, addr = range addrs {
		// 这个网卡地址是ip地址:ipv4 / ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}
	err = common.ERR_NO_LOCAL_IP_FOUND
	return
}
func main() {
	// 获取ip地址
	var localIp string
	var err error
	if localIp, err = getLocalIP(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(localIp)

	//var (
	//	oldJobVal common.Job
	//	err       error
	//)
	//str := `{"name":"job1","command":"echo hello","cronExpr":"* * * * *"}`
	//if err = json.Unmarshal([]byte(str), &oldJobVal); err != nil {
	//	return
	//}
	//fmt.Println("值:", oldJobVal)

	// test Time
	//t := time.Now().Add(10 * time.Second).Sub(time.Now())
	//fmt.Printf("%#v", t)

	//t := &T{"张三", 20}
	//i := 1
	//modifyT(t, &i)
	//go getT(t, &i)
	//time.Sleep(2 * time.Second)
	//fmt.Println("t:", t)
	//fmt.Println("i:", i)

}
func modifyT(t *T, i *int) {
	go func() {
		t.Name = "李四"
		*i = 2
	}()
}
func getT(t *T, i *int) {
	fmt.Println("getT t:", t)
	fmt.Println("getT i:", *i)
}
