package main

import (
	"log"

	"github.com/shige1114/job/backend/internal/account/infrastructure/persistence"
)

func main() {
	// 1. データベースの初期化
	db, err := persistence.NewDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// 2. Wire で生成された関数を呼び出し、依存関係が解決されたルーターを取得
	// 注意: wire コマンドを実行して wire_gen.go が生成されている必要があります
	r := InitializeRouter(db)

	// 3. サーバーの起動
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
