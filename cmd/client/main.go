package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "openai/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	scanner *bufio.Scanner
	client  pb.OpenAIServiceClient
)

func main() {
	fmt.Println("start gRPC Client.")

	// 1. 標準入力から文字列を受け取るスキャナを用意
	scanner = bufio.NewScanner(os.Stdin)

	// 2. gRPCサーバーとのコネクションを確立
	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()

	// 3. gRPCクライアントを生成
	client = pb.NewOpenAIServiceClient(conn)

	for {
		fmt.Println("1: send Request")
		fmt.Println("2: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			/*関数呼び出し*/
			apiConvert()
		case "2":
			fmt.Println("bye.")
			goto M
		}
	}
M:
}

func apiConvert() {
	fmt.Println("Please enter prompt.")
	scanner.Scan()
	prompt := scanner.Text()

	// タイムアウト時間を10秒に設定する例
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.ChatCompletionRequest{
		Prompt: prompt,
	}

	res, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, choice := range res.Choices {
			fmt.Println(choice.Text) // Assuming 'Text' is the field you want to print
		}
	}
}
