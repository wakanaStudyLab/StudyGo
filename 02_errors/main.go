package main

import (
	"errors"
	"fmt"
	"os"
)

// ============================================================
// Step 2: エラーハンドリング + defer
// 実行: go run ./02_errors/
// ============================================================

func main() {
	// --- 1. 基本のエラーハンドリング ---
	// Goは例外(try/catch)がない。多値返却でエラーを返す。
	result, err := safeDivide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("10 / 2 =", result)

	// ゼロ除算を試す
	_, err = safeDivide(10, 0)
	if err != nil {
		fmt.Println("Error:", err) // Error: division by zero
	}

	// --- 2. カスタムエラー型 ---
	_, err = validateAge(-5)
	if err != nil {
		// 型アサーションでカスタムエラーの詳細を取得
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("Validation failed: field=%s, msg=%s\n", ve.Field, ve.Message)
		}
	}

	// --- 3. エラーのラップ（errors.Is / errors.As）---
	// C#のInnerExceptionに相当する概念
	ErrNotFound := errors.New("not found")
	wrapped := fmt.Errorf("user lookup failed: %w", ErrNotFound) // %w でラップ
	fmt.Println(errors.Is(wrapped, ErrNotFound))                 // true

	// --- 4. defer（C#のusing / C++のRAIIに相当）---
	// defer: 関数終了時に実行される（panicでも実行される）
	deferExample()

	// --- 5. ファイル操作でdeferを使う典型例 ---
	err = readFile("C:/Users/harun/programming/go/sample/go.mod")
	if err != nil {
		fmt.Println("File error:", err)
	}
}

// --- safeDivide ---
// エラーを返す関数の基本パターン
func safeDivide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// --- カスタムエラー型 ---
type ValidationError struct {
	Field   string
	Message string
}

// error interface は Error() string を実装するだけでOK
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

func validateAge(age int) (int, error) {
	if age < 0 {
		return 0, &ValidationError{
			Field:   "age",
			Message: "must be non-negative",
		}
	}
	return age, nil
}

// --- defer の動作確認 ---
func deferExample() {
	fmt.Println("=== defer example ===")
	defer fmt.Println("defer 1") // 最後に実行
	defer fmt.Println("defer 2") // 後からdeferした方が先に実行（LIFO）
	defer fmt.Println("defer 3")
	fmt.Println("normal execution")
	// 出力順: normal execution → defer 3 → defer 2 → defer 1
}

// --- ファイル読み込み（deferの典型的な使い方）---
func readFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("readFile: %w", err) // エラーをラップして返す
	}
	defer f.Close() // ← ここでcloseを予約。return後も必ず実行される

	// ファイル操作...
	info, _ := f.Stat()
	fmt.Printf("File: %s, Size: %d bytes\n", info.Name(), info.Size())
	return nil
}
