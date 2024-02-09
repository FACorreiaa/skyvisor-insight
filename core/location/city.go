package location

import (
	"context"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryLocation struct {
	pgpool *pgxpool.Pool
}

func NewLocations(
	pgpool *pgxpool.Pool,

) *RepositoryLocation {
	return &RepositoryLocation{
		pgpool: pgpool,
	}
}

func (r *RepositoryLocation) getCityData(ctx context.Context, query string,
	args ...interface{}) ([]models.City, error) {
	var city []models.City

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.City
		err := rows.Scan(&c.ID, &c.CityName, &c.Timezone, &c.GMT,
			&c.Continent, &c.CountryName, &c.CurrencyName,
			&c.PhonePrefix, &c.Latitude, &c.Longitude,
		)

		if err != nil {
			return nil, err
		}
		city = append(city, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return city, nil
}

func (r *RepositoryLocation) GetCity(ctx context.Context, page, pageSize int,
	orderBy string, sortBy string) ([]models.City, error) {
	offset := (page - 1) * pageSize
	query := `SELECT
			    ct.id,
			    ct.city_name,
			    ct.timezone,
			    ct.gmt,
			    cou.continent,
			    cou.country_name,
			    cou.currency_name,
			    cou.phone_prefix,
			    ct.latitude,
			    ct.longitude
			FROM city ct
			LEFT JOIN country cou ON cou.country_iso2 = ct.country_iso2
			WHERE ct.city_name IS NOT NULL AND TRIM(UPPER(ct.city_name)) != ''
			ORDER BY
		    CASE WHEN $2 = 'City Name' AND $3 = 'ASC' THEN ct.city_name::text END ASC,
		    CASE WHEN $2 = 'City Name' AND $3 = 'DESC' THEN ct.city_name::text END DESC,
		    CASE WHEN $2 = 'Timezone' AND $3 = 'ASC' THEN ct.timezone::text END ASC,
		    CASE WHEN $2 = 'Timezone' AND $3 = 'DESC' THEN ct.timezone::text END DESC,
			CASE WHEN $2 = 'GMT' AND $3 = 'ASC' THEN ct.gmt::text END ASC,
		    CASE WHEN $2 = 'GMT' AND $3 = 'DESC' THEN ct.gmt::text END DESC,
		    CASE WHEN $2 = 'Continent' AND $3 = 'ASC' THEN cou.continent::text END ASC,
		    CASE WHEN $2 = 'Continent' AND $3 = 'DESC' THEN cou.continent::text END DESC,
		    CASE WHEN $2 = 'Country Name' and $3 = 'ASC' THEN cou.country_name::text END ASC,
		    CASE WHEN $2 = 'Country Name' and $3 = 'DESC' THEN cou.country_name::text END DESC,
		    CASE WHEN $2 = 'Currency Name' and $3 = 'ASC' THEN cou.currency_name::text END ASC,
		    CASE WHEN $2 = 'Currency Name' and $3 = 'DESC' THEN cou.currency_name::text END DESC,
 			CASE WHEN $2 = 'Phone Prefix' and $3 = 'ASC' THEN cou.phone_prefix::text END ASC,
		    CASE WHEN $2 = 'Phone Prefix' and $3 = 'DESC' THEN cou.phone_prefix::text END DESC,
		    CASE WHEN $2 = 'Latitude' and $3 = 'ASC' THEN ct.latitude::text END ASC,
		    CASE WHEN $2 = 'Latitude' and $3 = 'DESC' THEN ct.latitude::text END DESC,
		    CASE WHEN $2 = 'Longitude' and $3 = 'ASC' THEN ct.longitude::text END ASC,
		    CASE WHEN $2 = 'Longitude' and $3 = 'DESC' THEN ct.longitude::text END DESC
			OFFSET $1 LIMIT $4;`

	return r.getCityData(ctx, query, offset, orderBy, sortBy, pageSize)
}

func (r *RepositoryLocation) GetCityLocation(ctx context.Context) ([]models.City, error) {
	query := `select DISTINCT ON(ct.city_name)
              ct.id,
              ct.city_name, ct.timezone, ct.gmt,
              cou.continent, cou.country_name,
              cou.currency_name, cou.phone_prefix,
              ct.latitude, ct.longitude
            from city ct
            left join
              country cou on cou.country_iso2 = ct.country_iso2
            where
              ct.city_name IS NOT NULL AND TRIM(UPPER(ct.city_name)) != ''
            order by
              ct.city_name`

	return r.getCityData(ctx, query)
}

func (r *RepositoryLocation) GetCitySum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT COUNT(DISTINCT ct.city_name)
										FROM city ct
										LEFT JOIN country cou ON cou.country_iso2 = ct.country_iso2
										WHERE ct.city_name IS NOT NULL AND TRIM(UPPER(ct.city_name)) != '';`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RepositoryLocation) GetCityByName(ctx context.Context, page, pageSize int,
	name string, orderBy string, sortBy string) ([]models.City, error) {
	offset := (page - 1) * pageSize
	query := `SELECT ct.id,
              ct.city_name, ct.timezone, ct.gmt,
              cou.continent, cou.country_name,
              cou.currency_name, cou.phone_prefix,
              ct.latitude, ct.longitude
              FROM city ct
              LEFT JOIN
              	country cou on cou.country_iso2 = ct.country_iso2
              WHERE
              	ct.city_name IS NOT NULL AND TRIM(UPPER(ct.city_name)) != ''
              AND    Trim(Upper(city_name)) ILIKE trim(upper('%'
				                  || $1
				                  || '%'))
	            ORDER BY
			    CASE WHEN $2 = 'City Name' AND $3 = 'ASC' THEN ct.city_name::text END ASC,
			    CASE WHEN $2 = 'City Name' AND $3 = 'DESC' THEN ct.city_name::text END DESC,
			    CASE WHEN $2 = 'Timezone' AND $3 = 'ASC' THEN ct.timezone::text END ASC,
			    CASE WHEN $2 = 'Timezone' AND $3 = 'DESC' THEN ct.timezone::text END DESC,
				CASE WHEN $2 = 'GMT' AND $3 = 'ASC' THEN ct.gmt::text END ASC,
			    CASE WHEN $2 = 'GMT' AND $3 = 'DESC' THEN ct.gmt::text END DESC,
			    CASE WHEN $2 = 'Continent' AND $3 = 'ASC' THEN cou.continent::text END ASC,
			    CASE WHEN $2 = 'Continent' AND $3 = 'DESC' THEN cou.continent::text END DESC,
			    CASE WHEN $2 = 'Country Name' and $3 = 'ASC' THEN cou.country_name::text END ASC,
			    CASE WHEN $2 = 'Country Name' and $3 = 'DESC' THEN cou.country_name::text END DESC,
			    CASE WHEN $2 = 'Currency Name' and $3 = 'ASC' THEN cou.currency_name::text END ASC,
			    CASE WHEN $2 = 'Currency Name' and $3 = 'DESC' THEN cou.currency_name::text END DESC,
	            CASE WHEN $2 = 'Phone Prefix' and $3 = 'ASC' THEN cou.phone_prefix::text END ASC,
			    CASE WHEN $2 = 'Phone Prefix' and $3 = 'DESC' THEN cou.phone_prefix::text END DESC,
			    CASE WHEN $2 = 'Latitude' and $3 = 'ASC' THEN ct.latitude::text END ASC,
			    CASE WHEN $2 = 'Latitude' and $3 = 'DESC' THEN ct.latitude::text END DESC,
			    CASE WHEN $2 = 'Longitude' and $3 = 'ASC' THEN ct.longitude::text END ASC,
			    CASE WHEN $2 = 'Longitude' and $3 = 'DESC' THEN ct.longitude::text END DESC
	            OFFSET $4 LIMIT $5`

	return r.getCityData(ctx, query, name, orderBy, sortBy, offset, pageSize)
}
