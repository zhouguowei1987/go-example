package main

import (
	"fmt"
	s "strings"
)

var p45 = fmt.Println

func main() {
	p45("Contains:", s.Contains("test", "t"))
	p45("Count:", s.Count("test", "t"))
	p45("HasPrefix:", s.HasPrefix("test", "te"))
	p45("HasSuffix:", s.HasSuffix("test", "st"))
	p45("Index:", s.Index("test", "e"))
	p45("Join:", s.Join([]string{"a", "b"}, "-"))
	p45("Repeat:", s.Repeat("a", 5))
	p45("Replace:", s.Replace("foo", "o", "0", -1))
	p45("Replace:", s.Replace("foo", "o", "0", 1))
	p45("Split:", s.Split("a-b-c-d-e", "-"))
	p45("ToLower:", s.ToLower("TEST"))
	p45("ToUpper", s.ToUpper("test"))

	p45("Len:", len("hello"))
	p45("Char", "hello"[1])
}
