package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const TEXTS string = "/Users/xxx/Documents/texts/"
const IMAGES string = "/Users/xxx/Documents/images/"
const MISC string = "misc/"

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
	filetype := "null"
	if IsImage(file) {
		filetype = "img"
	} else if IsText(file) {
		filetype = "txt"
	} else {
		filetype = "misc"
	}

	return filetype
}

func UnsafeOrganize(path string) {
	if path == "" {
		current_directory, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path = current_directory
		fmt.Printf("organizing: %s \n", path)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.Default(int64(len(files)))

	for _, file := range files {
		old_path := path + file.Name()

		switch Type(file.Name()) {
		case "img":
			new_path := IMAGES + file.Name()
			os.Rename(old_path, new_path)
			bar.Add(1)
		case "txt":
			new_path := TEXTS + file.Name()
			os.Rename(old_path, new_path)
			bar.Add(1)
		default:
			new_path := MISC + file.Name()
			os.Rename(old_path, new_path)
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
	//format path for files within for custom input directory
	//weird bug, jump in len(files) ?
	//if !strings.HasSuffix(path, "/") {
	//	path += "/"
	//}

	//check if log exists, if not create (kind of useless tho ngl)
	log_status := CheckLogStatus()
	if log_status == false {
		_, err := os.Create("backup.log")
		if err != nil {
			log.Fatal(err)
		}
	}

	//default path declaration
	if path == "" {
		current_directory, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path = current_directory + "/" //proper formatting of original path
		fmt.Printf("organizing: %s \n", path)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.Default(int64(len(files)))

	for _, file := range files {
		old_path := path + file.Name()
		backup, err := os.OpenFile("backup.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := backup.Write([]byte(old_path + "\n")); err != nil {
			log.Fatal(err)
		}

		switch Type(file.Name()) {
		case "img":
			new_path := IMAGES + file.Name()
			os.Rename(old_path, new_path)
			bar.Add(1)
		case "txt":
			new_path := TEXTS + file.Name()
			os.Rename(old_path, new_path)
			bar.Add(1)
		default:
			new_path := MISC + file.Name()
			os.Rename(old_path, new_path)
			bar.Add(1)
		}
	}
}

func Revert() {}

func main() {
	SafeOrganize("/Users/xxx/Downloads/")
}
