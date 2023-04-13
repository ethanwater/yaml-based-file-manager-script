package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

var ORIGIN string
var TEXTS string
var IMAGES string
var MISC string

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

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		fmt.Println("nothing to organize")
	} else {
		fmt.Printf("organizing: %s \n", path)
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
}

func CheckLogStatus() bool { //for testing purposes
	exists := false
	if _, err := os.Stat("backup.log"); err == nil {
		exists = true
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
	if CheckLogStatus() == false {
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

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		fmt.Println("nothing to organize")
	} else {
		fmt.Printf("organizing: %s \n", path)
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
				newPath := MISC + file.Name()
				os.Rename(oldPath, newPath)
				bar.Add(1)
			}
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
		fmt.Println("nothing to revert")
	} else {
		fmt.Println("reverting changes")

		bar := progressbar.DefaultBytes(
			int64(len(logLines)),
			"reverting...",
		)

		for _, file := range logLines {
			originalPath := ORIGIN + file
			switch Type(file) {
			case "img":
				newPath := IMAGES + file
				os.Rename(newPath, originalPath)
				bar.Add(1)
			case "txt":
				newPath := TEXTS + file
				os.Rename(newPath, originalPath)
				bar.Add(1)
			default:
				newPath := MISC + file
				os.Rename(newPath, originalPath)
				bar.Add(1)
			}
		}

		Clear()

	}
}

func Clear() {
	if err := os.Truncate("backup.log", 0); err != nil {
		log.Printf("failed to clear backup %s", err)
	}
}

func Config() {
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[interface{}]interface{})
	err2 := yaml.Unmarshal(configFile, &data)

	if err2 != nil {

		log.Fatal(err2)
	}

	for k, v := range data {
		if k == "origin" {
			ORIGIN = fmt.Sprint(v)
			continue
		} else if k == "images" {
			IMAGES = fmt.Sprint(v)
			continue
		} else if k == "texts" {
			TEXTS = fmt.Sprint(v)
			continue
		} else if k == "misc" {
			MISC = fmt.Sprint(v)
			continue
		}

	}
}

func main() {
	Config()
	SafeOrganize(ORIGIN)
}
