package db

import (
	"fmt"

	"github.com/micro-plat/lib4go/db"
)

func ProcessingTask(db db.IDB, taskID int64) error {
	imap := map[string]interface{}{
		"task_id": taskID,
	}
	row, _, _, err := db.Execute(sqlProcessingTask, imap)
	if err != nil || row != 1 {
		return fmt.Errorf("修改任务为处理中(%d)失败 %v", taskID, err)
	}
	return nil
}

func FinishTask(db db.IDB, taskID int64) error {
	imap := map[string]interface{}{
		"task_id": taskID,
	}
	row, _, _, err := db.Execute(sqlFinishTask, imap)
	if err != nil || row != 1 {
		return fmt.Errorf("关闭任务(%d)失败 %v", taskID, err)
	}
	return nil
}
