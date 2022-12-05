package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ScrapeEnv struct {
	Endpoints []ScrapeEndpoint `json:"endpoints"`
}

type ScrapeEndpoint struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Port int    `json:"port"`
}

type ScrapeHost struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func main() {
	var scrapeEndpoints ScrapeEnv
	err := json.Unmarshal([]byte(os.Getenv("SCRAPE_ENDPOINTS")), &scrapeEndpoints)
	if err != nil {
		panic(err)
	}

	var scrapeConfig []ScrapeHost
	for _, endpoint := range scrapeEndpoints.Endpoints {
		host := ScrapeHost{
			Targets: []string{fmt.Sprintf("host.docker.internal:%v", endpoint.Port)},
			Labels: map[string]string{
				"job": endpoint.Name,
			},
		}
		scrapeConfig = append(scrapeConfig, host)
	}

	scrapeConfigJson, err := json.MarshalIndent(scrapeConfig, "", "  ")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Config queried")
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(scrapeConfigJson))
	})
	http.ListenAndServe(":8080", nil)
}
