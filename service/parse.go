package service

import (
	"auth/dao/pg"
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

	if fileName == "" {
		return nil, ErrMissedFile
	}

	fileInfo, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, ErrMissedFile
	}
	if fileInfo.IsDir() {
		return nil, ErrMissedFile
	}

	jobInfo := pg.CreateJob(fileName, user.Id)

	go runJob(jobInfo, terms)

	return jobInfo, nil
}

func runJob(jobInfo *model.JobInfo, terms []string) {

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
	finished := time.Now()
	jobInfo.Finished = &finished
	jobInfo.Status = utils.Job_status_finished

	pg.SaveJob(jobInfo)
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
