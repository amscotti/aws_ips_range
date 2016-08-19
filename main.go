package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//AWSIPs JSON format for AWS IPs Range
type AWSIPs struct {
	SyncToken  string `json:"syncToken"`
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix string `json:"ip_prefix"`
		Region   string `json:"region"`
		Service  string `json:"service"`
	} `json:"prefixes"`
}

func getIPRanges() AWSIPs {
	const url string = "https://ip-ranges.amazonaws.com/ip-ranges.json"

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var ipRanges AWSIPs
	json.Unmarshal(body, &ipRanges)

	return ipRanges
}

func main() {
	ipRanges := getIPRanges()
	fmt.Printf("%-25s%-20s%-10s\n", "IP Prefix", "Region", "Service")
	for _, ip := range ipRanges.Prefixes {
		fmt.Printf("%-25s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
	}
}
