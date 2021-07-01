![aksip](https://github.com/cmendible/aksip/workflows/aksip/badge.svg)

# aksip

Azure Kubernetes Service (AKS) advanced networking (CNI) address space calculator.

## Download

* Download the the latest version from the [releases](https://github.com/cmendible/aksip/releases) page.

## Usage

``` shell
aksip -h
Usage of aksip:
  -l int
        Number of expected internal LoadBalancer services (default 1)
  -n int
        Number of nodes (default 3)
  -p int
        Max pods per node (default 30)
  -s int
        Number of scale nodes (default 1)
```