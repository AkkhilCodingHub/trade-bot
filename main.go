package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

type NewsItem struct {
	Title           string        `json:"title"`
	Url             string        `json:"url"`
	Summary         string        `json:"summary"`
	TickerSentiment []interface{} `json:"ticker_sentiment"`
}

type StockData struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func fetchNews(w http.ResponseWriter, r *http.Request) {
	apiKey := getEnv("ALPHAVANTAGE_API_KEY", "")
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=NEWS_SENTIMENT&apikey=%s", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	feed, _ := data["feed"].([]interface{})
	news := []NewsItem{}
	for i, item := range feed {
		if i >= 10 {
			break
		}
		itm := item.(map[string]interface{})
		news = append(news, NewsItem{
			Title:           fmt.Sprintf("%v", itm["title"]),
			Url:             fmt.Sprintf("%v", itm["url"]),
			Summary:         fmt.Sprintf("%v", itm["summary"]),
			TickerSentiment: itm["ticker_sentiment"].([]interface{}),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func fetchStock(w http.ResponseWriter, r *http.Request) {
	apiKey := getEnv("ALPHAVANTAGE_API_KEY", "")
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		symbol = "RELIANCE.NSE"
	}
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&outputsize=compact&apikey=%s", symbol, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	series, ok := data["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "No data found"})
		return
	}
	result := []StockData{}
	count := 0
	for date, values := range series {
		if count >= 100 {
			break
		}
		val := values.(map[string]interface{})
		result = append(result, StockData{
			Date:   date,
			Open:   parseFloat(val["1. open"]),
			High:   parseFloat(val["2. high"]),
			Low:    parseFloat(val["3. low"]),
			Close:  parseFloat(val["4. close"]),
			Volume: parseFloat(val["5. volume"]),
		})
		count++
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func parseFloat(v interface{}) float64 {
	if s, ok := v.(string); ok {
		var f float64
		fmt.Sscanf(s, "%f", &f)
		return f
	}
	return 0
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}

	// Serve static files from ./static
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// fallback to static for any file under /static/
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "./static/index.html")
	})

	http.HandleFunc("/api/news", fetchNews)
	http.HandleFunc("/api/stock", fetchStock)
	log.Println("Go backend running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
