package main

import (//"fmt" 
		"net/http"
		"io/ioutil"
		"html/template"
		"encoding/xml"
		"sync"
)

var wg sync.WaitGroup

type Sitemapindex struct {
  Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	Keyword string
	Location string
}

type NewsAggPage struct {
    Title string
    News map[string]NewsMap
}

func newsRoutine(c chan News,Location string){
	defer wg.Done()
	var n News
	resp, _ := 	 http.Get(Location)
	bytes,_ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes,&n)
   	resp.Body.Close()	
   	c <- n
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := 	 http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes,_ := ioutil.ReadAll(resp.Body)
	var s Sitemapindex
	xml.Unmarshal(bytes,&s)	
	news_map := make(map[string]NewsMap)
    resp.Body.Close()
 	queue := make(chan News, 30)

	for _, Location := range s.Locations {
		wg.Add(1)
        go newsRoutine(queue,Location)		
	}

	wg.Wait()
    close(queue)

	for elem := range queue {
        for idx, _ := range elem.Keywords {
            news_map[elem.Titles[idx]] = NewsMap{elem.Keywords[idx], elem.Locations[idx]}
        }
    }

	p := NewsAggPage{Title: "Amazing News Aggregator", News: news_map}
    t, _ := template.ParseFiles("newsaggtemplate.html")
    t.Execute(w, p)	
}

func main() {
	http.HandleFunc("/", newsAggHandler)
	http.ListenAndServe(":80", nil) 
}