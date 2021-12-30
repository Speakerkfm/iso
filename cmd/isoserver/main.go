package main

import (
	"fmt"
	"plugin"

	_ "google.golang.org/protobuf/proto"

	"github.com/Speakerkfm/iso/pkg/models"
)

func main() {
	fmt.Println("Hello world!")

	plug, err := plugin.Open("example/proto/proto_struct.so")
	if err != nil {
		panic(err)
	}

	svcs, err := plug.Lookup("Services")
	if err != nil {
		panic(err)
	}

	s, ok := svcs.([]*models.ProtoService)
	if !ok {
		panic(":(")
	}

	fmt.Printf("svc: %+v\n", s)
}
