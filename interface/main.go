package main

import (
	"fmt"
	"iflytek.com/weipan4/learn-go/interface/infra"
	"iflytek.com/weipan4/learn-go/interface/real"
	"iflytek.com/weipan4/learn-go/interface/test"
)

type Retriever interface {
	Get(url string) string
}

type Poster interface {
	Post(url string, form map[string]string) string
}

type RetrieverPoster interface {
	Retriever
	Poster
}

func download(r Retriever) string {
	return r.Get("https://www.baidu.com")
}

const url = "https://www.imooc.com"

func session(s RetrieverPoster) string {
	s.Post(url, map[string]string{
		"content": "another fake imooc.com",
	})
	return s.Get(url)
}

func inspect(r Retriever) {
	switch v := r.(type) {
	case infra.UrlRetriever:
		break
	case real.Retriever:
		fmt.Println(v.UserAgent, v.TimeOut)
	case test.TestRetriever:
		fmt.Println(v.Content)
	}
}

func main() {
	var r Retriever
	r = infra.UrlRetriever{}
	fmt.Println(download(r))
	r = test.TestRetriever{Content: "this is fake interface"}
	fmt.Println(download(r))
	r = real.Retriever{}
	fmt.Println(download(r))
	fmt.Printf("%T %v\n", r, r)

	inspect(r)

	r = &real.Retriever{
		UserAgent: "Mozilla/5.0",
		TimeOut:   0,
	}
	retriever := r.(*real.Retriever)
	fmt.Printf("%T %v\n", retriever, retriever)

	if testRetriever, ok := r.(test.TestRetriever); ok {
		fmt.Println(testRetriever.Content)
	} else {
		fmt.Println("not a test interface")
	}

	fmt.Println("try a session")
	rp := test.TestRetriever{Content: "this is fake interface"}
	fmt.Println(session(rp))
}
