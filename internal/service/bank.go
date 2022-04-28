package service

import "context"

type BankRepo interface {
	AddBank(context.Context, *Bank) error
	GetBankList(context.Context) ([]*Bank, error)
}

type BankService struct {
	repo BankRepo
}

type Bank struct {
	Id   string
	Name string
}

func NewBankService(repo BankRepo) *BankService {
	return &BankService{
		repo: repo,
	}
}

func (bs *BankService) AddBank(ctx context.Context, bank *Bank) error {
	return bs.repo.AddBank(ctx, bank)
}

func (bs *BankService) GetBankList(ctx context.Context) ([]*Bank, error) {
	return bs.repo.GetBankList(ctx)
}
