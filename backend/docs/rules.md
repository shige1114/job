# Kotlinコーディング規約・静的解析ルール

## 公式スタイルガイド

JetBrains公式コーディング規約に準拠する。
https://kotlinlang.org/docs/coding-conventions.html

## ツール構成

| ツール | 用途 | 設定 |
|---|---|---|
| ktfmt | フォーマット | Googleスタイル準拠。極めて厳格なフォーマッター |
| detekt | 静的解析 | コード複雑度・コードスメル・潜在バグ検出 |

```kotlin
// build.gradle.kts
plugins {
    id("com.ncorti.ktfmt.gradle") version "0.16.0"
    id("io.gitlab.arturbosch.detekt") version "1.23.6"
}

ktfmt {
    googleStyle()
}
```

---

## Kotlin設計指針

### 1. `data class` を積極的に使う

Value Object、DTO、Command、Resultは全て `data class` で定義する。
`equals`, `hashCode`, `copy`, `toString` が自動生成される。

```kotlin
data class OrderId(val value: String)
data class Money(val amount: Int, val currency: String = "JPY")
data class CreateOrderCommand(val customerId: CustomerId, val items: List<OrderItemInput>)
```

### 2. `sealed class` でドメインの分岐を表現する

ドメインエラー・ユースケースの結果・状態遷移は `sealed class` で表現する。
`when` で網羅チェックされ、ケース漏れがコンパイルエラーになる。

```kotlin
sealed class WithdrawResult {
    data class Success(val balance: Money) : WithdrawResult()
    data class Rejected(val reason: String) : WithdrawResult()
}

// 呼び出し側 — 新しいケース追加時、ここに書かないとコンパイル通らない
when (result) {
    is WithdrawResult.Success -> ...
    is WithdrawResult.Rejected -> ...
}
```

### 3. 例外よりsealed classで結果を返す

例外はHTTPレベルのエラー（パラメータ不足等）にのみ使用する。
ドメインエラーは例外ではなくsealed classで表現する。

```kotlin
// ❌ 例外で制御フロー
fun withdraw(amount: Money): Money {
    if (balance < amount) throw InsufficientBalanceException()
    ...
}

// ✅ sealed classで明示
fun withdraw(amount: Money): WithdrawResult {
    if (balance.isLessThan(amount)) return WithdrawResult.Rejected("残高不足")
    ...
}
```

### 4. `null` は `?` で明示する。`!!` は禁止

見つからない可能性がある場合は戻り値を `?` 型にする。
`!!`（非nullアサーション）は使用禁止。detektで検出する。

```kotlin
// ✅ null安全
fun findById(id: OrderId): Order?

repo.findById(id)?.let { order ->
    order.confirm()
}

// ❌ NullPointerExceptionの温床
order!!.confirm()
```

### 5. 拡張関数でドメインの表現力を上げる

ドメインの語彙を拡張関数で表現する。ただし乱用しない。

```kotlin
fun Money.isAffordable(price: Money): Boolean = this.amount >= price.amount
```

### 6. スコープ関数の使い分け

| 関数 | 用途 | 例 |
|---|---|---|
| `let` | nullチェック後の処理 | `order?.let { it.confirm() }` |
| `apply` | オブジェクト初期化 | ビルダー的な初期化 |
| `also` | 副作用（ログ等） | `order.also { logger.info("created: ${it.id}") }` |
| `run` / `with` | 複数プロパティへのアクセス | — |

**禁止**: スコープ関数のネスト（2段以上）。読みにくくなる。

```kotlin
// ❌ ネスト禁止
repo.findById(id)?.let { it.also { it.run { ... } } }
```

### 7. 不変性を優先する

- `val` をデフォルトで使う。`var` は集約のミュータブルな内部状態にのみ使用する。
- コレクションは `List`（不変）をデフォルトにする。`MutableList` は内部実装にのみ使用する。

```kotlin
// ✅ 外部には不変リストを公開
class Order(
    val id: OrderId,
    private val _items: MutableList<OrderItem> = mutableListOf()
) {
    val items: List<OrderItem> get() = _items.toList()
}
```

### 8. 名前付き引数を活用する

引数が2つ以上の場合、名前付き引数で可読性を上げる。

```kotlin
// ✅ 意図が明確
val order = Order.create(
    customerId = CustomerId("C-001"),
    items = listOf(item1, item2),
    totalAmount = Money(amount = 3000)
)

// ❌ 何が何だかわからない
val order = Order.create(CustomerId("C-001"), listOf(item1, item2), Money(3000))
```

---

## detekt設定

```yaml
# detekt.yml

complexity:
  LongMethod:
    threshold: 20
  CyclomaticComplexity:
    threshold: 10
  LongParameterList:
    threshold: 5
  TooManyFunctions:
    threshold: 15

style:
  ForbiddenComment:
    values: ['TODO', 'FIXME']
  MaxLineLength:
    maxLineLength: 120
  WildcardImport:
    active: true

potential-bugs:
  UnsafeCallOnNullableType:
    active: true  # !! の検出

naming:
  FunctionNaming:
    functionPattern: '[a-z][a-zA-Z0-9]*'
  ClassNaming:
    classPattern: '[A-Z][a-zA-Z0-9]*'
  PackageNaming:
    packagePattern: '[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*'
```

---

## ドメイン層の汚染防止ルール

domain層（`modules/*/domain/`）のファイルに以下のimportがあったら違反:

- `import io.ktor.*`
- `import org.koin.*`
- `import org.jetbrains.exposed.*`
- `import kotlinx.serialization.*`

CIでgrepチェックを推奨:

```bash
# domain層にフレームワーク依存がないことを検証
grep -r "import io.ktor\|import org.koin\|import org.jetbrains.exposed\|import kotlinx.serialization" \
  src/main/kotlin/**/modules/*/domain/ && exit 1 || exit 0
```
