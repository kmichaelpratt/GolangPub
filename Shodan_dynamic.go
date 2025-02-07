package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//Important points reg script this is still in development.
//The purpose of this script is supporting a more dynaic use case IE instead of hardcoding an API key and hard coding org names, it asks that info from the user
//Once the user provides that info, the GET request is created. Still working to resolve an encoding issue though.

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

func createCSVfile(destination string) error {

	fmt.Print("Hello, What is your Shodan API key?")
	reader := bufio.NewReader(os.Stdin)

	input_api, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error received while reading input", err)
	}
	//Note that the above works good, need to figure out a solution for the array of org names

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Hello, What is your Organization Name ")
	input := multilineReader(scanner)

	fmt.Println("Input is ", input)

	//The below did not work as I would have hoped but thats fine
	//scanner := bufio.NewScanner(os.Stdin)
	//for {
	//	fmt.Print("Hello, What is your Organization Name ")
	//	// reads user input until \n by default
	//	scanner.Scan()
	//	// Holds the string that was scanned
	//	text := scanner.Text()
	//
	//	if len(text) != 0 {
	//		fmt.Println("Here are your org names")
	//		fmt.Println(text)
	//	} else {
	//		break
	//	}
	//
	//}

	//Remove the newline from the API info
	input_api = strings.TrimSuffix(input_api, "\n")

	//Printing out the API
	fmt.Println(input_api)

	//base URL
	u, _ := url.Parse("https://api.shodan.io/shodan/host/search")

	q := u.Query()
	q.Add("key", input_api)
	for _, value := range input {
		q.Add("query=org:", value)
	}
	u.RawQuery = q.Encode()

	fmt.Println(u)

	resp, err := http.Get(u.String())

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

	header := []string{"Host ASN", "Shodan Timestamp", "Host City", "Host Country", "IP address", "Port Open", "Host ISP", "Domains Listed on Host"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range resultFile.Host {
		var csvRow []string
		csvRow = append(csvRow, r.ASN, r.Timestamp, r.Location.City, r.Location.Country, r.IP_Addr, fmt.Sprint(r.Port), r.ISP, fmt.Sprint(r.Domain))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil

}

// function to read multiple lines
func multilineReader(scanner *bufio.Scanner) []string {
	input := []string{}

	for {
		//scans a line from console
		scanner.Scan()

		//Holds string that was scanned
		text := scanner.Text()

		if len(strings.TrimSpace(text)) != 0 {
			input = append(input, text)
		} else {
			break
		}
		fmt.Print(">")
	}
	return input
}

func main() {

	if err := createCSVfile("Shodan.csv"); err != nil {
		log.Fatal(err)

	}

}
