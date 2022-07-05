package baiducloud

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/util"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/mitchellh/go-homedir"
)

// DefaultTimeout timeout for common product, bcc e.g.
const DefaultTimeout = 180 * time.Second
const DefaultDebugMsg = "\n*************** %s Response *************** \n%+v\n%s******************************\n\n"

const (
	PaymentTimingPostpaid = "Postpaid"
	PaymentTimingPrepaid  = "Prepaid"
)

func debugOn() bool {
	for _, part := range strings.Split(os.Getenv("DEBUG"), ",") {
		if strings.TrimSpace(part) == "terraform" {
			return true
		}
	}
	return false
}

func addDebug(action, content interface{}) {
	if debugOn() {
		trace := "[DEBUG TRACE]:\n"
		for skip := 1; skip < 3; skip++ {
			_, filepath, line, _ := runtime.Caller(skip)
			trace += fmt.Sprintf("%s:%d\n", filepath, line)
		}

		//fmt.Printf(DefaultDebugMsg, action, content, trace)
		log.Printf(DefaultDebugMsg, action, content, trace)
	}
}

// write data to file
func writeToFile(filePath string, data interface{}) error {
	if strings.HasPrefix(filePath, "~") {
		usr, errCurrent := user.Current()
		if errCurrent != nil {
			return fmt.Errorf("get current user error: %s", errCurrent.Error())
		}
		if usr.HomeDir != "" {
			filePath = strings.Replace(filePath, "~", usr.HomeDir, 1)
		}
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stat file error: %s", err.Error())
	}

	if fileInfo != nil {
		if errRemove := os.Remove(filePath); errRemove != nil {
			return fmt.Errorf("delete old file error: %s", errRemove.Error())
		}
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json marshal error: %s", err.Error())
	}

	return ioutil.WriteFile(filePath, []byte(bytes), 0644)
}

// write data to file
func writeStringToFile(filePath string, data string) error {
	if strings.HasPrefix(filePath, "~") {
		usr, errCurrent := user.Current()
		if errCurrent != nil {
			return fmt.Errorf("get current user error: %s", errCurrent.Error())
		}
		if usr.HomeDir != "" {
			filePath = strings.Replace(filePath, "~", usr.HomeDir, 1)
		}
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stat file error: %s", err.Error())
	}

	if fileInfo != nil {
		if errRemove := os.Remove(filePath); errRemove != nil {
			return fmt.Errorf("delete old file error: %s", errRemove.Error())
		}
	}

	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

func buildClientToken() string {
	return util.NewUUID()
}

func buildStateConf(pending, target []string, timeout time.Duration, f resource.StateRefreshFunc) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Delay:      10 * time.Second,
		Pending:    pending,
		Refresh:    f,
		Target:     target,
		Timeout:    timeout,
		MinTimeout: 3 * time.Second,
	}
}

func stringInSlice(strs []string, value string) bool {
	for _, str := range strs {
		if value == str {
			return true
		}
	}

	return false
}

// check two strings are equal or not
// if both strings are one of defaultStr value, return true
func stringEqualWithDefault(s1, s2 string, defaultStr []string) bool {
	isDefaultS1 := stringInSlice(defaultStr, s1)
	isDefaultS2 := stringInSlice(defaultStr, s2)

	if isDefaultS1 != isDefaultS2 {
		return false
	}

	if s1 != s2 {
		if !isDefaultS1 {
			return false
		}
	}

	return true
}

func loadFileContent(v string) ([]byte, error) {
	filename, err := homedir.Expand(v)
	if err != nil {
		return nil, err
	}
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func zipFileDir(path string) ([]byte, error) {
	fileDir, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}

	zipFileBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipFileBuffer)

	err = filepath.Walk(fileDir, func(path string, info os.FileInfo, errs error) error {
		if info.IsDir() {
			return nil
		}

		zipFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = zipFile.Close()
		}()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.Replace(path, fileDir, "./", -1)
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if _, err := io.Copy(writer, zipFile); err != nil {
			return err
		}
		return zipWriter.Flush()
	})
	if err != nil {
		_ = zipWriter.Close()
		return nil, err
	}

	// Close() will write some final data to buffer, so zipWriter should be closed before read zip file from buffer
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}
	return zipFileBuffer.Bytes(), err
}
