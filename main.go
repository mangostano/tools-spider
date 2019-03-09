package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func main(){
	versionPageNumber := [6]string{"1.0", "1.1", "2.0", "2.1", "2.2", "3.0"}
	resultMap := make(map[string][]string)
	c := colly.NewCollector()

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		url := e.Request.URL.Path
		version := url[len(url)-3:]
		subVersions := make([]string, 0)
		e.ForEach("tr > td:nth-child(2) > h3", func(i int, element *colly.HTMLElement) {
			subVersions = append(subVersions, element.Text[4:])
		})
		resultMap[version] = subVersions
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	for _, page := range versionPageNumber {
		c.Visit("https://dotnet.microsoft.com/download/dotnet-core/" + page)
    }

	c.Wait()
	resultJson, _ := json.Marshal(resultMap)
	fmt.Println(string(resultJson))
	jsonFile, err := os.Create("./versions.json")
	if err != nil {
		log.Fatal("create json file failed!")
	}
	defer jsonFile.Close()
	jsonFile.WriteString(string(resultJson))
	jsonFile.Sync()
}

