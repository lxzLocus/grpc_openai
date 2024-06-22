package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

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
			callApi()
		case "2":
			fmt.Println("bye.")
			return
		}
	}
}

func callApi() {
	fmt.Println("Please enter prompt.")
	scanner.Scan()
	prompt := scanner.Text()

	req := &pb.ChatCompletionRequest{
		Prompt: prompt,
	}

	res, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Println("Error calling ApiConvert:", err)
	} else {
		for _, choice := range res.GetChoices() {
			fmt.Println("Response from CreateChatCompletion:", choice.GetText())
		}
	}
}
