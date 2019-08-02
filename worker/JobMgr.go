package worker

import (
	"context"
	"fmt"
	"github.com/sishen007/gocrontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	// 单例
	G_jobMgr *JobMgr
)

// 监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchRespChan      clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)
	// 1. 获取/cron/jobs/目录下的所有任务,并获知当前集群的revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	// 遍历当前有哪些任务
	for _, kvpair = range getResp.Kvs {
		// 反序列化json到job
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			// 将job同步给scheduler
			G_scheduler.PushJobEvent(jobEvent)
		}
	}
	// 2. 从当前revision向后监听变化事件
	go func() {
		// 监听协程
		watchStartRevision = getResp.Header.Revision + 1
		watchRespChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix(), clientv3.WithRev(watchStartRevision))
		//处理kv变化事件
		for watchResp = range watchRespChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 任务保存事件
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// 构造一个更新Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE: // 任务删除事件
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					// 构造一个删除Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, &common.Job{Name: jobName})
				}
				fmt.Println(*jobEvent)
				// 推送scheduler
				G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}

// 监听进程杀死
func (jobMgr *JobMgr) watchKiller() (err error) {
	var (
		job           *common.Job
		watchRespChan clientv3.WatchChan
		watchResp     clientv3.WatchResponse
		watchEvent    *clientv3.Event
		jobName       string
		jobEvent      *common.JobEvent
	)
	// 监听/cron/killer目录
	go func() {
		// 监听协程-监听/cron/killer目录的变化
		watchRespChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR, clientv3.WithPrefix())
		//处理监听事件
		for watchResp = range watchRespChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 杀死任务事件
					// 构造一个更新Event
					jobName = common.ExtractKillerName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JON_EVENT_KILL, job)

				case mvccpb.DELETE: //killer标记过期,被自动删除
				}
				// 推送scheduler
				G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}

// 初始化管理器
func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
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
	watcher = clientv3.Watcher(client)
	// 赋值单例
	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}
	// 启动任务监听
	if err = G_jobMgr.watchJobs(); err != nil {
		return
	}
	// 启动任务killer监听
	if err = G_jobMgr.watchKiller(); err != nil {
		return
	}
	return
}

// 创建分布式锁
func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	// 返回一把锁
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}
