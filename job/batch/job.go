package batch

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

func UploadCSV(ctx context.Context) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	bucketName := "sample-api-batch-20240220"
	objectName := "csv/test.csv"

	writer := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	// test.csvを取得して、cloud storage にアップロードする
	f, err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 書き込み
	if _, err = io.Copy(writer, f); err != nil {
		log.Fatal(err)
	}
}

// GCSに保存されたcsvを取得する
func GetCSV(ctx context.Context) *[]byte {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	bucketName := "sample-api-batch-20240220"
	objectName := "csv/test.csv"

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", data)
	return &data
}

func CreateCSV(ctx context.Context) {
	// 構造体のリスト
	data := []struct {
		id    string
		name  string
		email string
	}{
		{id: "value1", name: "value2", email: "value3"},
		{id: "value4", name: "value5", email: "value6"},
	}

	// CSVファイルの作成
	f, err := os.Create("test.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// ヘッダーの書き込み
	if err := w.Write([]string{"id", "name", "email"}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 構造体のリストの書き込み
	for _, d := range data {
		if err := w.Write([]string{d.id, d.name, d.email}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
