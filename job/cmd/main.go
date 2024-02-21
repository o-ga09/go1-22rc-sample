package main

import (
	"context"
	"fmt"
	"time"

	"github.com/o-ga09/go122rcsample/job/batch"
)

func main() {
	ctx := context.Background()
	start := time.Now()
	fmt.Println("start")
	batch.CreateCSV(ctx)
	fmt.Println("end")
	elapsed := time.Since(start)

	fmt.Println(elapsed)
}
