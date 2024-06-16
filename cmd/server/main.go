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

// api_convertにてREST APIと通信をする
func (s *myServer) ApiConvert(ctx context.Context, req *pb.ChatCompletionRequest) (*pb.ChatCompletionResponse, error) {
	// REST APIのエンドポイント
	const endpoint = "/v1/completions"
	const url = "http://192.168.10.30:5000" + endpoint

	// リクエストペイロードの作成
	payload := map[string]interface{}{
		"model":                      req.Model,
		"prompt":                     req.Prompt,
		"best_of":                    req.BestOf,
		"echo":                       req.Echo,
		"frequency_penalty":          req.FrequencyPenalty,
		"logit_bias":                 req.LogitBias,
		"logprobs":                   req.Logprobs,
		"max_tokens":                 req.MaxTokens,
		"n":                          req.N,
		"presence_penalty":           req.PresencePenalty,
		"stop":                       req.Stop,
		"stream":                     req.Stream,
		"suffix":                     req.Suffix,
		"temperature":                req.Temperature,
		"top_p":                      req.TopP,
		"user":                       req.User,
		"preset":                     req.Preset,
		"min_p":                      req.MinP,
		"dynamic_temperature":        req.DynamicTemperature,
		"dynatemp_low":               req.DynatempLow,
		"dynatemp_high":              req.DynatempHigh,
		"dynatemp_exponent":          req.DynatempExponent,
		"smoothing_factor":           req.SmoothingFactor,
		"top_k":                      req.TopK,
		"repetition_penalty":         req.RepetitionPenalty,
		"repetition_penalty_range":   req.RepetitionPenaltyRange,
		"typical_p":                  req.TypicalP,
		"tfs":                        req.Tfs,
		"top_a":                      req.TopA,
		"epsilon_cutoff":             req.EpsilonCutoff,
		"eta_cutoff":                 req.EtaCutoff,
		"guidance_scale":             req.GuidanceScale,
		"negative_prompt":            req.NegativePrompt,
		"penalty_alpha":              req.PenaltyAlpha,
		"mirostat_mode":              req.MirostatMode,
		"mirostat_tau":               req.MirostatTau,
		"mirostat_eta":               req.MirostatEta,
		"temperature_last":           req.TemperatureLast,
		"do_sample":                  req.DoSample,
		"seed":                       req.Seed,
		"encoder_repetition_penalty": req.EncoderRepetitionPenalty,
		"no_repeat_ngram_size":       req.NoRepeatNgramSize,
		"min_length":                 req.MinLength,
		"num_beams":                  req.NumBeams,
		"length_penalty":             req.LengthPenalty,
		"early_stopping":             req.EarlyStopping,
		"truncation_length":          req.TruncationLength,
		"max_tokens_second":          req.MaxTokensSecond,
		"prompt_lookup_num_tokens":   req.PromptLookupNumTokens,
		"custom_token_bans":          req.CustomTokenBans,
		"sampler_priority":           req.SamplerPriority,
		"auto_max_new_tokens":        req.AutoMaxNewTokens,
		"ban_eos_token":              req.BanEosToken,
		"add_bos_token":              req.AddBosToken,
		"skip_special_tokens":        req.SkipSpecialTokens,
		"grammar_string":             req.GrammarString,
	}

	// ペイロードをJSONに変換
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// HTTPクライアントを作成
	client := &http.Client{Timeout: 10 * time.Second}

	// HTTPリクエストを送信
	resp, err := client.Do(req)
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
		return nil, fmt.Errorf("API error: %v", apiError)
	}
}
