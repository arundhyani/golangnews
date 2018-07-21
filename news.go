package main

import ("fmt" 
		"net/http"
		"io/ioutil"
		"encoding/xml")

type Sitemapindex struct {
  Locations []Location `xml:"sitemap"`
}

type Location struct {
  Loc string `xml:"loc"`
}

func (l Location) String () string {
  return fmt.Sprintf(l.Loc)
}

func main() {
	resp, _ := 	 http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes,_ := ioutil.ReadAll(resp.Body)
	var s Sitemapindex
	xml.Unmarshal(bytes,&s)	
	for _, Location := range s.Locations {
		fmt.Printf("%s\n", Location)
	}
}