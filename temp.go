package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(string(hash)) 
}
// Вывод будет примерно: $2a$10$wTf7.KzDq2p4X5wM7yNn.eU7XzYq8IeJ1fL3pQ9gV4E0tYvC2G5O6