package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	pb "openai/pkg/grpc"

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
	pb.RegisterOpenAIServiceServer(s, NewMyServer())

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
	pb.UnimplementedOpenAIServiceServer
}

type APIError struct {
	Code int
	Msg  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: code = %d, msg = %s", e.Code, e.Msg)
}

// api_convertにてREST APIと通信をする
func (s *myServer) CreateChatCompletion(ctx context.Context, req *pb.ChatCompletionRequest) (*pb.ChatCompletionResponse, error) {
	// REST APIのエンドポイント
	const endpoint = "/v1/completions"
	const url = "http://192.168.10.30:5000" + endpoint

	// リクエストペイロードの作成
	promptPayload := map[string]interface{}{
		"prompt":                     req.Prompt,
		"max_tokens":                 512,
		"temperature":                0.7,
		"temperature_last":           false,
		"dynamic_temperature":        false,
		"dynatemp_low":               1,
		"dynatemp_high":              1,
		"dynatemp_exponent":          1,
		"smoothing_factor":           0,
		"top_p":                      0.9,
		"min_p":                      0,
		"top_k":                      20,
		"repetition_penalty":         1.15,
		"presence_penalty":           0,
		"frequency_penalty":          0,
		"repetition_penalty_range":   1024,
		"typical_p":                  1,
		"tfs":                        1,
		"top_a":                      0,
		"epsilon_cutoff":             0,
		"eta_cutoff":                 0,
		"guidance_scale":             1,
		"penalty_alpha":              0,
		"mirostat_mode":              0,
		"mirostat_tau":               5,
		"mirostat_eta":               0.1,
		"do_sample":                  true,
		"seed":                       -1,
		"encoder_repetition_penalty": 1,
		"no_repeat_ngram_size":       0,
		"min_length":                 0,
		"num_beams":                  1,
		"length_penalty":             1,
		"early_stopping":             false,
		"sampler_priority":           "temperature\ndynamic_temperature\nquadratic_sampling\ntop_k\ntop_p\ntypical_p\nepsilon_cutoff\neta_cutoff\ntfs\ntop_a\nmin_p\nmirostat",
	}

	// ペイロードをJSONに変換
	payload, err := json.Marshal(promptPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	// デバッグ用にリクエスト内容を出力
	log.Println("HTTP Request Payload:")
	log.Println(string(payload))

	// HTTPリクエストの作成
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// HTTPクライアントを作成
	client := &http.Client{Timeout: 20 * time.Second}

	// HTTPリクエストを送信
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// ステータスコードに応じた処理
	if resp.StatusCode == http.StatusOK {
		// 200 OK の場合の処理
		var apiResp pb.ChatCompletionResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
		}
		return &apiResp, nil
	} else {
		// エラーの場合の処理
		var apiError pb.ErrorResponse
		if err := json.Unmarshal(body, &apiError); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error response body: %v", err)
		}
		return nil, &APIError{Code: resp.StatusCode, Msg: apiError.String()} // カスタムエラー型を使用
	}

}
