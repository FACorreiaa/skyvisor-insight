package airline

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AirlineRepository struct {
	pgpool *pgxpool.Pool
}

func NewAirports(
	pgpool *pgxpool.Pool,

) *AirlineRepository {
	return &AirlineRepository{
		pgpool: pgpool,
	}
}

func (r *AirlineRepository) GetTax(ctx context.Context, page, pageSize int) ([]models.Tax, error) {
	var tax []models.Tax

	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, `SELECT t.id, t.tax_name, a.airline_name, a.country_name,  ct.city_name
											FROM tax t
													 INNER JOIN airline a ON t.iata_code = a.iata_code
													 INNER JOIN country c ON a.country_name = c.country_name
													 INNER JOIN city ct ON ct.country_iso2 = c.country_iso2
											WHERE t.tax_name IS NOT NULL AND t.tax_name != ''
       										OFFSET $1 LIMIT $2`,
		offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Tax
		err := rows.Scan(
			&t.ID, &t.TaxName, &t.AirlineName,
			&t.CountryName, &t.CityName,
		)

		if err != nil {
			return nil, err
		}
		tax = append(tax, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tax, nil
}

func (r *AirlineRepository) GetSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(id)
										FROM tax
										WHERE  tax.tax_name IS NOT NULL AND tax_name != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
