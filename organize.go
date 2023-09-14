package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

var ORIGIN string
var TEXTS string
var IMAGES string
var CONFIGS []Config
var CONFIG_FILE string

type Config struct {
	Name string
	Path string
	Ext  []string
}

func OpenConfig() {
	cmdOpenConfig := exec.Command("open", CONFIG_FILE)
	_, err := cmdOpenConfig.Output()
	if err != nil {
		fmt.Println(err)
	}
}

func SetConfig(config string) {
	CONFIG_FILE = config
}

func ClearBackup() {
	if err := os.Truncate("backup.log", 0); err != nil {
		log.Printf("failed to clear backup %s", err)
	}
}

func Configurations() {
	configFile, err := ioutil.ReadFile("example_config.yaml")
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
//Configuration initialization. This is the first function performed in order to
//allow organize.go to work properly. Due to the YAML integration, organize.go
//allows for a highly flexible and customizable configuration system to best fit
//the users needs.
	var (
		name string
		path string
		ext  string
	)

	if n, err := config["name"].(string); err {
		name = n
	} else {
		fmt.Println(err)
	}

	if p, err := config["path"].(string); err {
		path = p
		if name == "ORIGIN" {
			ORIGIN = path
		}
	} else {
		fmt.Println(err)
	}

	if e, err := config["ext"].(string); err {
		ext = e
	} else {
		fmt.Println(err)
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

	start := time.Now()
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

	end := time.Now()
	fmt.Printf("complete!\ntime elapsed: %f secs", end.Sub(start).Seconds())
}


func SafeOrganize() {
//SafeOrganize organizes the origin directory using the same file organization
//method used in the UnsafeOrganize function. However, it updates the backup.log
//file the program created in order to keep track of the previous paths of the
//now moved files. This way, in case of any error in the file organization, the
//user can call the Revert function to revert any files moved to their original
//locations.
	ClearBackup()
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
	start := time.Now()
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
	end := time.Now()
	fmt.Printf("complete!\ntime elapsed: %f secs", end.Sub(start).Seconds())
}

func Revert() {
//Reverts aby changes from the last SafeOrganize call. Can only revert the past
//organzation call. Any calls before that must be changed manually by the user.
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

	start := time.Now()
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
	end := time.Now()
	fmt.Printf("complete!\ntime elapsed: %f secs", end.Sub(start).Seconds())
}

func Test() {
//The Test function tests a "pseudo-organization" in order to assure that either
//Unsafe or Safe organization methods will work properly. In almost every case,
//if there is an error, the organzaiton methods will notice it and notify the
//user of the error.
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
		"testing organize...",
	)

	start := time.Now()
	for _, file := range files {
		oldPath := path + file.Name()
		for _, config := range CONFIGS {
			for _, extension := range config.Ext {
				if strings.Contains(oldPath, extension) {
					break
				}
			}
		}
		bar.Add(1)
	}
	end := time.Now()
	fmt.Printf("complete!\ntime elapsed: %f secs", end.Sub(start).Seconds())
}

func Scan() {
 	//Scans the given directory and returns information about the directory and
 	//the files within
 	Configurations()
 
 	path := ORIGIN
 	hiddenCount, directoryCount, matchConfig := 0, 0, 0
 
 	files, err := ioutil.ReadDir(path)
 	fileCount := len(files)
	if err != nil {
 		log.Fatal(err)
 	}
 	if len(files) == 0 {
 		log.Fatal("nothing to organize")
 	}
 
 	fmt.Printf("organizing: %s \n", path)
 	bar := progressbar.DefaultBytes(
 		int64(fileCount),
 		"scanning...",
	)
	start := time.Now()
 
	for _, file := range files {
 		if strings.HasPrefix(file.Name(), ".") {
 			hiddenCount += 1
 			bar.Add(1)
 			continue
 		} else if file.IsDir() {
 			directoryCount += 1
 			bar.Add(1)
 			continue
 		}
		oldPath := path + file.Name()
 		for _, config := range CONFIGS {
 			for _, extension := range config.Ext {
 				if strings.Contains(oldPath, extension) {
 					matchConfig += 1
 					break
 				}
 			}
			bar.Add(1)
		}
	}
	end := time.Now()
	fmt.Printf("complete!\ntime elapsed: %f secs", end.Sub(start).Seconds())
}

func DeepScan() {
//DeepScan returns statistics about the directories that are initialized in the
//configuration file.
	Configurations()

	path := ORIGIN
	hiddenCount, directoryCount, matchConfig := 0, 0, 0

	directory, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	directoryStat, err := directory.Stat()
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	fileCount := len(files)
	if fileCount == 0 {
		log.Fatal("nothing to organize")
	}

	fmt.Printf("deep scanning: %s \n", path)
	bar := progressbar.DefaultBytes(
		int64(fileCount),
		"scanning...",
	)

	extMap := make(map[string]int)
	start := time.Now()
	for _, config := range CONFIGS {
		if config.Name == "origin" {
			continue
		}
		count := 0
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				hiddenCount += 1
				bar.Add(1)
				continue
			} else if file.IsDir() {
				directoryCount += 1
				bar.Add(1)
				continue
			}
			for _, extension := range config.Ext {
				if strings.Contains(file.Name(), extension) {
					matchConfig += 1
					count += 1
					bar.Add(1)
					break
				}
			}
		}
		bar.Add(1)
		extMap[config.Name] = count
	}

	end := time.Now()

	fmt.Printf("complete!\ntime elapsed: %f secs", end.Sub(start).Seconds())
	fmt.Printf("STATISTICS: %+v\n", directoryStat)
	fmt.Printf("TOTAL: %d\n", fileCount)
	for key, value := range extMap {
		fmt.Printf("%s: %d\n", key, value)
	}
	fmt.Printf("HIDDEN: %d\n", hiddenCount)
	fmt.Printf("DIRECTORIES: %d\n", directoryCount)
	fmt.Printf("CONFIGURATIONS: %d\n", matchConfig)
	fmt.Println("FILES: ")
	for _, file := range files {
		fmt.Print(file.Name(), ", ")
	}
}

func main() {
	Scan()
}
