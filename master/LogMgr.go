package master

import (
	"context"
	"github.com/sishen007/gocrontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogMgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
	var (
		client *mongo.Client
	)

	// 建立mongodb连接
	// 建立连接
	if client, err = mongo.NewClient(options.Client().ApplyURI(G_config.MongodbUri)); err != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond)
	if err = client.Connect(ctx); err != nil {
		return
	}

	// 选择db和collection
	G_logMgr = &LogMgr{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

// 查询日志
func (logMgr *LogMgr) ListLog(name string, skip, limit int) (logArr []*common.JobLog, err error) {
	var (
		filter      *common.JobLogFilter
		logSort     *common.SortLogByStartTime
		cursor      *mongo.Cursor
		record      *common.JobLog
		findOptions *options.FindOptions
	)
	// 初始化logArr(为了len(logArr = 0))
	logArr = make([]*common.JobLog, 0)
	// 过滤条件
	filter = &common.JobLogFilter{
		JobName: name,
	}
	// 排序条件
	logSort = &common.SortLogByStartTime{
		SortOrder: -1,
	}
	findOptions = options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(logSort)
	if cursor, err = logMgr.logCollection.Find(context.TODO(), filter, findOptions); err != nil {
		return
	}

	// 遍历结果集
	for cursor.Next(context.TODO()) {
		record = &common.JobLog{}
		// 反序列化到bson到对象
		if err = cursor.Decode(record); err != nil {
			continue
		}
		logArr = append(logArr, record)
	}
	return
}
