FROM golang:1.15.2-alpine

# # go moduleを使用
# ENV GO111MODULE=on 

# アップデートとgitのインストール！！
RUN apk update && apk add git && go get -u github.com/oxequa/realize
# appディレクトリの作成
RUN mkdir /go/src/nfp-server
# アプリケーションを実行するディレクトリを指定
WORKDIR /go/src/nfp-server

ADD . /go/src/nfp-server
# # 上記のディレクトリにgo.modとgo.sumをコピー
# COPY go.mod go.sum ./
# # 上記のファイルに変更がなければキャッシュ利用できる
# RUN go mod download

# COPY . .
# RUN go build .

EXPOSE 8080

CMD ["realize", "start"]