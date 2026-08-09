package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type UniversitiesResponse struct {
	Page       int          `json:"page"`
	PerPage    int          `json:"per_page"`
	Total      int          `json:"total"`
	TotalPages int          `json:"total_pages"`
	Data       []University `json:"data"`
}

type University struct {
	University            string   `json:"university"`
	RankDisplay           string   `json:"rank_display"`
	Score                 float64  `json:"score"`
	Type                  string   `json:"type"`
	StudentFacultyRatio   int      `json:"student_faculty_ratio"`
	InternationalStudents string   `json:"international_students"`
	FacultyCount          string   `json:"faculty_count"`
	Location              Location `json:"location"`
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Region  string `json:"region"`
}

func main() {
	// ans := ""
	url := fmt.Sprintf("https://jsonmock.hackerrank.com/api/universities?page=%d", 1)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var data UniversitiesResponse
	country := "India"
	rank := math.MaxInt64
	school := ""
	json.Unmarshal(body, &data)
	for i := 1; i <= data.TotalPages; i++ {
		url := fmt.Sprintf("https://jsonmock.hackerrank.com/api/universities?page=%d", i)
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		var data UniversitiesResponse
		json.Unmarshal(body, &data)
		for _, ddd := range data.Data {
			if strings.EqualFold(ddd.Location.Country, country) {
				rankint, _ := strconv.Atoi(ddd.RankDisplay)
				// log.Println(ddd)
				if rankint < rank {

					rank = rankint
					school = ddd.University

				}
			}

		}

	}
	log.Println(school)
	log.Println("hiiii")

}
