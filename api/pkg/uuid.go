package pkg

import "github.com/google/uuid"

// generateIDはIDを生成するための仮の関数です。
func GenerateID() string {
	return uuid.NewString()
}
