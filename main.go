package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	var appendMode bool
	flag.BoolVar(&appendMode, "a", false, "Append the value instead of replacing it")
	flag.Parse()

	seen := make(map[string]url.Values)

	// read URLs on stdin, then replace the values in the query string
	// with some user-provided value
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", sc.Text(), err)
			continue
		}

		var key string

		if ((u.Scheme == "http" || u.Scheme == "https" ) && u.Port() == "") {
			key = fmt.Sprintf("%s://%s%s", u.Scheme, u.Hostname(), u.EscapedPath())
		} else {
			if(u.Port() == "80" || u.Port() == "443") {
				key = fmt.Sprintf("%s://%s%s", u.Scheme, u.Hostname(), u.EscapedPath())
			} else {
				key = fmt.Sprintf("%s://%s:%s%s", u.Scheme, u.Hostname(), u.Port(), u.EscapedPath())
			}
		}

		if _, exists := seen[key]; !exists {
			seen[key] = make(url.Values)
		}

		for k,v := range u.Query() {
			seen[key][k] = v
		}
	}

	for key, vals := range seen {
		u, err := url.Parse(key);
		if err != nil {
			fmt.Fprintf(os.Stderr, "Impossible, failed to parse url %s [%s]\n", sc.Text(), err)
			continue
		}
		qs := url.Values{}

		for param, vv := range vals {
			if appendMode {
				qs.Set(param, vv[0]+flag.Arg(0))
			} else {
				qs.Set(param, flag.Arg(0))
			}
		}

		u.RawQuery = qs.Encode()

		fmt.Printf("%s\n", u)

	}

}
