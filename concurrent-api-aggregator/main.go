package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type API struct {
	Name string
	URL  string
}
type Result struct {
	Name  string
	Data  string
	Error string
}

func main() {
	http.HandleFunc("/aggregate", aggregateHandler)
	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
func aggregateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	apis := []API{
		{
			Name: "Example1",
			URL:  "https://example.com",
		},
		{
			Name: "Example2",
			URL:  "https://www.google.com",
		},
		{
			Name: "Example3",
			URL:  "https://httpbin.org/get",
		},
		{
			Name: "SlowAPI",
			URL:  "https://httpbin.org/delay/5",
		},
	}
	var wg sync.WaitGroup
	result := make(chan Result)
	sem := make(chan struct{}, 2)
	for i := 0; i < len(apis); i++ {
		wg.Add(1)
		go Worker(apis[i], &wg, result, ctx, sem)
	}
	go func() {
		wg.Wait()
		close(result)
	}()

	var results []Result
	for res := range result {
		results = append(results, res)
		log.Printf("received result from %s", res.Name)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)

}

func Worker(api API, wg *sync.WaitGroup, result chan<- Result, ctx context.Context, sem chan struct{}) {

	defer wg.Done()

	select {
	case sem <- struct{}{}:
	case <-ctx.Done():
		result <- Result{
			Name:  api.Name,
			Data:  "0",
			Error: ctx.Err().Error(),
		}
		return
	}
	defer func() {
		<-sem
	}()
	log.Printf("starting %s", api.Name)
	apiCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, api.URL, nil)
	if err != nil {
		result <- Result{
			Name:  api.Name,
			Data:  "0",
			Error: err.Error(),
		}
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		result <- Result{
			Name:  api.Name,
			Data:  "0",
			Error: err.Error(),
		}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		result <- Result{
			Name:  api.Name,
			Data:  "0",
			Error: fmt.Sprintf("status code: %d", resp.StatusCode),
		}
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result <- Result{
			Name:  api.Name,
			Data:  "0",
			Error: err.Error(),
		}
		return
	}
	log.Printf("%s downloaded, bytes: %d", api.Name, len(body))
	result <- Result{
		Name: api.Name,
		Data: strconv.Itoa(len(body)),
	}

}
