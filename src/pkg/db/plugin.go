package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/marmotedu/iam/pkg/log"
)

const (
	callBackBeforeName = "core:before"
	callBackAfterName  = "core:after"
	startTime          = "_start_time"
)

// TracePlugin defines gorm plugin used to trace sql.
type TracePlugin struct{}

// Name returns the name of trace plugin.
func (op *TracePlugin) Name() string {
	return "tracePlugin"
}

// Initialize initialize the trace plugin.
func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)

	return
}

var _ gorm.Plugin = &TracePlugin{}

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
}

func after(db *gorm.DB) {
	_ts, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}

	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}
	// sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	log.Infof("sql cost time: %fs", time.Since(ts).Seconds())
}

// migrateDatabase run auto migration for given models, will only add missing fields,
// won't delete/change current data.
// nolint:unused // may be reused in the feature, or just show a migrate usage.

// GORM 的 AutoMigrate 方法，只对新增的字段或索引进行变更，理论上是没有风险的。
// 就是模型增加了什么字段或者索引，就自动在表里面加了？
//func migrateDatabase(db *gorm.DB) error {
//	if err := db.AutoMigrate(&v1.User{}); err != nil {
//		return errors.Wrap(err, "migrate user model failed")
//	}
//	if err := db.AutoMigrate(&v1.Policy{}); err != nil {
//		return errors.Wrap(err, "migrate policy model failed")
//	}
//	if err := db.AutoMigrate(&v1.Secret{}); err != nil {
//		return errors.Wrap(err, "migrate secret model failed")
//	}
//
//	return nil
//}

// resetDatabase resets the database tables.
// nolint:unused,deadcode // may be reused in the feature, or just show a migrate usage.
//func resetDatabase(db *gorm.DB) error {
//	if err := cleanDatabase(db); err != nil {
//		return err
//	}
//	if err := migrateDatabase(db); err != nil {
//		return err
//	}
//
//	return nil
//}

// cleanDatabase tear downs the database tables.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
//func cleanDatabase(db *gorm.DB) error {
//	if err := db.Migrator().DropTable(&v1.User{}); err != nil {
//		return errors.Wrap(err, "drop user table failed")
//	}
//	if err := db.Migrator().DropTable(&v1.Policy{}); err != nil {
//		return errors.Wrap(err, "drop policy table failed")
//	}
//	if err := db.Migrator().DropTable(&v1.Secret{}); err != nil {
//		return errors.Wrap(err, "drop secret table failed")
//	}
//
//	return nil
//}
