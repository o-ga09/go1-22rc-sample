package batch

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/storage"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

const CSV_RECORD = 1000000

func Run(ctx context.Context) {
	// データベースへの接続
	db, err := sql.Open("mysql", "api:P@ssw0rd@tcp(localhost:3306)/api?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	data := GetCSV(ctx)

	// データベースにインサート
	sql := "INSERT INTO users (name, email) VALUES"
	for _, record := range *data {
		sql += fmt.Sprintf("('%v','%v'),", record[2], record[1])
	}
	sql = sql[:len(sql)-1]
	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

}

func UploadCSV(ctx context.Context) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	bucketName := "sample-api-batch-20240220"
	objectName := "csv/test_2.csv"

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
func GetCSV(ctx context.Context) *[][]string {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	bucketName := "sample-api-batch-20240220"
	objectName := "csv/test_2.csv"

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(rc)
	reader.Read()

	var data [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// 各レコードを結合
		data = append(data, record)
	}
	return &data
}

func CreateCSV(ctx context.Context) {
	// 構造体のリスト
	data := []struct {
		id    string
		name  string
		email string
	}{}

	for i := range CSV_RECORD {
		d := struct {
			id    string
			name  string
			email string
		}{
			id: strconv.Itoa(i), name: fmt.Sprintf("test%d", i), email: uuid.NewString(),
		}
		data = append(data, d)
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
