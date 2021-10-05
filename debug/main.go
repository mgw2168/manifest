package main

import (
	"fmt"
	"time"
)

func main() {
	format := "2006-01-02 15:04:05"
	now := time.Now()
	//now, _ := time.Parse(format, time.Now().Format(format))
	a, _ := time.Parse(format, "2015-03-10 11:00:00")
	b, _ := time.Parse(format, "2015-03-10 16:00:00")

	fmt.Println(now.Format(format), a.Format(format), b.Format(format))
	fmt.Println(now.After(a))
	fmt.Println(now.Before(a))
	fmt.Println(now.After(b))
	fmt.Println(now.Before(b))
	fmt.Println(a.After(b))
	fmt.Println(a.Before(b))
	fmt.Println(now.Format(format), a.Format(format), b.Format(format))
	fmt.Println(now.Unix(), a.Unix(), b.Unix())
}
