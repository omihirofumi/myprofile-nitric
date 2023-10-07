package main

import (
	"fmt"
	"github.com/nitrictech/go-sdk/nitric"
)

func main() {
	profilesApi, err := nitric.NewApi("public")
	if err != nil {
		return
	}

	profiles, err := nitric.NewCollection("profiles").With(nitric.CollectionReading, nitric.CollectionWriting)
	if err != nil {
		return
	}

	if err := nitric.Run(); err != nil {
		fmt.Println(err)
	}
}
