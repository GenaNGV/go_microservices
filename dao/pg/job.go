package pg

import (
	"auth/enviroment"
	"auth/model"
	"auth/utils"
	"time"
)

func CreateJob(fileName string, userId uint) *model.JobInfo {

	jobInfo := &model.JobInfo{FileName: fileName, CreatedBy: userId, Status: utils.Assigned, Created: time.Now()}

	return jobInfo
}

func SaveJob(jobInfo *model.JobInfo) {
	enviroment.Env.DB.Save(jobInfo)
}

func SaveJobStatics(jobStatistics *model.JobStatics) {
	enviroment.Env.DB.Save(jobStatistics)
}
