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
var CONFIGS []Config

type Config struct {
	Name string
	Path string
	Ext  []string
}

func Clear() {
	if err := os.Truncate("backup.log", 0); err != nil {
		log.Printf("failed to clear backup %s", err)
	}
}

func Configurations() {
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]interface{})
	err2 := yaml.Unmarshal(configFile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}
	for _, v := range data {
		CONFIGS = append(
			CONFIGS,
			CreateConfig(v.(map[string]interface{})),
		)
		continue
	}
}

func CreateConfig(config map[string]interface{}) Config {
	var (
		name string
		path string
		ext  string
	)
	//name
	if n, ok := config["name"].(string); ok {
		name = n
	} else {
		fmt.Println(ok)
	}
	//path
	if p, ok := config["path"].(string); ok {
		path = p
		if name == "origin" {
			ORIGIN = path
		}
	} else {
		fmt.Println(ok)
	}
	//extensions
	if e, ok := config["ext"].(string); ok {
		ext = e
	} else {
		fmt.Println(ok)
	}

	extensions := strings.Split(ext, " ")
	freshConfig := Config{
		Name: name,
		Path: path,
		Ext:  extensions,
	}
	return freshConfig
}

func UnsafeOrganize() {
	Configurations()
	path := ORIGIN

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		log.Fatal("nothing to organize")
	}

	fmt.Printf("organizing: %s \n", path)
	bar := progressbar.DefaultBytes(
		int64(len(files)),
		"organizing...",
	)

	for _, file := range files {
		oldPath := path + file.Name()
		for _, config := range CONFIGS {
			for _, extension := range config.Ext {
				if strings.Contains(oldPath, extension) {
					newPath := config.Path + file.Name()
					os.Rename(oldPath, newPath)
					break
				}
			}
		}
		bar.Add(1)
	}
}

func SafeOrganize() {
	Clear()
	Configurations()
	path := ORIGIN

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		log.Fatal("nothing to organize")
	}
	fmt.Printf("safely organizing: %s \n", path)

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

		for _, config := range CONFIGS {
			for _, extension := range config.Ext {
				if strings.Contains(oldPath, extension) {
					newPath := config.Path + file.Name()
					os.Rename(oldPath, newPath)
					break
				}
			}
		}
		bar.Add(1)
	}
}

func Revert() {
	Configurations()

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

	for _, file := range logLines {
		oldPath := ORIGIN + file
		for _, config := range CONFIGS {
			for _, extension := range config.Ext {
				if strings.Contains(oldPath, extension) {
					newPath := config.Path + file
					os.Rename(newPath, oldPath)
					break
				}
			}
		}
		bar.Add(1)
	}
}

func main() {
	SafeOrganize()
}
