package service

import (
	"auth/dao/pg"
	"auth/enviroment"
	"auth/model"
	"auth/utils"
	"bufio"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ErrMissedFile = errors.New("missed file")
)

func Parse(fileName string, terms []string, user *model.UserAuth) (*model.JobInfo, error) {

	jobInfo := pg.CreateJob(fileName, user.Id)

	if fileName == "" {
		finishJob(jobInfo, utils.Error)

		return nil, ErrMissedFile
	}

	fileInfo, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		finishJob(jobInfo, utils.Error)
		return nil, ErrMissedFile
	}

	if fileInfo.IsDir() {
		finishJob(jobInfo, utils.Error)
		return nil, ErrMissedFile
	}

	go runJob(jobInfo, terms)

	return jobInfo, nil
}

func runJob(jobInfo *model.JobInfo, terms []string) {

	jobInfo.Status = utils.Processing
	enviroment.Env.DB.Save(jobInfo)

	log.WithFields(log.Fields{"file": jobInfo.FileName, "job": *jobInfo.Id}).Info("Parsing")

	var wg sync.WaitGroup

	channel := make(chan int, 1)
	channel <- 0

	for i := 0; i < len(terms); i++ {
		wg.Add(1)

		term := terms[i]

		go func() {
			defer wg.Done()
			findTerm(jobInfo, term, channel)
		}()
	}

	wg.Wait()

	log.WithFields(log.Fields{"job": *jobInfo.Id}).Info("Saving to database")
	finishJob(jobInfo, utils.Finished)
}

func findTerm(jobInfo *model.JobInfo, term string, channel chan int) {

	log.WithFields(log.Fields{"job": *jobInfo.Id, "term": term}).Info("Starting")

	f, err := os.Open(jobInfo.FileName)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	jobStatics := &model.JobStatics{JobInfoId: jobInfo.Id, Count: 0, Term: term}

	scanner := bufio.NewScanner(f)

	var old int
	for scanner.Scan() {
		line := scanner.Text()

		index := strings.Index(line, term)

		for index >= 0 {
			jobStatics.Count = jobStatics.Count + 1

			old = <-channel
			channel <- 1 + old

			log.WithFields(log.Fields{"job": *jobInfo.Id, "term": term, "found": jobStatics.Count, "total": 1 + old}).Info("Processing")

			line = line[index+1:]
			index = strings.Index(line, term)
		}
	}

	log.WithFields(log.Fields{"job": *jobInfo.Id, "term": term, "found": jobStatics.Count}).Info("Processed")

	pg.SaveJobStatics(jobStatics)
}

func finishJob(jobInfo *model.JobInfo, status utils.JobStatus) {

	finished := time.Now()
	jobInfo.Finished = &finished
	jobInfo.Status = status

	enviroment.Env.DB.Save(jobInfo)
}
