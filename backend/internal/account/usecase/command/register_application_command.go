// Package command 引数をもとにコマンド作成
package command

type RegisterApplicationCommand struct {
	Email string
}

type RegisterApplicationResult struct {
	ID string
}
