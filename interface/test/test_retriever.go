package test

import "fmt"

type TestRetriever struct {
	Content string
}

func (retriever TestRetriever) Get(url string) string {
	return retriever.Content
}

func (retriever TestRetriever) Post(url string, form map[string]string) string {
	retriever.Content = form["content"]
	fmt.Println(retriever.Content)
	return "ok"
}

func (tr *TestRetriever) String() string {
	return fmt.Sprintf("Retriever %v", tr.Content)
}
