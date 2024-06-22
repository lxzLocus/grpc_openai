# 手順系
---

## protocコマンドでコードを生成する
```
$ cd api
$ protoc --go_out=../pkg/grpc --go_opt=paths=source_relative --go-grpc_out=../pkg/grpc --go-grpc_opt=paths=source_relative protocols.proto
```

### TIPS
```
protocコマンドにつけていたオプションはそれぞれ以下の通りです。

--go_out: hello.pb.goファイルの出力先ディレクトリを指定
--go_opt: hello.pb.goファイル生成時のオプション。
今回はpaths=source_relativeを指定して--go_outオプションでの指定が相対パスであることを明示

--go-grpc_out: hello_grpc.pb.goファイルの出力先ディレクトリを指定
--go-grpc_opt: hello_grpc.pb.goファイル生成時のオプション。
今回はpaths=source_relativeを指定して--go-grpc_outオプションでの指定が相対パスであることを明示
```

## サーバーの起動
```
$ cd cmd/server
$ go run main.go
```

## サーバー内に実装されているサービス一覧の確認
```
$ grpcurl -plaintext localhost:8080 list
```

## サービスのメソッド一覧確認
```
$ grpcurl -plaintext localhost:8080 list openai.OpenAIService
```

## メソッド呼び出し
```
$ grpcurl -plaintext -d '{"prompt": "create simple python sample code"}' localhost:8080 openai.OpenAIService
```

---
# 呼び出し方法 grpcurl

## 悪い呼び出し
`$ grpcurl -plaintext -d '{"prompt": "create simple python sample code"}' localhost:8080 openai.OpenAIService`

## 正しい呼び出し
`$ grpcurl -plaintext -d '{"prompt": "create simple python sample code"}' localhost:8080 openai.OpenAIService/CreateChatCompletion`

### TIPS
max_tokens大きくすると，回答生成に時間がかかって，gRPCがタイムアウトする
```
% grpcurl -plaintext -d '{"prompt": "create simple python sample code"}' localhost:8080 openai.OpenAIService/CreateChatCompletion

ERROR:
  Code: Unknown
  Message: failed to send HTTP request: Post "http://192.168.10.30:5000/v1/completions": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```