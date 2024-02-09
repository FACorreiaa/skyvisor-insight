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

func (r *RepositoryLocation) GetCountry(ctx context.Context, page, pageSize int,
	orderBy string, sortBy string) ([]models.Country, error) {
	offset := (page - 1) * pageSize
	query := `SELECT
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
			    CASE WHEN $1 = 'Country Name' AND $2 = 'ASC' THEN cou.country_name::text END ASC,
			    CASE WHEN $1 = 'Country Name' AND $2 = 'DESC' THEN cou.country_name::text END DESC,
			    CASE WHEN $1 = 'Capital' AND $2 = 'ASC' THEN cou.capital::text END ASC,
			    CASE WHEN $1 = 'Capital' AND $2 = 'DESC' THEN cou.capital::text END DESC,
				CASE WHEN $1 = 'Continent' AND $2 = 'ASC' THEN cou.continent::text END ASC,
			    CASE WHEN $1 = 'Continent' AND $2 = 'DESC' THEN cou.continent::text END DESC,
			    CASE WHEN $1 = 'Currency Name' AND $2 = 'ASC' THEN cou.currency_name::text END ASC,
			    CASE WHEN $1 = 'Currency Name' AND $2 = 'DESC' THEN cou.currency_name::text END DESC,
			    CASE WHEN $1 = 'Currency Code' and $2 = 'ASC' THEN cou.currency_code::text END ASC,
			    CASE WHEN $1 = 'Currency Code' and $2 = 'DESC' THEN cou.currency_code::text END DESC,
			    CASE WHEN $1 = 'Population' and $2 = 'ASC' THEN cou.population::text END ASC,
			    CASE WHEN $1 = 'Population' and $2 = 'DESC' THEN cou.population::text END DESC,
	            CASE WHEN $1 = 'Phone Prefix' and $2 = 'ASC' THEN cou.phone_prefix::text END ASC,
			    CASE WHEN $1 = 'Phone Prefix' and $2 = 'DESC' THEN cou.phone_prefix::text END DESC,
			    CASE WHEN $1 = 'Latitude' and $2 = 'ASC' THEN ct.latitude::text END ASC,
			    CASE WHEN $1 = 'Latitude' and $2 = 'DESC' THEN ct.latitude::text END DESC,
			    CASE WHEN $1 = 'Longitude' and $2 = 'ASC' THEN ct.longitude::text END ASC,
			    CASE WHEN $1 = 'Longitude' and $2 = 'DESC' THEN ct.longitude::text END DESC
            OFFSET $3 LIMIT $4`

	return r.getCountryData(ctx, query, orderBy, sortBy, offset, pageSize)
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

func (r *RepositoryLocation) GetCountryByName(ctx context.Context, page, pageSize int,
	name string, orderBy string, sortBy string) ([]models.Country, error) {
	offset := (page - 1) * pageSize
	query := `SELECT
			    cou.country_name,
			    cou.id,
			    cou.capital,
			    cou.continent,
			    cou.population,
			    cou.currency_code,
			    cou.currency_name,
			    cou.phone_prefix,
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
			 AND    Trim(Upper(country_name)) ILIKE trim(upper('%'
				                  || $1
				                  || '%'))
			ORDER BY
			    CASE WHEN $2 = 'Country Name' AND $3 = 'ASC' THEN cou.country_name::text END ASC,
			    CASE WHEN $2 = 'Country Name' AND $3 = 'DESC' THEN cou.country_name::text END DESC,
			    CASE WHEN $2 = 'Capital' AND $3 = 'ASC' THEN cou.capital::text END ASC,
			    CASE WHEN $2 = 'Capital' AND $3 = 'DESC' THEN cou.capital::text END DESC,
				CASE WHEN $2 = 'Continent' AND $3 = 'ASC' THEN cou.continent::text END ASC,
			    CASE WHEN $2 = 'Continent' AND $3 = 'DESC' THEN cou.continent::text END DESC,
			    CASE WHEN $2 = 'Currency Name' AND $3 = 'ASC' THEN cou.currency_name::text END ASC,
			    CASE WHEN $2 = 'Currency Name' AND $3 = 'DESC' THEN cou.currency_name::text END DESC,
			    CASE WHEN $2 = 'Currency Code' and $3 = 'ASC' THEN cou.currency_code::text END ASC,
			    CASE WHEN $2 = 'Currency Code' and $3 = 'DESC' THEN cou.currency_code::text END DESC,
			    CASE WHEN $2 = 'Population' and $3 = 'ASC' THEN cou.population::text END ASC,
			    CASE WHEN $2 = 'Population' and $3 = 'DESC' THEN cou.population::text END DESC,
	            CASE WHEN $2 = 'Phone Prefix' and $3 = 'ASC' THEN cou.phone_prefix::text END ASC,
			    CASE WHEN $2 = 'Phone Prefix' and $3 = 'DESC' THEN cou.phone_prefix::text END DESC,
			    CASE WHEN $2 = 'Latitude' and $3 = 'ASC' THEN ct.latitude::text END ASC,
			    CASE WHEN $2 = 'Latitude' and $3 = 'DESC' THEN ct.latitude::text END DESC,
			    CASE WHEN $2 = 'Longitude' and $3 = 'ASC' THEN ct.longitude::text END ASC,
			    CASE WHEN $2 = 'Longitude' and $3 = 'DESC' THEN ct.longitude::text END DESC
            OFFSET $4 LIMIT $5`

	return r.getCountryData(ctx, query, name, orderBy, sortBy, offset, pageSize)
}
