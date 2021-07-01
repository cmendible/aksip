package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func main() {

	nodes := flag.Int("n", 3, "Number of nodes")
	scale := flag.Int("s", 1, "Number of scale nodes")
	maxPods := flag.Int("p", 30, "Max pods per node")
	isvc := flag.Int("l", 1, "Number of expected internal LoadBalancer services")

	flag.Parse()

	if *maxPods > 250 {
		fmt.Println("Max pods is higher than 250 (Limit per node).")
		os.Exit(1)
	}

	if *maxPods < 10 {
		fmt.Println("Max pods is lower than 10 (Minimum per node).")
		os.Exit(1)
	}

	if *maxPods**nodes < 30 {
		fmt.Println("Projected number of pods is lower than 30 (Minimum per 30 per cluster).")
		os.Exit(1)
	}

	if *nodes+*scale > 1000 {
		fmt.Println("Total number of nodes (nodes + scale) is higher than 1000 (Limit per cluster).")
		os.Exit(1)
	}

	requiredIPs := (*nodes + 1 + *scale) + ((*nodes + 1 + *scale) * *maxPods) + *isvc

	cidr := getCIDR(requiredIPs)

	data := [][]string{
		{strconv.Itoa(*nodes), strconv.Itoa(*scale), strconv.Itoa(*maxPods), strconv.Itoa(*isvc), strconv.Itoa(requiredIPs), cidr},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Nodes", "Scale", "maxPods", "isvc", "IPs", "CIDR"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	os.Exit(0)
}

func getCIDR(requiredIPs int) string {
	cidrs := map[int]int{}

	// Azure smallest supported subnet size is /29 and biggest /8
	// https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-faq#how-small-and-how-large-can-vnets-and-subnets-be 
	for i := 29; i >= 8; i-- {
		cidrs[i] = int(getAvailableHosts(i))
	}

	for cidr, ips := range cidrs {
		if ips > requiredIPs && requiredIPs > cidrs[cidr+1] {
			return fmt.Sprintf("/%s", strconv.Itoa(cidr))
		}
	}

	return ""
}

func getAvailableHosts(cidr int) float64 {
	// Remember Azure reserves 5 IPs in every subnet
	// https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-faq#are-there-any-restrictions-on-using-ip-addresses-within-these-subnets
	return math.Pow(2, float64(32-cidr)) - 5
}
