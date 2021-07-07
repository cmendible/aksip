package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func main() {

	nodes := flag.Int("n", 3, "Number of nodes")
	scale := flag.Int("s", 1, "Number of scale nodes")
	maxPods := flag.Int("p", 30, "Max pods per node")
	isvc := flag.Int("l", 1, "Number of expected internal LoadBalancer services")
	table := flag.Bool("t", false, "Set true for output table")

	flag.Parse()

	r := calculator{MaxPods: *maxPods, Nodes: *nodes, Scale: *scale, Isvc: *isvc}
	r.validate()
	r.calculateRequiredIPs()

	err := r.getCIDR()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if !*table {
		renderJson(r)
	} else {
		renderTable(r)
	}

	os.Exit(0)
}

func renderTable(r calculator) {
	data := [][]string{
		{
			strconv.Itoa(r.Nodes),
			strconv.Itoa(r.Scale),
			strconv.Itoa(r.MaxPods),
			strconv.Itoa(r.Isvc),
			strconv.Itoa(r.RequiredIPs),
			r.CIDR},
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
}

func renderJson(r calculator) {
	j, _:= json.MarshalIndent(r, "", "  ")
	fmt.Println(string(j))
}
