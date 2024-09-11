package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRequest struct {
	Hash            string      `json:"hash"`
	ReceiverAddress string      `json:"receiver_address"` // Address of the Receiver
	SenderAddress   string      `json:"sender_address"`   // Address of the Receiver
	Amount          string      `json:"amount"`           // Address of the Receiver
	Message         pgtype.Text `json:"message"`
}

type TransactionStatusResponse struct {
	TransactionHash string `json:"transaction_hash"`
	Status          string `json:"status"`
}

func Transaction(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ganacheUrl = os.Getenv("GANACHE_URL")
		var req TransactionRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("Failed to parse request body:", err)
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		if req.Hash == "" {
			log.Println("Missing transaction hahs")
			http.Error(w, "Missing transaction hash", http.StatusBadRequest)
			return
		}
		// Connect to the local Ethereum client (Ganache)

		client, err := ethclient.Dial(ganacheUrl)
		if err != nil {
			log.Println("Failed to connect to the Ethereum client:", err)
			http.Error(w, "Failed to connect to the Ethereum client", http.StatusInternalServerError)
			return
		}
		// Check the transaction status
		tx, pending, err := client.TransactionByHash(context.Background(), common.HexToHash(req.Hash))
		if err != nil {
			log.Println("Failed to get transaction:", err)
			http.Error(w, "Failed to get transaction", http.StatusInternalServerError)
			return
		}

		// If the transaction is pending, notify that it's not confirmed yet
		status := "not confirmed"
		if tx != nil && !pending {
			status = "confirmed"
			query := `INSERT INTO "Transaction" ("sender", "reciever", "amount", "message", "created_at") VALUES ($1, $2, $3, $4, $5)`
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := pool.Exec(ctx, query, req.SenderAddress, req.ReceiverAddress, req.Amount, req.Message, time.Now().UTC())
			if err != nil {
				log.Println("Failed to save transaction to the database:", err)
				http.Error(w, "Failed to save transaction to the database", http.StatusInternalServerError)
				return
			}
		}

		response := TransactionStatusResponse{
			TransactionHash: req.Hash,
			Status:          status,
		}
		utils.Send(w, http.StatusOK, response)
		return
	}
}
