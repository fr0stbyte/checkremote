package main

import (
	"flag"
)

type CheckproviderOpts struct {
	agent    *string
	count    *int
	capacity *int
	headers  *string
	http11   *bool
	url      *string
}

var opts CheckproviderOpts
var sem chan int

func init() {
	// initialize option flags
	opts.url = flag.String("u", "", "url to hit")
	opts.count = flag.Int("n", 10, "times to run the test")
	opts.capacity = flag.Int("c", 10, "simultaneous requests")
	opts.http11 = flag.Bool("http11", false, "use HTTP/1.0 by default")
	opts.headers = flag.String("H", "", "extra headers to pass")
	opts.agent = flag.String("A", "checkremote/0.0.1", "user agent string")
	flag.Parse()

	// upper bound for opts.capacity
	if *opts.capacity > 20 {
		*opts.capacity = 10
	}

	sem = make(chan int, *opts.capacity)
	for i := 0; i < *opts.capacity; i++ {
		sem <- 1
	}

}
