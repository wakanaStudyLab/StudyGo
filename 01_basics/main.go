package main

import "fmt"

// ============================================================
// Step 1: Go の基本文法
// 実行: go run ./01_basics/
// ============================================================

func main() {
	// --- 1. 変数宣言 ---
	// Goの変数宣言は3パターン
	var a int = 10       // 明示的な型指定
	var b = 20           // 型推論（C#のvar相当）
	c := 30              // 短縮宣言（関数内のみ使用可）
	fmt.Println(a, b, c) // 10 20 30

	// 定数
	const Pi = 3.14159
	fmt.Println(Pi)

	// --- 2. 基本的な型 ---
	var i int = 42
	var f float64 = 3.14
	var s string = "hello"
	var bl bool = true
	fmt.Println(i, f, s, bl)

	// --- 3. 関数呼び出し ---
	sum := add(3, 5)
	fmt.Println("3 + 5 =", sum)

	// Goは多値返却が基本
	q, r := divide(17, 5)
	fmt.Printf("17 / 5 = %d 余り %d\n", q, r)

	// --- 4. for ループ（Goにはforしかない）---
	// while相当
	n := 0
	for n < 3 {
		fmt.Println("n =", n)
		n++
	}

	// 通常のfor
	for i := 0; i < 3; i++ {
		fmt.Println("i =", i)
	}

	// --- 5. スライス（C++のvector / RustのVec相当）---
	nums := []int{10, 20, 30, 40, 50}
	fmt.Println("slice:", nums)
	fmt.Println("len:", len(nums), "cap:", cap(nums))

	// append で追加
	nums = append(nums, 60)
	fmt.Println("after append:", nums)

	// range でイテレート（C#のforeach相当）
	for idx, val := range nums {
		fmt.Printf("nums[%d] = %d\n", idx, val)
	}

	// --- 6. map（C++のunordered_map / C#のDictionary相当）---
	scores := map[string]int{
		"Alice": 90,
		"Bob":   75,
	}
	scores["Charlie"] = 85

	// 存在チェック
	val, ok := scores["Alice"]
	if ok {
		fmt.Println("Alice:", val)
	}

	// --- 7. 構造体 ---
	p := Person{Name: "Harun", Age: 25}
	fmt.Println(p)
	fmt.Println("Hello,", p.greet())

	// --- 8. ポインタ（C++と同じ概念、算術なし）---
	x := 42
	ptr := &x  // アドレス取得
	*ptr = 100 // デリファレンス
	fmt.Println("x =", x) // 100
}

// --- 関数定義 ---
// 引数の型は後ろに書く
func add(a int, b int) int {
	return a + b
}

// 多値返却（Goの慣用句）
func divide(a, b int) (int, int) {
	return a / b, a % b
}

// --- 構造体 ---
type Person struct {
	Name string
	Age  int
}

// メソッド（レシーバを先頭に書く）
// (p Person) は値レシーバ  → コピーで呼ばれる
// (p *Person) はポインタレシーバ → 変更が反映される
func (p Person) greet() string {
	return fmt.Sprintf("I'm %s, %d years old", p.Name, p.Age)
}
