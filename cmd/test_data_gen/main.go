package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/schollz/progressbar"
	"github.com/tjarratt/babble"
	"golang.org/x/sys/unix"

	"gitlab.ozon.dev/ninashvl/homework-1/config"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	expstore "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage/psql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	testUsersCount = 100
	randomWindow   = 100
)

// Тулза для генерации тестовых данных в бд

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		unix.SIGTERM, unix.SIGKILL, unix.SIGINT)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	db, err := sql.Open("pgx", cfg.PsqlDSN())
	if err != nil {
		log.Fatal(err)
	}
	bar := progressbar.New(testUsersCount)
	s := expstore.New(db, zerolog.Logger{})
	b := babble.NewBabbler()
	for i := 0; i < testUsersCount; i++ {
		_ = bar.Add(1)
		select {
		case <-ctx.Done():
			log.Println("Gracefully shutdown")
			return
		default:

			userID := rand.Int63()

			// Выставление пользовательских трат
			curr := models.RubCurrency
			switch i % 4 {
			case 1:
				curr = models.EurCurrency
			case 2:
				curr = models.UsdCurrency
			case 3:
				curr = models.CnyCurrency
			}
			err := s.SetCurrency(ctx, userID, curr)
			if err != nil {
				log.Println("[ERROR] SetCurrency error: ", err)
				continue
			}

			// Генерация пользовательских трать
			countOfExps := rand.Intn(randomWindow)
			for j := 0; j < countOfExps; j++ {
				cat := b.Babble()
				amount := rand.Float64() * math.Pow(10, float64(rand.Intn(4)))
				exp := &models.Expense{
					Amount:   amount,
					Category: cat,
					Date:     time.Now(),
				}
				err = s.Add(ctx, userID, exp)
				if err != nil {
					log.Println("[ERROR] Add expense error: ", err)
					continue
				}
			}
		}
	}
	fmt.Println()
	log.Println("[INFO] test data successfully generated")
}
