package main

import (
	"fmt"
	"log"
	"os/exec"
)

const name = "my name is %s"

func main() {
	var output []byte
	cmd := exec.Command("pwd")
	if output, err := cmd.Output(); err != nil {
		log.Fatalf("cmd %s failed: %v", string(output), err)
	}
	fmt.Println(string(output))
	fmt.Println(fmt.Sprintf(name, "weipan4"))
}
