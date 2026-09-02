# Go言語 無名関数 & クロージャ 完全理解ガイド (Go Anonymous Functions & Closures Deep Dive)

Go 言語における **第一級関数（First-class functions）**、**無名関数（Anonymous Functions）**、そして **クロージャ（Closures）** の完全解説書です。

「Go の関数型と型定義」「クロージャによる外側変数のキャプチャ機構」「**Go 1.22 で解決された歴史的超有名バグ（ループ変数キャプチャ問題）**」「エスケープ解析によるヒープ退避」「`defer` や Functional Options などの実務パターン」まで徹底的に解き明かします。

---

## 📑 目次

1. [Go における「関数」の第一級性（First-class Functions）](#1-go-における関数の第一級性first-class-functions)
2. [無名関数の基本構文と即時実行 (IIFE)](#2-無名関数の基本構文と即時実行-iife)
3. [クロージャの正体：変数キャプチャと参照モデル](#3-クロージャの正体変数キャプチャと参照モデル)
4. [歴史的超重要トピック：Go 1.22 ループ変数キャプチャの仕様変更](#4-歴史的超重要トピックgo-122-ループ変数キャプチャの仕様変更)
5. [アンダー・ザ・フード：エスケープ解析とヒープ退避](#5-アンダーザフードエスケープ解析とヒープ退避)
6. [実務実践パターン集](#6-実務実践パターン集)
   - 6-1. `defer` と無名関数（名前付き戻り値の変更・時間計測）
   - 6-2. Functional Options パターン（Go 流エレガントなコンストラクタ）
   - 6-3. HTTP ミドルウェアパターン
   - 6-4. goroutine とクロージャ
7. [他言語エンジニア向け比較表 (Go vs Java vs C++ vs Rust vs Python)](#7-他言語エンジニア向け比較表-go-vs-java-vs-c-vs-rust-vs-python)
8. [理解度チェッククイズ](#8-理解度チェッククイズ)

---

## 1. Go における「関数」の第一級性（First-class Functions）

Go 言語では、関数は数値や文字列と同じ **「第一級の値」** です。
- 変数に関数を代入できる
- 関数の引数に関数を渡せる（高階関数）
- 関数の戻り値として関数を返せる

### 関数型（Function Type）と `type` エイリアス
Go では関数のシグネチャそのものが「型」になります。

```go
// int を2つ受け取って int を返す関数型
type BinaryOp func(a, b int) int

func calculate(a, b int, op BinaryOp) int {
    return op(a, b)
}

func main() {
    add := func(a, b int) int { return a + b }
    fmt.Println(calculate(10, 20, add)) // 30
}
```

---

## 2. 無名関数の基本構文と即時実行 (IIFE)

### 2-1. 基本構文
関数名をつけずに `func(引数) 戻り値 { 処理 }` と記述します。

```go
// 変数に代入して後から呼び出す
multiply := func(x, y int) int {
    return x * y
}
result := multiply(4, 5) // 20
```

### 2-2. 即時実行無名関数 (IIFE: Immediately Invoked Function Expression)
定義したその場で末尾に `(引数)` を付けて即座に実行するパターンです。

```go
func() {
    fmt.Println("即座に実行されます")
}()

// 初期化スコープの限定や一時的な変数の局所化に有用
configValue := func() string {
    if isProd {
        return "https://api.production.com"
    }
    return "http://localhost:8080"
}()
```

---

## 3. クロージャの正体：変数キャプチャと参照モデル

### 3-1. クロージャ（Closure）とは？
無名関数が、**自分自身のスコープの外側にあるローカル変数を参照（キャプチャ）している状態** の関数を「クロージャ」と呼びます。

```go
func createCounter() func() int {
    count := 0 // 外側のローカル変数
    
    // この無名関数は count 変数をキャプチャして閉包（Closure）している
    return func() int {
        count++
        return count
    }
}

func main() {
    c1 := createCounter()
    fmt.Println(c1()) // 1
    fmt.Println(c1()) // 2
    
    c2 := createCounter() // 新しい独立した環境を持つクロージャ
    fmt.Println(c2()) // 1
}
```

### 3-2. Go のキャプチャは「参照（ポインタ）」である
Java のラムダ式は値を「コピー（実質的 final）」しますが、**Go のクロージャは外側の変数を「直接の参照（メモリ番地）」としてキャプチャ**します。

そのため、クロージャの内部から外側の変数を直接書き換えることができますし、逆に外側の変数が書き換わるとクロージャ内から見える値も変わります。

---

## 4. 歴史的超重要トピック：Go 1.22 ループ変数キャプチャの仕様変更

Go 言語の歴史において、最も多くのバグを生み出し、Go 開発者全員を悩ませてきたのが **「ループ変数の参照キャプチャ（Loop Variable Capture Trap）」** です。

### 4-1. Go 1.21 以前の地獄（すべて最後の要素が出力されるバグ）
Go 1.21 以前では、`for` ループの変数（`v` や `i`）は **「ループ全体で1つのメモリアドレス」** を使い回す仕様でした。

```go
// ❌ Go 1.21 以前で頻発した致命的バグ
values := []string{"a", "b", "c"}
var funcs []func()

for _, v := range values {
    // 全てのクロージャが同一のメモリ番地 &v を参照してしまう！
    funcs = append(funcs, func() {
        fmt.Println(v)
    })
}

for _, f := range funcs {
    f()
}
// 期待: "a", "b", "c"
// 実際 (Go 1.21以前): "c", "c", "c" (全部最後の要素！)
```

### 4-2. Go 1.21 以前の回避策（イディオム）
かつてはループ内で変数を再シャドウイング（`v := v`）するか、引数として渡す必要がありました。
```go
for _, v := range values {
    v := v // ループ内のローカル変数として再コピー
    funcs = append(funcs, func() { fmt.Println(v) })
}
```

### 4-3. Go 1.22（2024年2月リリース）の神仕様変更 ★現代の標準★
Go 1.22 以降、**`for` ループのイテレーションごとに新しい変数が自動的に生成される（Loop variable per-iteration）** ように言語仕様が修正されました！

```go
// ⭕ Go 1.22 以降: 何も小細工をしなくても正しく動く！
for _, v := range values {
    funcs = append(funcs, func() {
        fmt.Println(v) // 各ループで新しい v がキャプチャされる
    })
}
// 出力: "a", "b", "c"
```
> **⚠️ 現場知識**: 既存の古いコードベース（Go 1.21 未満）を保守・移行する際や、コーディング面接では今でも頻出の知識です。

---

## 5. アンダー・ザ・フード：エスケープ解析とヒープ退避

### 5-1. 通常のスタック変数の寿命
通常、関数のローカル変数はスタックに確保され、関数をリターンすると消滅します。
しかし、クロージャがその変数をキャプチャして関数の外に返された場合、どうなるのでしょうか？

```go
func makeGreeter() func() string {
    msg := "Hello from stack?"
    return func() string {
        return msg // msg は関数 makeGreeter 終了後も使われる！
    }
}
```

### 5-2. エスケープ解析 (Escape Analysis)
Go コンパイラはビルド時に **エスケープ解析** を行います。
1. ローカル変数が関数の外に生き残る（クロージャにキャプチャされて返される）ことを検知。
2. その変数をスタックではなく、**ガベージコレクタ（GC）が管理する「ヒープ領域」に自動的に割り当て（エスケープ）** します。
3. クロージャ自身は、**「関数ポインタ」＋「ヒープ上のキャプチャ環境構造体へのポインタ」のペア（Fat Pointer）** としてメモリ上に表現されます。

> **検証コマンド**:  
> `go build -gcflags="-m"` を実行すると、コンパイラが `moved to heap: msg` とエスケープ解析したログを確認できます。

---

## 6. 実務実践パターン集

### 6-1. `defer` と無名関数
`defer` に無名関数を渡すことで、関数の終了時にクリーンアップやメトリクス計測、**名前付き戻り値の事後書き換え** が可能です。

```go
// 1. 名前付き戻り値のエラー書き換え・ロギング
func processItem(id int) (err error) {
    defer func() {
        if err != nil {
            log.Printf("Failed processing id=%d: %v", id, err)
        }
    }()
    
    // 何らかの処理...
    return errors.New("database timeout")
}

// 2. 実行時間の計測 (Timer)
func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    fmt.Printf("%s took %s\n", name, elapsed)
}

func heavyTask() {
    defer timeTrack(time.Now(), "heavyTask")
    time.Sleep(100 * time.Millisecond)
}
```

### 6-2. Functional Options パターン（Go 流コンストラクタ）
Go には関数のオーバーロードやオプショナル引数がありません。
クロージャを使った **Functional Options パターン** は、拡張可能で美しい設定オブジェクトを作成する実務のデファクトスタンダードです。

```go
type Server struct {
    host string
    port int
    tls  bool
}

// Option は Server を変更するクロージャ
type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTLS(tls bool) Option {
    return func(s *Server) {
        s.tls = tls
    }
}

// コンストラクタ
func NewServer(host string, opts ...Option) *Server {
    srv := &Server{
        host: host,
        port: 8080, // デフォルト値
        tls:  false,
    }
    for _, opt := range opts {
        opt(srv) // 各クロージャを実行して設定を適用
    }
    return srv
}

// 呼び出し側（極めて可読性が高い！）
s1 := NewServer("localhost")
s2 := NewServer("example.com", WithPort(443), WithTLS(true))
```

### 6-3. HTTP ミドルウェアパターン
Go の標準ライブラリ `net/http` では、無名関数とクロージャを使ってミドルウェアを連鎖させます。

```go
type Middleware func(http.HandlerFunc) http.HandlerFunc

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[%s] %s", r.Method, r.URL.Path)
        next(w, r) // 次のハンドラを呼び出す
    }
}
```

### 6-4. goroutine とクロージャ
並行ワーカーにクロージャを渡すパターンです。

```go
var wg sync.WaitGroup
for i := 1; i <= 3; i++ {
    wg.Add(1)
    go func(workerID int) {
        defer wg.Done()
        fmt.Printf("Worker %d finished\n", workerID)
    }(i) // 引数として渡すことで確実に分離
}
wg.Wait()
```

---

## 7. 他言語エンジニア向け比較表 (Go vs Java vs C++ vs Rust vs Python)

| 項目 | Go | Java (21+) | Modern C++ (20+) | Rust | Python |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **記法** | `func(x int) int {}` | `(x) -> x*2` | `[](auto x){}` | `\|x\| x*2` | `lambda x: x*2` |
| **複数行・文** | ⭕ **完全サポート** | ⭕ サポート | ⭕ サポート | ⭕ サポート | ❌ **式のみ** |
| **外部変数の書き換え**| ⭕ **直接可能 (参照)** | ❌ 不可 (final限定) | ⭕ `[&]` や `mutable` | ⭕ `mut` 借用 / move | ⚠️ `nonlocal` が必要 |
| **メモリ配置** | 自動エスケープ (ヒープ) | 初回動的生成 (ヒープ) | スタック (無名構造体) | スタック (無名構造体) | ヒープ (`cell`) |
| **ゼロオーバーヘッド**| ❌ (GC・ヒープ退避) | 僅少 (Bootstrap後) | ⭕ **完全インライン** | ⭕ **完全インライン** | ❌ (関数呼出コスト) |

---

## 8. 理解度チェッククイズ

### Q1. Go 1.22 より前のバージョンで、以下のコードを実行すると何が起きる可能性がありましたか？
```go
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)
    }()
}
```
<details>
<summary>▶ 解答と解説</summary>

**解答**: すべての goroutine が同じ変数 `i` のメモリ番地を参照するため、goroutine が実行される頃にはループが終了しており、`3 3 3` と出力される競合状態（Race Condition）が発生していました。Go 1.22 からはイテレーションごとに新しい `i` が作られるため `0, 1, 2` が出力されます。
</details>

### Q2. Go のクロージャが関数の外側のローカル変数を参照している場合、そのローカル変数はメモリのどこに配置されますか？
- A. 必ずスタック領域
- B. コンパイラのエスケープ解析によりヒープ領域に自動退避される
- C. テキスト（コード）セグメント
- D. グローバル変数領域

<details>
<summary>▶ 解答と解説</summary>

**正解: B**
Go コンパイラはエスケープ解析を行い、関数のスコープを超えて生存する必要がある変数を自動的にヒープ領域に割り当てます。
</details>

---

## まとめ

1. **「Go の関数は第一級、クロージャは参照キャプチャ」**: 外側の変数を直接書き換えられる柔軟性を持つ。
2. **「Go 1.22 のループ仕様変更を理解する」**: 歴史的な `v := v` の背景を知ることで古いコードも読める。
3. **「実務は `defer` と Functional Options」**: クロージャを活用した Go らしい洗練されたアーキテクチャを身につける。
