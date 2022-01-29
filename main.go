package main

import "fmt"

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	for i:=0;i<totalPages;i++ {
		extractedJob := getPage(i)
		jobs = append(jobs, extractedJob...)
	}

	fmt.Println(jobs)
}