package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/fatih/color"
)

type Area struct {
	Id             int    `json:"id"`
	DisplayName    string `json:"displayName"`
	Areas          []Area `json:"areas"`
	TotalConfirmed uint32 `json:"TotalConfirmed"`
	TotalDeaths    uint32 `json:"TotalDeaths"`
	TotalRecovered uint32 `json:"TotalRecovered"`
}

type Data struct {
	Id             string `json:"id"`
	DisplayName    string `json:"displayName"`
	Areas          []Area `json:"areas"`
	TotalConfirmed uint32 `json:"totalConfirmed"`
	TotalDeaths    uint32 `json:"totalDeaths"`
	TotalRecovered uint32 `json:"totalRecovered"`
	LastUpdated    string `json:"lastUpdated"`
}

func main() {
	response, err := http.Get("https://bing.com/covid/data")

	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err.Error)
	}

	var covid Data
	json.Unmarshal(body, &covid)

	sort.SliceStable(covid.Areas, func(i, j int) bool {
		return covid.Areas[i].TotalConfirmed > covid.Areas[j].TotalConfirmed
	})

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	w := tabwriter.NewWriter(os.Stdout, 35, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "Ülke Sayısı", yellow("Toplam Hasta Sayısı"), red("Toplam Ölüm"), green("Toplam İyileşen"))
	w.Flush()

	for _, area := range covid.Areas {
		w := tabwriter.NewWriter(os.Stdout, 35, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", area.DisplayName, yellow(area.TotalConfirmed), red(area.TotalDeaths), green(area.TotalRecovered))
		w.Flush()
	}

	tw := tabwriter.NewWriter(os.Stdout, 35, 0, 20, ' ', tabwriter.TabIndent)
	fmt.Fprintf(tw, "\n\n%s\t%s\t%s\n", yellow("Dünya Geneli Toplam Hasta Sayısı"), red("Dünya Geneli Toplam Ölüm"), green("Dünya Geneli Toplam İyileşen"))
	fmt.Fprintf(tw, "%s\t%s\t%s\n", yellow(covid.TotalConfirmed), red(covid.TotalDeaths), green(covid.TotalRecovered))
	tw.Flush()
}
