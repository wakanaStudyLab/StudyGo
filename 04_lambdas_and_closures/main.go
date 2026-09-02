package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// Step 4: Go の無名関数・クロージャ・実践パターン
// 実行: go run ./04_lambdas_and_closures/
// ============================================================================

func main() {
	demoBasicAnonymousFunctions()
	demoClosuresAndState()
	demoLoopVariableCapture()
	demoDeferWithAnonymousFunction()
	demoFunctionalOptionsPattern()
	demoHigherOrderFunctions()
}

// ----------------------------------------------------------------------------
// 1. 基本的な無名関数と即時実行 (IIFE)
// ----------------------------------------------------------------------------
func demoBasicAnonymousFunctions() {
	fmt.Println("=== 1. Anonymous Functions & IIFE ===")

	// 変数への関数代入
	add := func(a, b int) int {
		return a + b
	}
	fmt.Printf("Add(15, 25) = %d\n", add(15, 25))

	// 即時実行無名関数 (IIFE) によるスコープ限定の初期化
	envMode := 2
	endpoint := func() string {
		switch envMode {
		case 1:
			return "http://localhost:8080"
		case 2:
			return "https://staging.api.internal:8443"
		default:
			return "https://api.production.com"
		}
	}()
	fmt.Printf("Configured endpoint via IIFE: %s\n", endpoint)
}

// ----------------------------------------------------------------------------
// 2. クロージャによる状態保持 (Stateful Closure)
// ----------------------------------------------------------------------------
func createCounter(start int) func() int {
	count := start // 外側のローカル変数 (エスケープ解析でヒープに退避される)

	// count への直接参照（ポインタ）を持つクロージャを返す
	return func() int {
		count++
		return count
	}
}

func demoClosuresAndState() {
	fmt.Println("\n=== 2. Closures & State Capture ===")

	counterA := createCounter(0)
	counterB := createCounter(100)

	fmt.Printf("Counter A: %d, %d, %d\n", counterA(), counterA(), counterA())
	fmt.Printf("Counter B: %d, %d\n", counterB(), counterB())
}

// ----------------------------------------------------------------------------
// 3. ループ変数キャプチャ (Go 1.22+ Per-Iteration Scope)
// ----------------------------------------------------------------------------
func demoLoopVariableCapture() {
	fmt.Println("\n=== 3. Loop Variable Capture in Goroutines (Go 1.22+) ===")

	values := []string{"Go", "Rust", "Python", "Java"}
	var wg sync.WaitGroup

	fmt.Println("Spawning goroutines inside loop:")
	for _, v := range values {
		wg.Add(1)
		// Go 1.22 以降: ループ変数 v はイテレーションごとに新しいインスタンスが生成される
		// そのため、引数渡しにしなくても各 goroutine が正しく固有の値を参照できる
		go func() {
			defer wg.Done()
			fmt.Printf("  Goroutine received: %s\n", v)
		}()
	}
	wg.Wait()
}

// ----------------------------------------------------------------------------
// 4. defer と無名関数 (名前付き戻り値の変更 & 計測)
// ----------------------------------------------------------------------------
func calculateSafely(val int) (result int, err error) {
	// defer に無名関数を渡すことで、関数終了直前に戻り値 result, err を書き換え可能
	defer func() {
		if val < 0 {
			result = 0
			err = fmt.Errorf("negative input not allowed: %d", val)
		}
	}()

	return val * 10, nil
}

func timeTracker(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("  [%s] Elapsed: %v\n", name, time.Since(start))
	}
}

func demoDeferWithAnonymousFunction() {
	fmt.Println("\n=== 4. defer with Anonymous Functions ===")

	// 1. 名前付き戻り値の書き換え検証
	res1, err1 := calculateSafely(5)
	fmt.Printf("Calculate(5)  -> Result: %d, Error: %v\n", res1, err1)

	res2, err2 := calculateSafely(-3)
	fmt.Printf("Calculate(-3) -> Result: %d, Error: %v\n", res2, err2)

	// 2. defer による関数実行時間の計測イディオム
	func() {
		defer timeTracker("SimulatedIOWait")()
		time.Sleep(10 * time.Millisecond)
	}()
}

// ----------------------------------------------------------------------------
// 5. Functional Options パターン (実務標準コンストラクタ)
// ----------------------------------------------------------------------------
type Server struct {
	Host    string
	Port    int
	Timeout time.Duration
	TLS     bool
}

type ServerOption func(*Server)

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.Port = port
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

func WithTLS(enabled bool) ServerOption {
	return func(s *Server) {
		s.TLS = enabled
	}
}

func NewServer(host string, opts ...ServerOption) *Server {
	// デフォルト値
	srv := &Server{
		Host:    host,
		Port:    8080,
		Timeout: 30 * time.Second,
		TLS:     false,
	}

	// 各オプションクロージャを適用
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

func demoFunctionalOptionsPattern() {
	fmt.Println("\n=== 5. Functional Options Pattern ===")

	s1 := NewServer("localhost")
	fmt.Printf("Default Server: Host=%s, Port=%d, TLS=%v, Timeout=%v\n",
		s1.Host, s1.Port, s1.TLS, s1.Timeout)

	s2 := NewServer("api.example.com",
		WithPort(443),
		WithTLS(true),
		WithTimeout(5*time.Second),
	)
	fmt.Printf("Custom Server:  Host=%s, Port=%d, TLS=%v, Timeout=%v\n",
		s2.Host, s2.Port, s2.TLS, s2.Timeout)
}

// ----------------------------------------------------------------------------
// 6. 高階関数 (Filter & Map)
// ----------------------------------------------------------------------------
func filter[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func mapSlice[T any, R any](items []T, mapper func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = mapper(item)
	}
	return result
}

func demoHigherOrderFunctions() {
	fmt.Println("\n=== 6. Higher-Order Functions (Generics Filter & Map) ===")

	names := []string{"Go", "Rust", "TypeScript", "Python", "C++", "Java"}

	// 4文字以上の言語をフィルタ
	longerThan3 := filter(names, func(s string) bool {
		return len(s) > 3
	})
	fmt.Printf("Filtered (>3 chars): %v\n", longerThan3)

	// 大文字に変換
	uppercased := mapSlice(longerThan3, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Printf("Uppercased:          %v\n", uppercased)
}
