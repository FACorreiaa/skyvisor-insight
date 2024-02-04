package airline

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

// tax

func (r *RepositoryAirline) GetTax(ctx context.Context, page, pageSize int) ([]models.Tax, error) {
	var tax []models.Tax

	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, `SELECT DISTINCT ON(t.id)
    										t.id, t.tax_name, a.airline_name, a.country_name
											FROM tax t
											         INNER JOIN airline a ON t.iata_code = a.iata_code
											         INNER JOIN country c ON a.country_name = c.country_name
											WHERE t.tax_name IS NOT NULL AND t.tax_name != ''
											ORDER BY t.id
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
			&t.CountryName,
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

func (r *RepositoryAirline) GetSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(DISTINCT t.id)
										FROM tax t
										INNER JOIN airline a ON t.iata_code = a.iata_code
										INNER JOIN country c ON a.country_name = c.country_name
										WHERE t.tax_name IS NOT NULL AND t.tax_name != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
