package master

import (
	"encoding/json"
	"github.com/sishen007/gocrontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

// 任务的http接口
type ApiServer struct {
	httpServer *http.Server
}

var G_apiServer *ApiServer

// 保存任务接口至etcd
// POST job = {"name":"job1","command":"echo hello","cronExpr":"* * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	// 解析POST表单
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 获取表单中的job字段
	postJob = req.PostForm.Get("job")
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	// 保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	// 返回正常应答{errno:0,msg:"",data:{}}
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 删除任务接口
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	// 解析POST表单
	var (
		err    error
		bytes  []byte
		name   string
		oldJob *common.Job
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	// 获取删除的键名
	// 获取表单中的job字段
	name = req.PostForm.Get("name")

	// 删除操作
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	// 返回正常应答{errno:0,msg:"",data:{}}
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 获取任务列表
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	// 解析POST表单
	var (
		err     error
		bytes   []byte
		jobList []*common.Job
	)
	// 删除操作
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}
	// 返回正常应答{errno:0,msg:"",data:{}}
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 杀死任务
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	// 解析POST表单
	var (
		err   error
		bytes []byte
		name  string
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	// 获取删除的键名
	// 获取表单中的job字段
	name = req.PostForm.Get("name")

	// 删除操作
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}
	// 返回正常应答{errno:0,msg:"",data:{}}
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 初始化服务
func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir     // 静态文件根目录
		staticHandler http.Handler // 静态文件的Http回调
	)
	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	// 静态文件目录
	staticDir = http.Dir(G_config.StaticDir)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	// 启动tcp监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}
	// 创建http服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		Handler:      mux,
	}
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}
	// 启动了服务端
	go httpServer.Serve(listener)

	return
}
