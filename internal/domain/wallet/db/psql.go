package db

import (
	"context"
	"log/slog"
	wallet2 "wallet/internal/domain/wallet"
)

type Storage struct {
	Client wallet2.PsqlClient
	logger *slog.Logger
}

func NewRepository(client wallet2.PsqlClient, logger *slog.Logger) wallet2.Storage {
	return &Storage{client, logger}
}

func (s *Storage) CreateWallet(ctx context.Context, userID string) error {
	const op = "wallet.db.CreateWallet"
	log := s.logger.With(slog.String("op", op))

	q := `INSERT INTO wallet(
                    user_id,
                    balance_eur,
                    balance_usd,
                    balance_rub) 
		  VALUES ($1, $2, $3, $4)`
	_, err := s.Client.Exec(ctx, q, userID, 0, 0, 0)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Storage) GetWalletByUserID(ctx context.Context, UserID string) (wallet2.Wallet, error) {
	const op = "wallet.db.GetWalletByUserID"
	log := s.logger.With(slog.String("op", op))

	q := `SELECT id, 
       			 balance_eur, 
       			 balance_usd, 
       			 balance_rub 
		  FROM wallet WHERE user_id=$1`

	var w wallet2.Wallet
	err := s.Client.QueryRow(ctx, q, UserID).Scan(&w.UUID, &w.BalanceEUR, &w.BalanceUSD, &w.BalanceRUB)
	if err != nil {
		log.Error(err.Error())
		return w, err
	}
	w.UserUUID = UserID
	return w, nil
}

func (s *Storage) UpdateWallet(ctx context.Context, userID string, wallet wallet2.Wallet) error {
	const op = "wallet.db.UpdateWallet"
	log := s.logger.With(slog.String("op", op))

	q := `UPDATE wallet SET 
                  balance_eur = $1, 
                  balance_usd = $2, 
                  balance_rub = $3
		  WHERE user_id=$4`

	_, err := s.Client.Exec(ctx, q, wallet.BalanceEUR, wallet.BalanceUSD, wallet.BalanceRUB, userID)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Storage) DeleteWallet(ctx context.Context, userID string) error {
	const op = "wallet.db.DeleteWallet"
	log := s.logger.With(slog.String("op", op))

	q := `DELETE FROM wallet WHERE user_id=$1`
	_, err := s.Client.Exec(ctx, q, userID)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil

}
