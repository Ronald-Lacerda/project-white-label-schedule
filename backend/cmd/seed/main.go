// cmd/seed/main.go — Seed de dados para desenvolvimento.
// Uso: make seed
// Cria um estabelecimento e um usuário gestor se não existirem.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"schedule/internal/shared"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	dsn = strings.TrimPrefix(dsn, "mysql://")

	db, err := sqlx.Connect("mysql", dsn+"?parseTime=true&loc=UTC&charset=utf8mb4")
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Verifica se já existe
	var count int
	db.GetContext(ctx, &count, `SELECT COUNT(*) FROM users WHERE email = ?`, "admin@exemplo.com")
	if count > 0 {
		fmt.Println("Seed já aplicado — usuário admin@exemplo.com já existe.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt: %v", err)
	}

	estID := shared.NewID()
	userID := shared.NewID()
	now := time.Now().UTC()

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO establishments (id, name, slug, timezone, active, google_calendar_connected, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		estID, "Barbearia Silva", "barbearia-silva", "America/Sao_Paulo", true, false, now, now,
	)
	if err != nil {
		log.Fatalf("insert establishment: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO whitelabel_configs (establishment_id, primary_color)
		VALUES (?, ?)`,
		estID, "#1a1a2e",
	)
	if err != nil {
		log.Fatalf("insert whitelabel: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, establishment_id, name, email, password_hash, role, active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, estID, "Admin", "admin@exemplo.com", string(hash), "owner", true, now,
	)
	if err != nil {
		log.Fatalf("insert user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("commit: %v", err)
	}

	fmt.Println("✓ Seed aplicado com sucesso!")
	fmt.Printf("  Estabelecimento: Barbearia Silva (slug: barbearia-silva)\n")
	fmt.Printf("  Usuário: admin@exemplo.com\n")
	fmt.Printf("  Senha:   123456\n")
	fmt.Printf("  ID:      %s\n", estID)
}
