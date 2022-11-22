package main

import (
	"fmt"
	"math"
	"net/url"
	"strings"
)

func main() {
	s := "postgres://user:pass@host.com:5432/path?k=v#f"

	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(u.Scheme)

	fmt.Println(u.User)
	fmt.Println(u.User.Username())
	p, _ := u.User.Password()
	fmt.Println(p)

	fmt.Println(u.Host)
	h := strings.Split(u.Host, ":")
	fmt.Println(h[0])
	fmt.Println(h[1])

	fmt.Println(u.Path)
	fmt.Println(u.Fragment)

	fmt.Println(u.RawQuery)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m)
	fmt.Println(m["k"][0])

	thumbnail := "https://vecarbon-p.cn-sh2.ufileos.com/d/liuting/1471841930349401956/1471841930349401947/thumbnail_1666070897412.jpeg"
	fileName := strings.Split(thumbnail, "/")
	fmt.Println(fileName[len(fileName)-1])

	fmt.Println(math.Trunc(9.815*1e2+0.5) * 1e-2)
}
