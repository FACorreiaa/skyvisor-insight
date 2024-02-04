package location

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func (r *RepositoryLocation) getCountryData(ctx context.Context, query string,
	args ...interface{}) ([]models.Country, error) {
	var country []models.Country

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Country
		err := rows.Scan(&c.CountryName,
			&c.ID, &c.Capital, &c.CurrencyName,
			&c.Continent,
			&c.Population, &c.CurrencyCode,
			&c.CurrencyName, &c.Latitude, &c.Longitude,
		)

		if err != nil {
			return nil, err
		}
		country = append(country, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return country, nil
}

func (r *RepositoryLocation) GetCountry(ctx context.Context, page, pageSize int) ([]models.Country, error) {
	offset := (page - 1) * pageSize
	query := `SELECT
			    DISTINCT ON (cou.country_name)
			    cou.country_name,
			    cou.id, cou.capital,
			    cou.currency_name,
			    cou.continent,
			    cou.population,
			    cou.currency_code,
			    cou.currency_name,
			    ct.latitude,
			    ct.longitude
			FROM
			    country cou
			        LEFT JOIN
			    city ct ON ct.city_name = cou.capital
			WHERE
			    cou.country_name IS NOT NULL
			  AND TRIM(UPPER(cou.country_name)) != ''
			  AND ct.latitude IS NOT NULL
  			  AND ct.longitude IS NOT NULL
			ORDER BY
			    cou.country_name
            OFFSET $1 LIMIT $2`

	return r.getCountryData(ctx, query, offset, pageSize)
}

func (r *RepositoryLocation) GetCountryLocation(ctx context.Context) ([]models.Country, error) {
	query := `SELECT
			    DISTINCT ON (cou.country_name)
			    cou.country_name,
			    cou.id, cou.capital,
			    cou.continent,
			    cou.currency_name,
			    cou.population,
			    cou.currency_code,
			    cou.currency_name,
			    ct.latitude,
			    ct.longitude
			FROM
			    country cou
			        LEFT JOIN
			    city ct ON ct.city_name = cou.capital
			WHERE
			    cou.country_name IS NOT NULL
			  AND TRIM(UPPER(cou.country_name)) != ''
			  AND ct.latitude IS NOT NULL
  			  AND ct.longitude IS NOT NULL
			ORDER BY
			    cou.country_name`

	return r.getCountryData(ctx, query)
}

func (r *RepositoryLocation) GetCountrySum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `
				SELECT COUNT(DISTINCT cou.country_name)
				FROM country cou
				LEFT JOIN
			    	city ct ON ct.city_name = cou.capital
				WHERE cou.country_name IS NOT NULL
				AND ct.latitude IS NOT NULL
  			  	AND ct.longitude IS NOT NULL
				AND TRIM(UPPER(cou.country_name)) != '';`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
