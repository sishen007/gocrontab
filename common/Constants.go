package common

const (
	// 任务保存前缀
	JOB_SAVE_DIR = "/cron/jobs/"

	// 杀死任务前缀
	JOB_KILLER_DIR = "/cron/killer/"
	// 所路径
	JOB_LOCK_DIR = "/cron/lock/"

	// 保存任务事件
	JOB_EVENT_SAVE = 1
	// 删除任务事件
	JOB_EVENT_DELETE = 2
)
