package worker

import (
	"context"
	"github.com/sishen007/gocrontab/common"
	"os/exec"
	"runtime"
	"time"
)

// 任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

// 执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
		)
		// 任务结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		// 初始化分布式锁
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)
		// 记录开始执行时间
		result.StartTime = time.Now()

		// 首先获取分布式锁
		// 任务执行完成后,释放锁
		err = jobLock.TryLock()
		defer jobLock.Unlock()
		if err != nil { // 上锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else {
			// 上锁成功,重置任务启动时间
			result.StartTime = time.Now()

			// 执行shell命令
			if runtime.GOOS == "windows" {
				cmd = exec.CommandContext(context.TODO(), "cmd", "/C", info.Job.Command)
			} else {
				cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
			}
			// 执行并捕获输出
			output, err = cmd.CombinedOutput()
			// 记录结束执行时间
			result.EndTime = time.Now()
			result.Err = err
			result.Output = output
		}
		// 任务执行完成后,把执行结果返回给Scheduler,Scheduler会从executingTable中删除执行记录
		G_scheduler.PushJobResult(result)
	}()
}

// 初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
