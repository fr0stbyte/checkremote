package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

func main() {

	if 0 == len(*opts.url) {
		fmt.Println("checkremote needs a url to hit")
		flag.PrintDefaults()
		os.Exit(1)
	}

	result_q := make(chan *CurlResponse, *opts.capacity)

	go startChecks(result_q, *opts.url, *opts.count)
	printResults(result_q, *opts.count)
}

func startChecks(result_q chan *CurlResponse, provider string, count int) {
	for i := 0; i < count; i++ {
		<-sem
		go func() {
			checkCurl(provider, result_q)
			sem <- 1
		}()
	}
}

func printResults(result_q chan *CurlResponse, count int) {

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 10, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(w, "\tName\tDNS\tConnect\tPretransfer\tFirstByte\tTotal\tBytes\tStatus")
	for i := 0; i < count; i++ {
		res := <-result_q
		fmt.Fprintf(w, "%3d:%s\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%d\n", i, res.name, res.dns, res.connect, res.preTransfer, res.startTransfer, res.total, res.bytes, res.statusCode)
		w.Flush()
	}

}
