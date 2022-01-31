package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ChoiBoyoon/learngo/scrapper"
	"github.com/labstack/echo"
)

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

const fileName string = "jobs.csv"

func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q="+term+"&limit=50"
	
	var jobs []extractedJob	//slice to be writed in file

	c := make(chan []extractedJob) //channel
	totalPages := getPageNum(baseURL)

	for i:=0;i<totalPages;i++ {
		go scrapeOnePage(i, baseURL, c) //receive slice of extractedJobs through channel
	}

	for i:=0;i<totalPages;i++ {
		extractedJobs := <-c //temporary variable to store received data through channel
		jobs = append(jobs, extractedJobs)
	}

	writeJobs(jobs)

	fmt.Println("Extracted ", len(jobs), " jobs.")
}

func getPageNum(page int) []extractedJob {
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

func scrapeOnePage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob //slice to be channeled to main()

	c := make(chan extractedJob) //channel to receive from extractJob()

	pageURL := url + "&start=" + strconv.Itoa(page*50) //URL of each page

	fmt.Println("requesting: ", pageURL)

	res,err := http.Get(pageURL)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close() //close after using. prevent memory leaks

	doc,err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	//crawling
	searchCards := doc.Find(".tapItem") //tapItem contains one job offer
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})
	for i:=0;i<searchCards.Length();i++ {
		job := <-c //temp var to receive a extractedJob through channel from extractJob()
		jobs = append(jobs, job)
	}

	mainC <- jobs //send slice to main() through mainC(main channel)
}

func extractJob(card *goquery.Selection) extractedJob {
	id,_ := card.Attr("data-jk")
	title := cleanString(card.Find("h2>span").Text())
	location := cleanString(card.Find(".company_location>div").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())
	c<-extractedJob{
		id:	id,
		title: title,
		location: location,
		salary: salary,
		summary: summary
	}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func writeJobs(jobs []extractedJobs) { //write slice of extractedJobs to .csv file
	file, err := os.Create(fileName)
	checkErr(err)
	
	//encoding for Korean
	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	file.Write(utf8bom)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, jobs := range jobs { //write in file line by line
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk="+job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func checkErr(err error) {
	if err!=nil{ 			//if there's an error
		log.Fatalln(err) 	//stop and exit function
	}
}

func checkStatusCode(res *http.Response) { //check response code after requesting html page
	if res.StatusCode != 200 { //200 means OK
		log.Fatalln("Request failed with status: ", res.StatusCode)
	}
}