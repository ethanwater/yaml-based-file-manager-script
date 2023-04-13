package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const ORIGIN string = "/Users/ethanwater/Downloads/"
const TEXTS string = "/Users/ethanwater/Documents/texts/"
const IMAGES string = "/Users/ethanwater/Documents/images/"

func IsImage(file string) bool {
	image := false
	if strings.HasSuffix(file, ".jpg") {
		image = true
	} else if strings.HasSuffix(file, ".jpeg") {
		image = true
	} else if strings.HasSuffix(file, ".png") {
		image = true
	} else if strings.HasSuffix(file, ".gif") {
		image = true
	}

	return image
}

func IsText(file string) bool {
	text := false
	if strings.HasSuffix(file, ".txt") {
		text = true
	} else if strings.HasSuffix(file, ".log") {
		text = true
	} else if strings.HasSuffix(file, ".pdf") {
		text = true
	} else if strings.HasSuffix(file, ".pages") {
		text = true
	}
	return text
}

func Type(file string) string {
	fileType := "null"
	if IsImage(file) {
		fileType = "img"
	} else if IsText(file) {
		fileType = "txt"
	} else {
		fileType = "misc"
	}

	return fileType
}

func UnsafeOrganize(path string) {
	if path == "" {
		currentDirectory, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path = currentDirectory
	}
	fmt.Printf("organizing: %s \n", path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.DefaultBytes(
		int64(len(files)),
		"organizing...",
	)

	for _, file := range files {
		oldPath := path + file.Name()

		switch Type(file.Name()) {
		case "img":
			newPath := IMAGES + file.Name()
			os.Rename(oldPath, newPath)
			bar.Add(1)
		case "txt":
			newPath := TEXTS + file.Name()
			os.Rename(oldPath, newPath)
			bar.Add(1)
		default:
			bar.Add(1)
		}
	}
}

func CheckLogStatus() bool { //can possibly be revoked
	exists := false
	if _, err := os.Stat("backup.log"); err == nil {
		exists = true
	} else if errors.Is(err, os.ErrNotExist) {
		exists = false
	}

	return exists
}

func SafeOrganize(path string) {
	Clear()
	//format path for files within for custom input directory
	//weird bug, jump in len(files) ?
	//if !strings.HasSuffix(path, "/") {
	//	path += "/"
	//}

	//check if log exists, if not create (kind of useless tho ngl)
	logStatus := CheckLogStatus()
	if logStatus == false {
		_, err := os.Create("backup.log")
		if err != nil {
			log.Fatal(err)
		}
	}

	//default path declaration
	if path == "" {
		currentDirectory, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path = currentDirectory + "/" //proper formatting of original path
	}
	fmt.Printf("organizing: %s \n", path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.DefaultBytes(
		int64(len(files)),
		"organizing...",
	)

	for _, file := range files {
		oldPath := path + file.Name()
		backup, err := os.OpenFile("backup.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := backup.Write([]byte(file.Name() + "\n")); err != nil {
			log.Fatal(err)
		}

		switch Type(file.Name()) {
		case "img":
			newPath := IMAGES + file.Name()
			os.Rename(oldPath, newPath)
			bar.Add(1)
		case "txt":
			newPath := TEXTS + file.Name()
			os.Rename(oldPath, newPath)
			bar.Add(1)
		default:
			bar.Add(1)
		}
	}
}

func Revert() {
	readLog, err := os.Open("backup.log")
	if err != nil {
		log.Fatal(err)
	}

	logScanner := bufio.NewScanner(readLog)
	logScanner.Split(bufio.ScanLines)
	var logLines []string

	for logScanner.Scan() {
		logLines = append(logLines, logScanner.Text())
	}
	readLog.Close()

	if len(logLines) == 0 {
		log.Fatal("nothing to revert")
	}
	fmt.Println("reverting changes")

	bar := progressbar.DefaultBytes(
		int64(len(logLines)),
		"reverting...",
	)

	for _, line := range logLines {
		originalPath := ORIGIN + line
		switch Type(line) {
		case "img":
			newPath := IMAGES + line
			os.Rename(newPath, originalPath)
			bar.Add(1)
		case "txt":
			newPath := TEXTS + line
			os.Rename(newPath, originalPath)
			bar.Add(1)
		default:
			bar.Add(1)
		}

	}

	Clear()
}

func Clear() {
	if err := os.Truncate("backup.log", 0); err != nil {
		log.Printf("failed to clear backup %s", err)
	}
}

func main() {}
