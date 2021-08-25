package report

import (
// "fmt"
)

type TaskReport struct {
	TaskId  string `json:"taskId" form:"taskId"`
	Key     string `json:"key" form:"key"`
	Value   string `json:"value" form:"value"`
	Detail  string `json:"detail" form:"detail"`
	Millsec int64  `json:"millsec" form:"millsec"`
}

// func (callback func(TaskReport) error) gin.HandlerFunc
