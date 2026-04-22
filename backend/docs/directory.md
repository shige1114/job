# ディレクトリ構成ルール（Kotlin / Ktor + Koin）

## 全体構造

```
myapp/
├── src/main/kotlin/com/example/myapp/
│   ├── Application.kt
│   │
│   ├── api/
│   │   ├── Routing.kt
│   │   └── middleware/
│   │       ├── AuthPlugin.kt
│   │       ├── LoggingPlugin.kt
│   │       └── ExceptionHandler.kt
│   │
│   ├── config/
│   │   ├── AppConfig.kt
│   │   └── DatabaseConfig.kt
│   │
│   ├── infra/
│   │   ├── db/
│   │   │   ├── DatabaseFactory.kt
│   │   │   ├── TransactionRunner.kt
│   │   │   └── UnitOfWorkImpl.kt
│   │   └── logger/
│   │       └── AppLogger.kt
│   │
│   ├── shared/
│   │   ├── apperror/
│   │   │   └── AppError.kt
│   │   ├── uow/
│   │   │   └── UnitOfWork.kt
│   │   └── domain/
│   │       └── valueobject/
│   │           ├── Money.kt
│   │           └── Pagination.kt
│   │
│   ├── di/
│   │   ├── AppModule.kt
│   │   └── InfraModule.kt
│   │
│   └── modules/
│       ├── order/
│       │   ├── domain/
│       │   │   ├── aggregate/
│       │   │   │   └── Order.kt
│       │   │   ├── entity/
│       │   │   │   └── OrderItem.kt
│       │   │   ├── valueobject/
│       │   │   │   ├── OrderId.kt
│       │   │   │   ├── OrderStatus.kt
│       │   │   │   └── CustomerId.kt
│       │   │   └── repository/
│       │   │       └── OrderRepository.kt
│       │   ├── usecase/
│       │   │   ├── CreateOrder.kt
│       │   │   ├── ConfirmOrder.kt
│       │   │   └── GetOrder.kt
│       │   ├── infrastructure/
│       │   │   └── persistence/
│       │   │       ├── OrderRepositoryImpl.kt
│       │   │       ├── OrderTable.kt
│       │   │       └── OrderMapper.kt
│       │   ├── presentation/
│       │   │   ├── dto/
│       │   │   │   ├── CreateOrderRequest.kt
│       │   │   │   ├── CreateOrderResponse.kt
│       │   │   │   └── GetOrderResponse.kt
│       │   │   └── route/
│       │   │       └── OrderRoute.kt
│       │   └── OrderModule.kt
│       │
│       ├── customer/
│       │   └── ...（同構造）
│       │
│       └── payment/
│           └── ...（同構造）
│
├── src/main/resources/
│   ├── application.conf
│   └── db/migration/
│
├── src/test/kotlin/com/example/myapp/
│   └── modules/
│       ├── order/
│       │   ├── domain/
│       │   │   └── OrderTest.kt
│       │   ├── usecase/
│       │   │   └── CreateOrderTest.kt
│       │   └── presentation/
│       │       └── OrderRouteTest.kt
│       └── ...
│
├── build.gradle.kts
└── settings.gradle.kts
```

## ルール

### エントリポイント
- `Application.kt` にKtorサーバーの起動とKoinの初期化を置く。

```kotlin
fun main() {
    embeddedServer(Netty, port = 8080) {
        install(Koin) {
            modules(appModule, infraModule, orderModule, customerModule, paymentModule)
        }
        configureRouting()
    }.start(wait = true)
}
```

### HTTP層（api/）
- `Routing.kt` で全モジュールのルートを集約する。
- `middleware/` にKtor Plugin（認証・ログ・例外ハンドリング）を置く。

### 設定（config/）
- `application.conf`（HOCON形式）に環境依存の値を外出しする。
- `AppConfig.kt` で型安全に読み込む。

### 外部依存（infra/）
- DB関連は `infra/db/`。
- ロガーは `infra/logger/`。
- フレームワーク依存はこの層に閉じ込める。

### 共有（shared/）
- 全モジュール横断のコードを置く。
- `shared/apperror/` — 共通エラー型（sealed class）。
- `shared/uow/` — UnitOfWorkインターフェース。
- `shared/domain/valueobject/` — 共通Value Object。

### DI（di/）
- `AppModule.kt` — 共通・横断的なDI定義。
- `InfraModule.kt` — インフラ層のDI定義。
- 各モジュールのDIは `modules/<module>/OrderModule.kt` に置く。

### モジュール（modules/）
- 業務機能は `modules/<module>/` に分割する。
- 各モジュールは以下の層を持つ。

### 依存方向
```
presentation → usecase → domain ← infrastructure
                           ↑
                         shared
```
- `domain/` は他の層に依存しない（純粋Kotlin。Ktor/Koinのimport禁止）。
- `usecase/` も純粋Kotlin（フレームワーク依存禁止）。
- `infrastructure/` は `domain/` のインターフェースを実装する（依存性逆転）。
- `presentation/` は `usecase/` を呼び出す。`domain/` を直接触らない。
- `shared/` は全モジュールから参照可能。`shared/` がモジュールに依存するのは禁止。

### Kotlin固有の型活用
- Value Object → `data class` + `val`（不変性をコンパイラが保証）。
- ドメインイベント・エラー型 → `sealed class`（網羅的パターンマッチ）。
- リポジトリ → `interface`（Kotlinネイティブ）。
- Nullable → `?` 型で明示（`null` の暗黙利用禁止）。

### 命名規則
- ファイル名・クラス名はPascalCase（`OrderStatus.kt`）。
- パッケージ名はlowercase（`valueobject`, `usecase`）。
- Ktorのルート定義は `Route` の拡張関数として書く。

---

# Presentation層ルール

## 原則
- presentation層はHTTPの関心事だけを扱う。ビジネスロジックを書かない。
- 責務: リクエスト受信 → DTOへの変換 → usecase呼び出し → レスポンス変換 → 返却。
- Ktor依存はこの層に閉じ込める。

## Route

`presentation/route/` にKtorの `Route` 拡張関数として定義する。

```kotlin
// modules/order/presentation/route/OrderRoute.kt

fun Route.orderRoutes() {
    val createOrder by inject<CreateOrder>()
    val getOrder by inject<GetOrder>()

    route("/orders") {
        post {
            val req = call.receive<CreateOrderRequest>()
            val command = req.toCommand()
            val result = createOrder.execute(command)
            when (result) {
                is CreateOrderResult.Success ->
                    call.respond(HttpStatusCode.Created, CreateOrderResponse.from(result))
                is CreateOrderResult.ValidationFailed ->
                    call.respond(HttpStatusCode.BadRequest, ErrorResponse(result.message))
            }
        }

        get("/{id}") {
            val id = call.parameters["id"] ?: throw BadRequestException("id is required")
            val result = getOrder.execute(OrderId(id))
            when (result) {
                is GetOrderResult.Found ->
                    call.respond(HttpStatusCode.OK, GetOrderResponse.from(result.order))
                is GetOrderResult.NotFound ->
                    call.respond(HttpStatusCode.NotFound, ErrorResponse("注文が見つかりません"))
            }
        }
    }
}
```

### ルール
- 1モジュール1ルートファイル。
- `Route` の拡張関数として定義する（クラスにしない）。
- usecaseは `by inject<>()` でKoinから取得する。
- `call.receive<>()` でリクエストDTOにデシリアライズ。
- usecaseの戻り値（sealed class）を `when` で網羅的に分岐し、HTTPステータスコードを決定する。
- presentation層でドメインオブジェクトを直接返さない。必ずResponse DTOに変換する。

## DTO

`presentation/dto/` にリクエスト・レスポンスのdata classを置く。

```kotlin
// Request DTO
@Serializable
data class CreateOrderRequest(
    val customerId: String,
    val items: List<OrderItemRequest>
) {
    fun toCommand() = CreateOrderCommand(
        customerId = CustomerId(customerId),
        items = items.map { it.toDomain() }
    )
}

@Serializable
data class OrderItemRequest(
    val productId: String,
    val quantity: Int
) {
    fun toDomain() = OrderItemInput(productId = productId, quantity = quantity)
}

// Response DTO
@Serializable
data class CreateOrderResponse(
    val orderId: String,
    val status: String
) {
    companion object {
        fun from(result: CreateOrderResult.Success) = CreateOrderResponse(
            orderId = result.orderId.value,
            status = result.status.name
        )
    }
}

// 共通エラーレスポンス
@Serializable
data class ErrorResponse(val message: String)
```

### ルール
- `@Serializable`（kotlinx.serialization）を付ける。
- Request DTOに `toCommand()` / `toQuery()` メソッドを持たせ、ドメインへの変換責務をDTOに置く。
- Response DTOに `companion object { fun from() }` を持たせ、ドメインからの変換責務をDTOに置く。
- DTOのフィールドはプリミティブ型（String, Int, List等）。Value Objectを直接持たない。

## Routingの集約

`api/Routing.kt` で全モジュールのルートを集約する。

```kotlin
// api/Routing.kt
fun Application.configureRouting() {
    routing {
        route("/api/v1") {
            orderRoutes()
            customerRoutes()
            paymentRoutes()
        }
    }
}
```

## 例外ハンドリング

`api/middleware/ExceptionHandler.kt` でグローバルに処理する。

```kotlin
fun Application.configureExceptionHandling() {
    install(StatusPages) {
        exception<BadRequestException> { call, cause ->
            call.respond(HttpStatusCode.BadRequest, ErrorResponse(cause.message ?: "Bad Request"))
        }
        exception<Throwable> { call, cause ->
            call.respond(HttpStatusCode.InternalServerError, ErrorResponse("Internal Server Error"))
        }
    }
}
```

### ルール
- ドメインエラーはusecaseの戻り値（sealed class）で表現する。例外を投げない。
- HTTPレベルのエラー（パラメータ不足等）のみ例外で処理する。
- presentation層で例外をcatchしてビジネスロジックの分岐をしない。

---

# Repository実装パターン

## 構造

```
domain/repository/
  └── OrderRepository.kt          ← interface（純粋Kotlin）

infrastructure/persistence/
  ├── OrderRepositoryImpl.kt      ← 実装（Exposed依存）
  ├── OrderTable.kt               ← テーブル定義
  └── OrderMapper.kt              ← Record ↔ Domain の変換
```

## domain層: インターフェース

```kotlin
// modules/order/domain/repository/OrderRepository.kt

interface OrderRepository {
    fun findById(id: OrderId): Order?
    fun findByCustomerId(customerId: CustomerId): List<Order>
    fun save(order: Order)
    fun delete(id: OrderId)
}
```

### ルール
- 純粋Kotlinのinterface。フレームワークのimport禁止。
- メソッド名はドメインの意図を表す（`findById`, `save`）。SQLの用語（`select`, `insert`）は使わない。
- 戻り値はドメインオブジェクト。DBのRecord型を返さない。

## infrastructure層: テーブル定義（Exposed）

```kotlin
// modules/order/infrastructure/persistence/OrderTable.kt

object OrderTable : Table("orders") {
    val id = varchar("id", 36)
    val customerId = varchar("customer_id", 36)
    val status = varchar("status", 20)
    val totalAmount = integer("total_amount")
    val currency = varchar("currency", 3)
    val createdAt = datetime("created_at")

    override val primaryKey = PrimaryKey(id)
}

object OrderItemTable : Table("order_items") {
    val id = varchar("id", 36)
    val orderId = varchar("order_id", 36) references OrderTable.id
    val productId = varchar("product_id", 36)
    val quantity = integer("quantity")
    val unitPrice = integer("unit_price")

    override val primaryKey = PrimaryKey(id)
}
```

## infrastructure層: マッパー

```kotlin
// modules/order/infrastructure/persistence/OrderMapper.kt

object OrderMapper {

    fun toDomain(row: ResultRow, items: List<ResultRow>): Order {
        return Order(
            id = OrderId(row[OrderTable.id]),
            customerId = CustomerId(row[OrderTable.customerId]),
            status = OrderStatus.valueOf(row[OrderTable.status]),
            items = items.map { toOrderItem(it) },
            totalAmount = Money(row[OrderTable.totalAmount], row[OrderTable.currency])
        )
    }

    private fun toOrderItem(row: ResultRow): OrderItem {
        return OrderItem(
            productId = row[OrderItemTable.productId],
            quantity = row[OrderItemTable.quantity],
            unitPrice = Money(row[OrderItemTable.unitPrice])
        )
    }

    fun toRecord(order: Order): Map<Column<*>, Any> {
        return mapOf(
            OrderTable.id to order.id.value,
            OrderTable.customerId to order.customerId.value,
            OrderTable.status to order.status.name,
            OrderTable.totalAmount to order.totalAmount.amount,
            OrderTable.currency to order.totalAmount.currency,
            OrderTable.createdAt to order.createdAt
        )
    }
}
```

### ルール
- `toDomain()`: DBのResultRow → ドメインオブジェクト。
- `toRecord()`: ドメインオブジェクト → DBカラム値のMap。
- マッパーはドメインロジックを持たない。純粋な変換のみ。

## infrastructure層: リポジトリ実装

```kotlin
// modules/order/infrastructure/persistence/OrderRepositoryImpl.kt

class OrderRepositoryImpl : OrderRepository {

    override fun findById(id: OrderId): Order? {
        val row = OrderTable.select { OrderTable.id eq id.value }.singleOrNull()
            ?: return null
        val items = OrderItemTable.select { OrderItemTable.orderId eq id.value }.toList()
        return OrderMapper.toDomain(row, items)
    }

    override fun findByCustomerId(customerId: CustomerId): List<Order> {
        return OrderTable
            .select { OrderTable.customerId eq customerId.value }
            .map { row ->
                val items = OrderItemTable
                    .select { OrderItemTable.orderId eq row[OrderTable.id] }
                    .toList()
                OrderMapper.toDomain(row, items)
            }
    }

    override fun save(order: Order) {
        val exists = OrderTable.select { OrderTable.id eq order.id.value }.count() > 0
        if (exists) {
            OrderTable.update({ OrderTable.id eq order.id.value }) {
                OrderMapper.toRecord(order).forEach { (col, value) ->
                    @Suppress("UNCHECKED_CAST")
                    it[col as Column<Any>] = value
                }
            }
        } else {
            OrderTable.insert {
                OrderMapper.toRecord(order).forEach { (col, value) ->
                    @Suppress("UNCHECKED_CAST")
                    it[col as Column<Any>] = value
                }
            }
        }
        saveItems(order)
    }

    private fun saveItems(order: Order) {
        OrderItemTable.deleteWhere { OrderItemTable.orderId eq order.id.value }
        order.items.forEach { item ->
            OrderItemTable.insert {
                it[id] = item.id
                it[orderId] = order.id.value
                it[productId] = item.productId
                it[quantity] = item.quantity
                it[unitPrice] = item.unitPrice.amount
            }
        }
    }

    override fun delete(id: OrderId) {
        OrderItemTable.deleteWhere { OrderItemTable.orderId eq id.value }
        OrderTable.deleteWhere { OrderTable.id eq id.value }
    }
}
```

### ルール
- Exposed（Kotlin製ORM）のDSLを使う。
- `save()` はupsert（存在すればupdate、なければinsert）。
- 集約単位で保存する（Order + OrderItems をまとめて）。
- トランザクション管理はリポジトリの外（usecase層 or UnitOfWork）で行う。
- リポジトリ実装はドメインロジックを持たない。永続化の関心事のみ。

## Koin DI定義

```kotlin
// modules/order/OrderModule.kt

val orderModule = module {
    single<OrderRepository> { OrderRepositoryImpl() }
    factory { CreateOrder(get()) }
    factory { ConfirmOrder(get()) }
    factory { GetOrder(get()) }
}
```

### ルール
- `single` — リポジトリ等のステートレスなもの。
- `factory` — usecase等、呼び出しごとに生成するもの。
- interfaceに対して実装をバインドする（`single<OrderRepository> { OrderRepositoryImpl() }`）。

## トランザクション管理

```kotlin
// infra/db/TransactionRunner.kt

object TransactionRunner {
    fun <T> run(block: () -> T): T {
        return transaction {
            block()
        }
    }
}

// usecase内での使用
class CreateOrder(private val repo: OrderRepository) {
    fun execute(command: CreateOrderCommand): CreateOrderResult {
        return TransactionRunner.run {
            // ここが1トランザクション
            val order = Order.create(command)
            repo.save(order)
            CreateOrderResult.Success(order.id, order.status)
        }
    }
}
```
