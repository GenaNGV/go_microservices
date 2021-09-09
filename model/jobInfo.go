package model

import (
	"time"
)

type JobInfo struct {
	Id        uint       `json:"id" gorm:"primaryKey"`
	Created   time.Time  `json:"created"`
	CreatedBy uint       `json:"createdBy"`
	Status    uint       `json:"status"`
	Finished  *time.Time `json:"finished,omitempty"`

	Statistics []*JobStatics `json:"statistics" gorm:"one2many:job_statistic"`
}

type JobStatics struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	JobInfoId uint   `json:"-"`
	Term      string `json:"term"`
	Count     uint   `json:"count"`
}

func (jobInfo *JobInfo) TableName() string {
	return "job_info"
}

func (jobStatics *JobStatics) TableName() string {
	return "job_statistic"
}
