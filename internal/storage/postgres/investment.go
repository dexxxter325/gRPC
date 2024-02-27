package postgres

import (
	"GRPC/internal/domain/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvestmentPostgres struct {
	db *pgxpool.Pool
}

func NewInvestmentPostgres(db *pgxpool.Pool) *InvestmentPostgres {
	return &InvestmentPostgres{db: db}
}

func (p *InvestmentPostgres) Create(ctx context.Context, amount int64, currency string) (investmentId int64, err error) {
	query := `insert into investments(amount,currency) values ($1,$2) returning id`
	row := p.db.QueryRow(ctx, query, amount, currency)
	if err = row.Scan(&investmentId); err != nil {
		return 0, fmt.Errorf("scan failed in create:%s", err)
	}
	return investmentId, err
}

func (p *InvestmentPostgres) Get(ctx context.Context) (investment []models.Investment, err error) {
	query := `select * from investments ORDER BY id ASC`
	row, err := p.db.Query(ctx, query)
	if err != nil {
		return []models.Investment{}, fmt.Errorf("query failed in get:%s", err)
	}
	defer row.Close()
	for row.Next() { //читает каждую строку из результата sql query ,чтобы затем добавить в срез
		var investments models.Investment
		if err = row.Scan(&investments.ID, &investments.Amount, &investments.Currency); err != nil {
			return []models.Investment{}, fmt.Errorf("scan failed in get:%s", err)
		}
		investment = append(investment, investments) //с каждой успешно прочитанной строкой,добавляем эл-т в срез
	}
	return investment, nil
}

func (p *InvestmentPostgres) Delete(ctx context.Context, investmentId int64) error {
	query := `delete from investments where id=$1`
	row, err := p.db.Exec(ctx, query, investmentId)
	if err != nil {
		return fmt.Errorf("exec failed in delete:%s", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("already deleted or doesn't exist:%s", err)
	}
	return nil
}
