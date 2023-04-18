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
	//establishes a config form the YAML
	var (
		name string
		path string
		ext  string
	)
	//name
	if n, err := config["name"].(string); err {
		name = n
	} else {
		fmt.Println(err)
	}
	//path
	if p, err := config["path"].(string); err {
		path = p
		if name == "origin" {
			ORIGIN = path
		}
	} else {
		fmt.Println(err)
	}
	//extensions
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
	//Organizes the files in an unsafe manner, void of logging
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
	fmt.Println("complete!")
	fmt.Println("time elapsed: ", end.Sub(start))

}

func SafeOrganize() {
	//Uses a logging system to log every file move executed. Creates/writes to a
	//backup file, which is necessary to call Revert
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
	fmt.Println("complete!")
	fmt.Println("time elapsed: ", end.Sub(start))
}

func Revert() {
	//Reverts the previous SafeOrganization actions, if the backup is empty,
	//this will do nothing
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
	fmt.Println("complete!")
	fmt.Println("time elapsed: ", end.Sub(start))
}

func Test() {
	//Test is used to make sure that the Safe and Unsafe Organization functions
	//run smoothly without any error. No files will be renamed/moved during this
	//process.
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
	fmt.Println("complete!")
	fmt.Println("time elapsed: ", end.Sub(start))
}

//func Scan() {
//	//Scans the given directory and returns information about the directory and
//	//the files within
//	Configurations()
//
//	path := ORIGIN
//	hiddenCount, directoryCount, matchConfig := 0, 0, 0
//
//	files, err := ioutil.ReadDir(path)
//	fileCount := len(files)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if len(files) == 0 {
//		log.Fatal("nothing to organize")
//	}
//
//	fmt.Printf("organizing: %s \n", path)
//	bar := progressbar.DefaultBytes(
//		int64(fileCount),
//		"scanning...",
//	)
//
//	for _, file := range files {
//		if strings.HasPrefix(file.Name(), ".") {
//			hiddenCount += 1
//			bar.Add(1)
//			continue
//		} else if file.IsDir() {
//			directoryCount += 1
//			bar.Add(1)
//			continue
//		}
//
//		oldPath := path + file.Name()
//		for _, config := range CONFIGS {
//			for _, extension := range config.Ext {
//				if strings.Contains(oldPath, extension) {
//					matchConfig += 1
//					break
//				}
//			}
//		}
//		bar.Add(1)
//	}
//
//	fmt.Println("\nRESULTS")
//	fmt.Printf("total files: %d\n", fileCount)
//	fmt.Printf("hidden files: %d\n", hiddenCount)
//	fmt.Printf("directory count: %d\n", directoryCount)
//	fmt.Printf("matching configurations: %d\n", matchConfig)
//	fmt.Printf("non-matching configurations: %d\n", fileCount-matchConfig)
//}

func DeepScan() {
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
	timeElapsed := end.Sub(start)

	fmt.Println("complete!")
	fmt.Println("time elapsed: ", timeElapsed)
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
	DeepScan()
}
