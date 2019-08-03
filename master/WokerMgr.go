package master

import (
	"context"
	"github.com/sishen007/gocrontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_workerMgr *WorkerMgr
)

// 获取在线的worker列表
func (workerMgr *WorkerMgr) ListWorkers() (workerArr []string, err error) {
	var (
		getResp  *clientv3.GetResponse
		kv       *mvccpb.KeyValue
		workerIP string
	)
	// 初始化workerArr
	workerArr = make([]string, 0)
	// 获取目录下的所有kvs
	if getResp, err = workerMgr.kv.Get(context.TODO(), common.JOB_WORKER_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	for _, kv = range getResp.Kvs {
		workerIP = common.ExtractWorkerIPName(string(kv.Key))
		workerArr = append(workerArr, workerIP)
	}
	return
}
func InitWorkerMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
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
	// 获取KV和Lease的API子集
	// 用于读取etcd键值对
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 赋值单例
	G_workerMgr = &WorkerMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}
