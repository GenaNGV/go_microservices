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

func Parse(fileName string, terms []string, user *model.UserAuth, rule string) (*model.JobInfo, error) {

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

	jobInfo.Status = utils.Processing
	enviroment.Env.DB.Save(jobInfo)

	if "concurrent" == rule {
		go runConcurrentJob(jobInfo, terms)
	} else if "ext_concurrent" == rule {
		go runExtConcurrentJob(jobInfo, terms)
	} else {
		go runJob(jobInfo, terms)
	}

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
	finishJob(jobInfo, utils.Finished)
}

func runConcurrentJob(jobInfo *model.JobInfo, terms []string) {

	log.WithFields(log.Fields{"file": jobInfo.FileName, "job": *jobInfo.Id}).Info("Parsing concurrent")

	f, err := os.Open(jobInfo.FileName)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)

	channel := make(chan int)

	total := 0

	var wg sync.WaitGroup

	for scanner.Scan() {
		line := scanner.Text()

		for i := 0; i < len(terms); i++ {
			wg.Add(1)

			term := terms[i]

			go func(c chan<- int) {
				defer wg.Done()
				log.WithFields(log.Fields{"job": *jobInfo.Id, "term": term, "line": line}).Info("locate...")

				findConcurrentTerm(term, c, line)
			}(channel)
		}
	}

	go func(c chan int) {
		wg.Wait()
		close(c)

		log.WithFields(log.Fields{"job": *jobInfo.Id, "total": total}).Info("calculate...")
	}(channel)

	for v := range channel {
		log.WithFields(log.Fields{"job": *jobInfo.Id, "value": v}).Info("summary...")
		total += v
	}

	log.WithFields(log.Fields{"job": *jobInfo.Id, "total": total}).Info("Saving to database")
	finishJob(jobInfo, utils.Finished)
}

func runExtConcurrentJob(jobInfo *model.JobInfo, terms []string) {

	log.WithFields(log.Fields{"file": jobInfo.FileName, "job": *jobInfo.Id}).Info("Parsing concurrent")

	f, err := os.Open(jobInfo.FileName)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)

	channel := make(chan *model.JobStatics)

	total := 0

	var wg sync.WaitGroup

	for scanner.Scan() {
		line := scanner.Text()

		for i := 0; i < len(terms); i++ {
			wg.Add(1)

			term := terms[i]

			go func(c chan<- *model.JobStatics) {
				defer wg.Done()
				log.WithFields(log.Fields{"job": *jobInfo.Id, "term": term, "line": line}).Info("locate...")

				findExtConcurrentTerm(term, c, line)
			}(channel)
		}
	}

	go func(c chan *model.JobStatics) {
		wg.Wait()
		close(c)
	}(channel)

	m := make(map[string]*model.JobStatics)

	for v := range channel {
		log.WithFields(log.Fields{"job": *jobInfo.Id, "term": v.Term, "value": v.Count}).Info("summary...")

		jobStatics := m[v.Term]

		if jobStatics == nil {
			v.JobInfoId = jobInfo.Id
			m[v.Term] = v
		} else {
			jobStatics.Count += v.Count
		}
	}

	for _, value := range m {
		pg.SaveJobStatics(value)
	}

	log.WithFields(log.Fields{"job": *jobInfo.Id, "total": total}).Info("Saving to database")
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

func findConcurrentTerm(term string, channel chan<- int, line string) {

	total := 0

	index := strings.Index(line, term)

	for index >= 0 {
		total++

		line = line[index+1:]
		index = strings.Index(line, term)
	}

	channel <- total
}

func findExtConcurrentTerm(term string, channel chan<- *model.JobStatics, line string) {

	jobInfo := &model.JobStatics{Count: 0, Term: term}

	index := strings.Index(line, term)

	for index >= 0 {
		jobInfo.Count++

		line = line[index+1:]
		index = strings.Index(line, term)
	}

	channel <- jobInfo
}

func finishJob(jobInfo *model.JobInfo, status utils.JobStatus) {

	finished := time.Now()
	jobInfo.Finished = &finished
	jobInfo.Status = status

	enviroment.Env.DB.Save(jobInfo)
}
