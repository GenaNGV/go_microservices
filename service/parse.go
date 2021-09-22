package service

import (
	"auth/dao/pg"
	"auth/model"
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

func Parse(fileName string, arr []string, user *model.UserAuth) (*model.JobInfo, error) {

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

	go runJob(jobInfo, arr)

	return jobInfo, nil
}

func runJob(jobInfo *model.JobInfo, arr []string) {

	var wg sync.WaitGroup

	defer wg.Done()

	channel := make(chan int, 1)
	channel <- 0

	for i := 0; i < len(arr); i++ {
		wg.Add(1)
		go findTerm(jobInfo, arr[i], channel, i)
	}
	wg.Wait()

	log.WithFields(log.Fields{"file": jobInfo.FileName, "total": <-channel}).Info("Finished Parsing file")

	finished := time.Now()
	jobInfo.Finished = &finished
	jobInfo.Status = 3

	pg.SaveJob(jobInfo)

	// out results

}

func findTerm(jobInfo *model.JobInfo, term string, channel chan int, index int) {

	log.WithFields(log.Fields{"file": jobInfo.FileName, "term": term}).Info("Parsing file")

	f, err := os.Open(jobInfo.FileName)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()

	jobStatics := &model.JobStatics{JobInfoId: jobInfo.Id, Count: 0, Term: term}

	scanner := bufio.NewScanner(f)
	row := 0

	var old int
	for scanner.Scan() {
		line := scanner.Text()

		index := strings.Index(line, term)
		row++

		for index >= 0 {
			jobStatics.Count = jobStatics.Count + 1

			old = <-channel
			channel <- 1 + old

			log.WithFields(log.Fields{"file": jobInfo.FileName, "term": term, "row": row, "found": jobStatics.Count, "total": 1 + old}).Info("Processing file")

			line = line[index+1:]
			index = strings.Index(line, term)
		}
	}

	log.WithFields(log.Fields{"file": jobInfo.FileName, "term": term, "found": jobStatics.Count}).Info("Parsed file")

	pg.SaveJobStatics(jobStatics)
}
