package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const TEXTS string = "/Users/ethanwater/Documents/texts/"
const IMAGES string = "/Users/ethanwater/Documents/images/"
const MISC string = "/Users/ethanwater/Documents/misc/"

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

func Organize(path string) {
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

func main() {
	Organize("")
}
