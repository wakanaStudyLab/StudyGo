package main

import (
	"fmt"
	"math"
)

// ============================================================
// Step 3: interface
// 実行: go run ./03_interfaces/
// ============================================================

// --- interface 定義 ---
// メソッドのシグネチャだけを定義
type Shape interface {
	Area() float64
	Perimeter() float64
}

// --- 実装1: Circle ---
type Circle struct {
	Radius float64
}

// implements宣言不要。メソッドを揃えるだけでShapeを満たす
func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// --- 実装2: Rectangle ---
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// --- 実装3: Triangle ---
type Triangle struct {
	A, B, C float64 // 3辺の長さ
}

func (t Triangle) Area() float64 {
	// ヘロンの公式
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}
func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }

// --- interface を受け取る関数 ---
// 具体的な型ではなく Shape として受け取る
func printShapeInfo(s Shape) {
	// 型スイッチで元の型を判定
	switch v := s.(type) {
	case Circle:
		fmt.Printf("Circle(r=%.1f)", v.Radius)
	case Rectangle:
		fmt.Printf("Rectangle(%.1fx%.1f)", v.Width, v.Height)
	case Triangle:
		fmt.Printf("Triangle(%.1f, %.1f, %.1f)", v.A, v.B, v.C)
	default:
		fmt.Printf("Unknown shape: %T", v)
	}
	fmt.Printf(" → Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func main() {
	// --- 1. 各図形を作成 ---
	c := Circle{Radius: 5}
	r := Rectangle{Width: 4, Height: 6}
	t := Triangle{A: 3, B: 4, C: 5}

	// --- 2. interface型のスライスにまとめられる ---
	shapes := []Shape{c, r, t}

	fmt.Println("=== Shape Info ===")
	for _, s := range shapes {
		printShapeInfo(s)
	}

	fmt.Printf("\nTotal Area: %.2f\n", totalArea(shapes))

	// --- 3. 型アサーション ---
	var s Shape = Circle{Radius: 3}

	// 安全な型アサーション（okで確認）
	if circle, ok := s.(Circle); ok {
		fmt.Printf("\nIt's a circle! Radius = %.1f\n", circle.Radius)
	}

	// --- 4. interface の合成 ---
	// Stringer: fmt.Println が自動的に呼ぶ interface
	p := Point{X: 3, Y: 4}
	fmt.Println("\nPoint:", p)        // String() が自動で呼ばれる
	fmt.Println("Distance:", p.Distance())
}

// --- Stringer interface の実装 ---
// fmt.Stringer: String() string を実装すると fmt.Println で使われる
type Point struct {
	X, Y float64
}

func (p Point) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", p.X, p.Y)
}

func (p Point) Distance() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}
