package worker

import (
	"context"
	"fmt"
	"github.com/sishen007/gocrontab/common"
	"go.etcd.io/etcd/clientv3"
	"net"
	"time"
)

type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP string // 本机IP
}

var (
	G_register *Register
)

// 获取ip地址
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

// 服务注册租约
func (register *Register) keepOnline() {
	var (
		registerKey    string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		ctx            context.Context
		cancelFunc     context.CancelFunc
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		err            error
	)
	// 注册路径
	registerKey = common.JOB_WORKER_DIR + register.localIP

	for {
		cancelFunc = nil
		// 申请一个10s的租约
		if leaseGrantResp, err = register.lease.Grant(context.TODO(), 5); err != nil {
			goto RETRY
		}
		// 获取租约ID
		leaseId = leaseGrantResp.ID
		// 租约续租
		ctx, cancelFunc = context.WithCancel(context.TODO())
		//defer cancelFunc()
		//defer register.lease.Revoke(context.TODO(), leaseId)
		if keepRespChan, err = register.lease.KeepAlive(ctx, leaseId); err != nil {
			goto RETRY
		}
		// 将租约注册到etcd
		if _, err = register.kv.Put(context.TODO(), registerKey, "", clientv3.WithLease(leaseId)); err != nil {
			goto RETRY
		}
		// 处理续租应答
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					goto RETRY
				}
			}
		}
	RETRY:
		if cancelFunc != nil {
			cancelFunc()
		}
		time.Sleep(1 * time.Second)
	}
}

// 初始化配置
func InitRegister() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
	)
	// 初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints, // 集群列表
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}
	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	// 获取ip地址
	if localIp, err = getLocalIP(); err != nil {
		return
	}
	// 获取KV和Lease的API子集
	// 用于读取etcd键值对
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 赋值单例
	G_register = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIp,
	}
	fmt.Println(G_register)
	// 服务注册
	go G_register.keepOnline()
	return
}
