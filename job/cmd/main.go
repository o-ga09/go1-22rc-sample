package main

import (
	"context"
	"fmt"

	"github.com/o-ga09/go122rcsample/job/batch"
)

func main() {
	ctx := context.Background()
	fmt.Println("start")
	// batch.UploadCSV(ctx)
	batch.GetCSV(ctx)
	fmt.Println("end")
}
