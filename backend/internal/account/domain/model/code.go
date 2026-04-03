package model

// Code はメール確認用の認証コードです
type Code struct {
	value string
}

func NewCode(v string) Code {
	return Code{value: v}
}

func (c Code) String() string {
	return c.value
}
