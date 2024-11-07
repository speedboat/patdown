package common

import (
	"fmt"
	"os"
)

func scan(nameservers []Nameserver, threads, delay int, recursive, single bool) {
	queries := make(chan Query)
	tab := make(chan interface{})

	if !recursive {
		Info(fmt.Sprintf("performing non-recursive lookups against %d resolvers...", len(nameservers)))
		for i := 0; i < threads; i++ {
			go RunQuery(queries, tab, delay)
		}

		for _, ns := range nameservers {
			for vendor, domains := range Vendors {
				for _, domainpair := range domains {
					queries <- Query{Nameserver: ns.Nameserver, Vendor: vendor, DomainPair: domainpair}
				}
			}
		}
	} else {
		Warn("recursive snooping can only be done once, as it populates the nameserver's cache")
		Info(fmt.Sprintf("recursively snooping on %d resolvers...", len(nameservers)))
		for i := 0; i < threads; i++ {
			go RunQueryRA(queries, tab, delay)
		}

		if !single {
			for _, ns := range nameservers {
				for vendor, domains := range Vendors {
					for _, domainpair := range domains {
						queries <- Query{Nameserver: ns.Nameserver, Vendor: vendor, DomainPair: domainpair}
					}
				}
			}
		} else {
			for vendor, domains := range Vendors {
				for _, domainpair := range domains {
					queries <- Query{Nameserver: nameservers[0].Nameserver, Vendor: vendor, DomainPair: domainpair}
				}
			}
		}
	}

	close(queries)
}

func Takeoff(nameservers []Nameserver) {
	var nonrns, rns []Nameserver
	for _, ns := range nameservers {
		if ns.Recursive {
			rns = append(rns, ns)
		}
		if ns.NonRA {
			nonrns = append(nonrns, ns)
		}
	}

	if len(nonrns) == 0 && len(rns) == 0 {
		Fatal("no valid nameservers available for probing, they may be down or they don't like your IP")
	}

	recursive := false

	for {
		if !recursive {
			if len(nonrns) > 0 {
				scan(nonrns, Params.Threads, Params.Delay, false, false)
			} else {
				for {
					Info(fmt.Sprintf("non-recursive lookups not viable on these servers, perform recursive snooping? %s(less reliable, can only be done once per server)%s",
						ColorRed, ColorReset))
					fmt.Printf("%s `--(y/n):%s ", ColorCyan, ColorReset)
					var input string
					fmt.Scanln(&input)
					if input == "y" {
						recursive = true
						break
					}
					if input == "n" {
						os.Exit(0)
					}
				}
				continue
			}
		} else {
			autodetected := Params.Domain != "" && len(Params.Nservers) == 0
			scan(rns, Params.Threads, Params.Delay, true, autodetected)
		}

	}
}
