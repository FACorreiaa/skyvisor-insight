package airline

import (
	"context"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

// tax

func (r *RepositoryAirline) getTaxData(ctx context.Context, query string,
	args ...interface{}) ([]models.Tax, error) {
	var tax []models.Tax

	rows, err := r.pgpool.Query(ctx, query, args...)
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

func (r *RepositoryAirline) GetTax(ctx context.Context, page, pageSize int,
	orderBy string, sortBy string, name string) ([]models.Tax, error) {
	query := `SELECT 
    										t.id, t.tax_name, a.airline_name, a.country_name
											FROM tax t
											         JOIN airline a ON t.iata_code = a.iata_code
											         JOIN country c ON a.country_name = c.country_name
											WHERE t.tax_name IS NOT NULL
											AND t.tax_name != ''
											AND    Trim(Upper(tax_name))
											           ILIKE trim(upper('%' || $1 || '%'))
											ORDER BY
			    CASE WHEN $2 = 'Tax Name' AND $3 = 'ASC' THEN t.tax_name::text END ASC,
			    CASE WHEN $2 = 'Tax Name' AND $3 = 'DESC' THEN t.tax_name::text END DESC,
			    CASE WHEN $2 = 'Airline Name' AND $3 = 'ASC' THEN a.airline_name::text END ASC,
			    CASE WHEN $2 = 'Airline Name' AND $3 = 'DESC' THEN a.airline_name::text END DESC,
			    CASE WHEN $2 = 'Country Name' AND $3 = 'ASC' THEN a.country_name::text END ASC,
			    CASE WHEN $2 = 'Country Name' AND $3 = 'DESC' THEN a.country_name::text END DESC
       										OFFSET $4 LIMIT $5`

	offset := (page - 1) * pageSize
	return r.getTaxData(ctx, query, name, orderBy, sortBy, offset, pageSize)
}

func (r *RepositoryAirline) GetSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(DISTINCT t.id)
										FROM tax t
										JOIN airline a ON t.iata_code = a.iata_code
										JOIN country c ON a.country_name = c.country_name
										WHERE t.tax_name IS NOT NULL AND t.tax_name != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
