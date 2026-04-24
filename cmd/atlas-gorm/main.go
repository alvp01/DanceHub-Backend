package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/alvp01/DanceHub-Backend/internal/academy"
	"github.com/alvp01/DanceHub-Backend/internal/student"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&academy.Academy{},
		&academy.RefreshToken{},
		&student.Student{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	_, _ = io.WriteString(os.Stdout, stmts)
}
