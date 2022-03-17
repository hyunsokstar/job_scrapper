package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id              string
	title           string
	companyName     string
	companyLocation string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	// 10 페이지 분량의 리스트에 대한 url 출력해 보기
	for i := 0; i < totalPages; i++ {
		// getPage(i)
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}
	// fmt.Println(jobs)
	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	file.Write(utf8bom)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "CompanyName", "CompanyLocation"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		// jobSlice := []string{job.id, job.title, job.companyName, job.companyLocation}
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.companyName, job.companyLocation}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}

}

// func getPage(page int) {
// 	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // strconv.Itoa() <=> 문자열로 바꾸기
// 	fmt.Println("Requesting", pageURL)
// }

func getPage(page int) []extractedJob {
	var jobs []extractedJob
	// strconv.Itoa() <=> 문자열로 바꾸기
	// start 는 0 , 50 , 100 , 150, 200 이런식이어야 함
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)

	// fmt.Println("Requesting", pageURL)

	// 추가 11 start
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".sponTapItem ")

	searchCards.Each(func(i int, card *goquery.Selection) {
		// id, _ := card.Attr("data-jk")
		// // title := card.Find(".jobTitle > span").Text()
		// title := cleanString(card.Find(".jobTitle > span").Text())
		// company_name := cleanString(card.Find(".companyName").Text())
		// companyLocation := cleanString(card.Find(".companyLocation").Text())
		// fmt.Println(id, title, company_name, companyLocation)
		job := extractJob(card)
		jobs = append(jobs, job)

	})
	return jobs

}

func extractJob(card *goquery.Selection) extractedJob {

	id, _ := card.Attr("data-jk")
	// title := card.Find(".jobTitle > span").Text()
	title := cleanString(card.Find(".jobTitle > span").Text())
	companyName := cleanString(card.Find(".companyName").Text())
	companyLocation := cleanString(card.Find(".companyLocation").Text())
	// fmt.Println(id, title, companyName, companyLocation)

	return extractedJob{
		id:              id,
		title:           title,
		companyName:     companyName,
		companyLocation: companyLocation,
	}

}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// 특정 페이지 정보를 받아오는 함수
func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)

	checkErr(err)  // error check
	checkCode(res) // status code check

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// fmt.Println(s.Find("a").Length())
		pages = s.Find("a").Length()

	})
	// return 0
	return pages
}

// error check
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// status code check
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status", res.StatusCode)
	}

}
