# Modern Go Crash Course (For Rust, C#, Python, Java Developers)

Rust, C#, Python, Java などの言語を習得済みのエンジニアが、**最短でモダン Go（Go 1.22+ / 1.23+）をマスターするための実践リファレンス**です。

---

## 🚀 クイックスタート (実行方法)

```powershell
cd C:\Users\harun\programming\go\sample

# 全モジュールを一括実行
go run main.go

# または各モジュールを個別実行
go run ./01_basics/
go run ./02_errors/
go run ./03_interfaces/
go run ./04_lambdas_and_closures/
```

---

## 🗺️ 言語対比マッピング早見表 (Go vs Rust vs C# vs Java vs Python)

| 概念・機能 | Go | Rust | C# | Java | Python |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **型推論** | `x := 10` | `let x = 10;` | `var x = 10;` | `var x = 10;` | `x = 10` |
| **不変性** | `const` (基本型のみ) | デフォルト不変 (`mut`) | `readonly` / `record` | `final` / `record` | `frozen=True` |
| **エラーハンドリング** | 多値返却 `(T, error)` | `Result<T, E>` | `try-catch` | `try-catch` | `try-except` |
| **ポリモーフィズム** | **暗黙的インターフェース** | `trait` (明示的 `impl`)| `interface` (明示的) | `interface` (明示的) | ダックタイピング |
| **リソース解放** | `defer f.Close()` | `Drop` トレイト | `using` | `try-with-resources` | `with f:` |
| **無名関数 / ラムダ** | `func(x int) int` | `\|x\| x * 2` | `x => x * 2` | `x -> x * 2` | `lambda x: x * 2` |
| **並行処理モデル** | **goroutine + channel** | OSスレッド + tokio | `async / await` | Virtual Threads | `async / await` |
| **Null 安全性** | `nil` (ポインタ/インタフェース)| `Option<T>` | `Nullable<T>` / `?` | `Optional<T>` | `None` |

---

## ⚠️ 他言語経験者が最もハマる Go の「罠」と作法

### 1. `nil` インターフェースの罠（型情報と値の不一致）
- Go のインターフェースは `(型情報, 値)` のタプルです。
- 中身が `nil` のポインタを入れても、インターフェース変数自体は **`nil` になりません**。
  ```go
  var p *MyStruct = nil
  var i any = p
  fmt.Println(i == nil) // ❌ false! (型情報 *MyStruct が入っているため)
  ```

### 2. ループ変数キャプチャ問題（Go 1.22 で仕様変更）
- Go 1.21 以前: `for _, v := range slice { go func() { fmt.Println(v) }() }` はすべて最後の要素を参照するバグを生んでいました。
- **Go 1.22 以降**: ループのイテレーションごとに新しい `v` が生成されるように仕様改善され、そのまま正しく動くようになりました。

### 3. スライスの再割り当て（Capacity と共有）
- スライスは「ポインタ・長さ・容量」のヘッダ構造体です。
- `sub := arr[1:3]` で切り出したスライスは、元の配列のメモリを共有しています。`append` 時に容量（`cap`）を超えると新しいメモリへ再確保されて共有が切れます。

---

## 📁 提供サンプルコードの解説

| ディレクトリ | テーマ | 主な学習内容 |
| :--- | :--- | :--- |
| [`./01_basics`](file:///C:/Users/harun/programming/go/sample/01_basics/main.go) | **基本構文・型・関数** | 変数宣言3パターン、多値返却、`for` ループ、スライス vs 配列、`map`、ポインタと値渡し |
| [`./02_errors`](file:///C:/Users/harun/programming/go/sample/02_errors/main.go) | **エラーハンドリング & defer** | `error` インターフェース、`errors.Is` / `errors.As`、カスタムエラー型、`defer` の LIFO 順序 |
| [`./03_interfaces`](file:///C:/Users/harun/programming/go/sample/03_interfaces/main.go) | **インターフェース & ジェネリクス** | 暗黙的インターフェース（ダックタイピング）、型アサーション、Type Switch、Go 1.18+ Generics (`[T any]`) |
| [`./04_lambdas_and_closures`](file:///C:/Users/harun/programming/go/sample/04_lambdas_and_closures/main.go) | **無名関数 & クロージャ** | 即時実行無名関数 (IIFE)、状態保持クロージャ、Go 1.22+ ループ変数、`defer` での戻り値書き換え、Functional Options パターン、高階関数 |
| [`main.go`](file:///C:/Users/harun/programming/go/sample/main.go) | **統合エントリーポイント** | 全モジュールを順番に一括実行するメインランナー |

> 📖 **Go言語 無名関数 & クロージャの完全理解ガイド**:  
> なぜクロージャの変数はヒープに退避されるのか（エスケープ解析）、Go 1.22 の歴史的仕様変更、実務の Functional Options パターンまで完全網羅した解説は [**`LAMBDA.md`**](file:///C:/Users/harun/programming/go/sample/LAMBDA.md) を参照してください。
