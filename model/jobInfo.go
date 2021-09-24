package model

import (
	"auth/utils"
	"time"
)

type JobInfo struct {
	Id        *uint           `json:"id" gorm:"primaryKey"`
	Created   time.Time       `json:"created"`
	CreatedBy uint            `json:"createdBy"`
	FileName  string          `json:"fileName"`
	Status    utils.JobStatus `json:"status" gorm:"column:job_status_id"`
	Finished  *time.Time      `json:"finished,omitempty"`

	Statistics []*JobStatics `json:"statistics" gorm:"one2many:job_statistic"`
}

type JobStatics struct {
	Id        *uint  `json:"id" gorm:"primaryKey"`
	JobInfoId *uint  `json:"-" gorm:"column:job_info_id"`
	Term      string `json:"term"`
	Count     uint   `json:"count"`
}

func (jobInfo *JobInfo) TableName() string {
	return "job_info"
}

func (jobStatics *JobStatics) TableName() string {
	return "job_statistic"
}
