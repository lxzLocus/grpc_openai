package main

import (
	// (一部抜粋)
	"context"
	"fmt"
	"log"
	"net"
	protocolspb "openai/pkg/grpc"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 自作サービス構造体のコンストラクタを定義
func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer()

	// 3. gRPCサーバーにOpenAIServiceを登録
	protocolspb.RegisterOpenAIServiceServer(s, NewMyServer())

	// 4. サーバーリフレクションの設定
	reflection.Register(s)

	// 5. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 5.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

type myServer struct {
	protocolspb.UnimplementedOpenAIServiceServer
}

// ナゾ多き箇所
func (s *myServer) api_convert(ctx context.Context, req *protocolspb.ChatCompletionRequest) (*protocolspb.ChatCompletionResponse, error) {
	//レスポンス箇所
	return &protocolspb.ChatCompletionResponse{}, nil
}
