# organize.go

This Go program is a file organization utility that allows you to organize files in a directory based on user-defined configurations specified in a YAML file. It provides options for safely organizing files, reverting changes, and performing deep scans to gather statistics about the files in the directory.

## Usage
Create a YAML configuration file named config.yaml with your desired file organization rules. Example:
```yaml
- name: origin
  path: /path/to/origin/
  ext: txt docx

- name: documents
  path: /path/to/documents/
  ext: pdf doc docx

- name: images
  path: /path/to/images/
  ext: jpg jpeg png
```
Open the configuration file using the open command:
```sh
go run main.go -open /path/to/config.yaml
```
Organizing Files
To organize files in the specified origin directory based on the configurations:
```sh
go run main.go -organize
```
Reverting Changes
To revert the changes made during the organization process:
```sh
go run main.go -revert
```
Testing
You can test the organization process without actually moving files:
```sh
go run main.go -test
```

Deep Scan and Statistics
To perform a deep scan of the files in the origin directory and gather statistics:
```sh
go run main.go -deepscan
```
