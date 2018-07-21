package main

import ("fmt" 
		"net/http"
		"io/ioutil"
		"html/template"
		"encoding/xml")

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


func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := 	 http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes,_ := ioutil.ReadAll(resp.Body)
	var s Sitemapindex
	xml.Unmarshal(bytes,&s)	
	var n News
	news_map := make(map[string]NewsMap)

	for _, Location := range s.Locations {
		resp, _ := 	 http.Get(Location)
		bytes,_ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(bytes,&n)
		for idx, _ := range n.Keywords {
			news_map[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}
		fmt.Println(n)
	}
	p := NewsAggPage{Title: "Amazing News Aggregator", News: news_map}
    t, _ := template.ParseFiles("newsaggtemplate.html")
    t.Execute(w, p)	
}

func main() {
	http.HandleFunc("/", newsAggHandler)
	http.ListenAndServe(":80", nil) 
}