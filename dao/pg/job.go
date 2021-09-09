package pg

import (
	"auth/enviroment"
	"auth/model"
)

func SaveJob(fileName string, userId uint) *model.JobInfo {

	jobInfo := &model.JobInfo{CreatedBy: userId}

	enviroment.Env.DB.Save(jobInfo)

	return jobInfo
}
