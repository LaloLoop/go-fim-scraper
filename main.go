package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type LoginParams struct {
	viewState     string
	viewGenerator string
	username      string
	password      string
}

type Location struct {
	Lat  string
	Long string
}

type LocationRecord struct {
	ReportedSpeed         string
	AccumulatedKilometers string
	VehicleState          string
	ReportDate            string
	Location              Location
}

func main() {
	fName := "location-records.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}

	defer file.Close()

	c := colly.NewCollector()

	locationsCollector := c.Clone()

	// TODO replace for env vars for credentials
	loginParams := LoginParams{username: "<your-username>", password: "<your-password>"}
	locationRecords := make([]LocationRecord, 0)
	url := "https://fimtrack.com"

	c.OnHTML("input", func(e *colly.HTMLElement) {
		id := e.Attr("id")

		if id == "__VIEWSTATE" {
			loginParams.viewState = e.Attr("value")
		} else if id == "__VIEWSTATEGENERATOR" {
			loginParams.viewGenerator = e.Attr("value")
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	err = c.Visit(url)

	if err != nil {
		log.Fatal("Failed to visit URL", err)
	}

	log.Println("Go Login params", loginParams)

	loginFormData := map[string]string{
		"user_login":           loginParams.username,
		"user_pass":            loginParams.password,
		"__VIEWSTATE":          loginParams.viewState,
		"__VIEWSTATEGENERATOR": loginParams.viewGenerator,
		"btSubmit":             "Entrar",
	}

	err = locationsCollector.Post(fmt.Sprintf("%s/login.aspx", url), loginFormData)

	if err != nil {
		log.Fatal("Failed to login", err)
	}

	locationsCollector.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})

	locationsCollector.OnError(func(response *colly.Response, err error) {
		log.Fatalf("Failed to perform request to %s (HttpCode: %d): %s\n", response.Request.URL.String(), response.StatusCode, err)
	})

	locationsCollector.OnHTML("table[id=tableloc]", func(e *colly.HTMLElement) {
		log.Println("Table with locations found")

		e.ForEach("tbody > tr", func(_ int, el *colly.HTMLElement) {
			log.Println("Found tr element")

			var locationRecord LocationRecord

			el.ForEach("td", func(idx int, td *colly.HTMLElement) {
				log.Printf("Found td element at idx %d\n", idx)

				switch idx {
				case 1:
					locationRecord.ReportedSpeed = td.Text
				case 2:
					locationRecord.AccumulatedKilometers = td.Text
				case 3:
					locationRecord.VehicleState = td.ChildText("span")
				case 4:
					locationRecord.ReportDate = td.Text
				case 5:
					locationRecord.Location.Lat = td.ChildAttr("div a", "lat")
					locationRecord.Location.Long = td.ChildAttr("div a", "lon")
				}
			})

			locationRecords = append(locationRecords, locationRecord)
		})
	})

	err = locationsCollector.Visit(fmt.Sprintf("%s/recorrido_c.aspx?idMenu=4", url))

	if err != nil {
		log.Fatal("Failed to retrieve location records", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")

	err = enc.Encode(locationRecords)

	if err != nil {
		log.Fatal("Failed to encode location records as JSON", err)
	}

	err = locationsCollector.Request("GET", fmt.Sprintf("%s/doLogOut.aspx", url), nil, nil, nil)

	if err != nil {
		log.Fatal("Failed to logout", err)
	}
}
