package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const url string = "https://ip-ranges.amazonaws.com/ip-ranges.json"

//AWSIPs JSON format for AWS IPs Range
type Prefix struct {
	IPPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

func (ip *Prefix) print() {
	fmt.Printf("%-25s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
}

func (ip *Prefix) toRow() []string {
	return []string{ip.IPPrefix, ip.Region, ip.Service}
}

type IPv6Prefix struct {
	IPPrefix           string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

func (ip *IPv6Prefix) print() {
	fmt.Printf("%-25s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
}

func (ip *IPv6Prefix) toRow() []string {
	return []string{ip.IPPrefix, ip.Region, ip.Service}
}

type AWSIPs struct {
	SyncToken    string       `json:"syncToken"`
	CreateDate   string       `json:"createDate"`
	Prefixes     []Prefix     `json:"prefixes"`
	IPv6Prefixes []IPv6Prefix `json:"ipv6_prefixes"`
}

func (ipRanges *AWSIPs) print() {
	fmt.Printf("%-25s%-20s%-10s\n", "IP Prefix", "Region", "Service")
	for _, ip := range ipRanges.Prefixes {
		ip.print()
	}

	fmt.Printf("\n%-25s%-20s%-10s\n", "IPv6 Prefix", "Region", "Service")
	for _, ip := range ipRanges.IPv6Prefixes {
		ip.print()
	}
}

func (ipRanges *AWSIPs) toCSV() {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	if err := w.Write([]string{"IP Prefix", "Region", "Service"}); err != nil {
		log.Fatalln("error writing header to file", err)
	}

	for _, ip := range ipRanges.Prefixes {
		if err := w.Write(ip.toRow()); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	for _, ip := range ipRanges.IPv6Prefixes {
		if err := w.Write(ip.toRow()); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

}

func getIPRanges() (AWSIPs, error) {
	var ipRanges AWSIPs

	res, err := http.Get(url)
	if err != nil {
		return ipRanges, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ipRanges, err
	}

	if err := json.Unmarshal(body, &ipRanges); err != nil {
		return ipRanges, err
	}

	return ipRanges, nil
}

func main() {
	var printCSV bool

	flag.BoolVar(&printCSV, "csv", false, "Output as CSV")
	flag.Parse()

	ipRanges, err := getIPRanges()
	if err != nil {
		log.Fatal(err)
	}

	if printCSV {
		ipRanges.toCSV()
	} else {
		ipRanges.print()
	}
}
