package file_utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"gopkg.in/yaml.v3"
)

func ReadFilesFromMultiPath(filePaths []string) ([]interface{}, error) {
	var err error
	var result []interface{}
	var tmpResult []interface{}
	for _, filePath := range filePaths {
		if tmpResult, err = ReadFiles(filePath); err != nil {
			return tmpResult, err
		}
		result = append(result, tmpResult...)
	}
	return result, err
}

func ReadFiles(filePath string) ([]interface{}, error) {
	var result []interface{}

	fileStat, err := os.Stat(filePath)
	if err != nil {
		return result, err
	}

	if fileStat.IsDir() {
		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			return result, err
		}
		for _, file := range files {
			path := filepath.Join(filePath, file.Name())
			var tmpResult []interface{}
			if tmpResult, err = ReadFiles(path); err != nil {
				return result, err
			}
			result = append(result, tmpResult...)
		}
		return result, nil
	}

	var tmpResult []byte
	if tmpResult, err = ioutil.ReadFile(filePath); err != nil {
		return result, err
	}
	splitedResult := bytes.Split(tmpResult, []byte("\n---"))
	for _, tmp := range splitedResult {
		data := make(map[string]interface{})
		if err = yaml.Unmarshal(tmp, &data); err != nil {
			return result, err
		}
		if len(data) == 0 {
			continue
		}
		result = append(result, data)
	}
	return result, err
}

func ReadFile(tctx *logger.TraceContext, filePath string, data interface{}) (err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(filePath); err != nil {
		return
	}
	err = yaml.Unmarshal(bytes, data)
	return
}
