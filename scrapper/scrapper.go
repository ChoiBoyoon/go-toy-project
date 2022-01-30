package main

import (
	"fmt"
	"net/http"
	"strconv"
)

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

func getPage(page int) []extractedJob {
	var jobs []extractedJob
	
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("requesting: ", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() //prevent memory leaks

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErro(err)

	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection){
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {
	id,_ := card.Attr("data-jk")
	title := cleanString(card.Find("h2>span").Text())
	location := cleanString(card.Find(".company_location>div").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())
	c<-extractedJob{
		id:	id,
		title:title,
		location: location,
		salary:salary,
		summary:summary
	}
}