package pg

import (
	"auth/enviroment"
	"auth/model"
	"auth/utils"
	"time"
)

func CreateJob(fileName string, userId uint) *model.JobInfo {

	jobInfo := &model.JobInfo{FileName: fileName, CreatedBy: userId, Status: utils.Job_status_assigned, Created: time.Now()}

	enviroment.Env.DB.Create(jobInfo)

	return jobInfo
}

func SaveJob(jobInfo *model.JobInfo) {
	enviroment.Env.DB.Save(jobInfo)
}

func SaveJobStatics(jobStatistics *model.JobStatics) {
	enviroment.Env.DB.Save(jobStatistics)
}
