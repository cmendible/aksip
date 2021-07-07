package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type calculator struct {
	Nodes       int `json:"nodes"`
	Scale       int `json:"scale"`
	MaxPods     int `json:"maxPods"`
	Isvc        int `json:"isvc"`
	RequiredIPs int `json:"requiredIPs"`
	CIDR        string `json:"cidr"`
}

func (r *calculator) validate() {
	if r.MaxPods > 250 {
		fmt.Println("Max pods is higher than 250 (Limit per node).")
		os.Exit(1)
	}

	if r.MaxPods < 10 {
		fmt.Println("Max pods is lower than 10 (Minimum per node).")
		os.Exit(1)
	}

	if r.MaxPods*r.Nodes < 30 {
		fmt.Println("Projected number of pods is lower than 30 (Minimum per 30 per cluster).")
		os.Exit(1)
	}

	if r.Nodes+r.Scale > 1000 {
		fmt.Println("Total number of nodes (nodes + scale) is higher than 1000 (Limit per cluster).")
		os.Exit(1)
	}
}

func (r *calculator) calculateRequiredIPs() {
	r.RequiredIPs = (r.Nodes + 1 + r.Scale) + ((r.Nodes + 1 + r.Scale) * r.MaxPods) + r.Isvc
}

func (r *calculator) getCIDR() error {
	cidrs := map[int]int{}

	// Azure smallest supported subnet size is /29 and biggest /8
	// https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-faq#how-small-and-how-large-can-vnets-and-subnets-be
	for i := 29; i >= 8; i-- {
		cidrs[i] = int(getAvailableHosts(i))
	}

	for cidr, ips := range cidrs {
		if ips > r.RequiredIPs && r.RequiredIPs > cidrs[cidr+1] {
			r.CIDR = fmt.Sprintf("/%s", strconv.Itoa(cidr))
			return nil
		}
	}

	return fmt.Errorf("no CIDR found")
}

func getAvailableHosts(cidr int) float64 {
	// Remember Azure reserves 5 IPs in every subnet
	// https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-faq#are-there-any-restrictions-on-using-ip-addresses-within-these-subnets
	return math.Pow(2, float64(32-cidr)) - 5
}
