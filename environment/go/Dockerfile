#デプロイ用コンテナに含めるバイナリを作成するコンテナ
FROM golang:1.22rc2-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -trimpath -ldflags "-w -s" -o app ./api/cmd/main.go

#-----------------------------------------------
#デプロイ用コンアテナ
FROM centos:centos7 as deploy

RUN yum -y update

EXPOSE "8080"

COPY --from=deploy-builder /app/app .

CMD ["./app"]

#-----------------------------------------------
#ローカル開発環境で利用するホットリロード環境
FROM golang:1.22.0-bullseye as dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go install github.com/cosmtrek/air@latest
CMD ["air"]