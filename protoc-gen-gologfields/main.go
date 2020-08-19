// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package main

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/improbable-io/go-proto-logfields/internal/genlogfields"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if f.Generate {
				genlogfields.GenerateFile(gen, f)
			}
		}
		return nil
	})
}
