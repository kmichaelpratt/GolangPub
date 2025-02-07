package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Host []Host `json:"matches"`
}

type Host struct {
	ASN       string   `json:"asn"`
	Host      []string `json:"hostnames"`
	Location  Location `json:"location"`
	Org       string   `json:"org"`
	Port      int      `json:"port"`
	IP_Addr   string   `json:"ip_str"`
	Domain    []string `json:"domains"`
	Timestamp string   `json:"timestamp"`
	ISP       string   `json:"isp"`
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country_name"`
}

//Important points reg script this is GA and ready for production
//Note in line 40 a GET request is issued, replace API_KEY with your Shodan API key.
//Also replace "ORG_NAME" with the name of your organization, if only one name remove "ORG_NAME_B", if many different orgnames replicate syntax to reflect 'ORG_NAME_A','ORG_NAME_B,'ORG_NAME_C' etc.
// Once the script completes a csv file named "Shodan.csv"  will be created in the directory where the script is located.

func createCSVfile(destination string) error {

	resp, err := http.Get("https://api.shodan.io/shodan/host/search?key=API_KEY&query=org:'ORG_NAME_A','ORG_NAME_B'&facets=ip&limit=800")

	if err != nil {
		fmt.Println("Error received", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte
	//fmt.Println(string(body))          // convert to string before print
	// note based on the above we know that the GET request is working properly
	//if you uncomment out the above fmt.Println statement it will show the Get request is working

	var resultFile Response
	if err := json.Unmarshal(body, &resultFile); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	//Create file to store CSV data
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	//write the header of the CSV file and the successive rows by iterating thru JSON struct array
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"Host ASN", "Shodan Timestamp", "Org Name", "Host City", "Host Country", "IP address", "Port Open", "Host ISP", "Domains Listed on Host"}
	if err := writer.Write(header); err != nil {
		return err
	}

	//Range thru results and add the pertinent info for each host IE their port, ASN, Org_Name etc
	for _, r := range resultFile.Host {
		var csvRow []string
		csvRow = append(csvRow, r.ASN, r.Timestamp, r.Org, r.Location.City, r.Location.Country, r.IP_Addr, fmt.Sprint(r.Port), r.ISP, fmt.Sprint(r.Domain))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil

}

func main() {

	if err := createCSVfile("Shodan.csv"); err != nil {
		log.Fatal(err)

	}

}
