package service

import (
	"auth/model"
	"github.com/pkg/errors"
	"os"
)

var (
	ErrMissedFile = errors.New("missed file")
)

func Parse(fileName string, chars uint) (*model.JobInfo, error) {

	if fileName == "" {
		return nil, ErrMissedFile
	}

	fileInfo, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, ErrMissedFile
	}
	if !fileInfo.IsDir() {
		return nil, ErrMissedFile
	}


	return nil, nil
}


