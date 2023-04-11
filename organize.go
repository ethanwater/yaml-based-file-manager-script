package main

import (
    "log"
    "io/ioutil"
	"os"
	"strings"
)

const DOWNLOADS string = "/Users/xxx/Downloads/"
const TEXTS string = "/Users/xxx/Documents/texts/"
const IMAGES string = "/Users/xxx/Documents/images/"
const MISC string =  "/Users/xxx/Documents/misc/"


func IsImage(file string) bool {
	image := false
	if strings.HasSuffix(file, ".jpg") {
		image = true
	} else if strings.HasSuffix(file, ".jpeg") {
		image = true
	} else if strings.HasSuffix(file, ".png") {
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

func main() {
	files, err := ioutil.ReadDir(DOWNLOADS)
	if err != nil { 
		log.Fatal(err) 
	}

    for _, file := range files {
		old_path := DOWNLOADS + file.Name()

		switch Type(file.Name()) {
			case "img":
				new_path := IMAGES + file.Name()
				os.Rename(old_path, new_path)
			case "txt":
				new_path := TEXTS + file.Name()
				os.Rename(old_path, new_path)
			default:
				new_path := MISC + file.Name()
				os.Rename(old_path, new_path)
		}
    }
}
