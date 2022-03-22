# Goroutines-Concurency-Merge-Multiple-Files
A simple script to merge files concurently in a path with max_line count. When max_line count was reached another output file will be created. 

## Running Time Enviroments
arg1= input path, arg2= output path, arg3= output files base name, arg4= max_line count

## Build and Compile
```bash
$ go build
$ go run main.go arg1 arg2 arg3 arg4
```
Example: 
```bash
$ go build
$ go run main.go ./data ./out base_out 4
```
