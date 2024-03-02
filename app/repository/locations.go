package repository

import (
	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationsRepository struct {
	pgpool *pgxpool.Pool
}

func NewLocationsRepository(db *pgxpool.Pool) *LocationsRepository {
	return &LocationsRepository{pgpool: db}
}

// Country
func (r *LocationsRepository) getCountryData(ctx context.Context, query string,
	args ...interface{}) ([]models.Country, error) {
	var country []models.Country

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Country
		err := rows.Scan(
			&c.CountryName,
			&c.ID,
			&c.Capital,
			&c.CurrencyName,
			&c.Continent,
			&c.Population,
			&c.CurrencyCode,
			&c.CurrencyName,
			&c.Latitude,
			&c.Longitude,
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

func (r *LocationsRepository) GetCountry(ctx context.Context, page, pageSize int,
	orderBy, sortBy, name string) ([]models.Country, error) {
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
			JOIN
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

func (r *LocationsRepository) GetCountryLocation(ctx context.Context) ([]models.Country, error) {
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
			JOIN
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

func (r *LocationsRepository) GetCountrySum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `
				SELECT COUNT(DISTINCT cou.country_name)
				FROM country cou
				JOIN
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

func (r *LocationsRepository) GetCountryByName(ctx context.Context, name string) (models.Country, error) {
	var c models.Country
	query := `SELECT
			    cou.country_name,
			    cou.id,
			    cou.capital,
			    cou.currency_name,
			    cou.continent,
			    cou.population,
			    cou.currency_code,
			    cou.currency_name,
			    cou.country_iso2,
			    cou.country_iso3,
			    cou.country_iso_numeric,
			    ct.latitude,
			    ct.longitude
			FROM
			    country cou
			JOIN
			    city ct ON ct.city_name = cou.capital
			WHERE
			    cou.country_name IS NOT NULL
			  AND TRIM(UPPER(cou.country_name)) != ''
			  AND ct.latitude IS NOT NULL
  			  AND ct.longitude IS NOT NULL
			  AND country_name = $1
			`
	err := r.pgpool.QueryRow(ctx, query, name).Scan(
		&c.CountryName,
		&c.ID,
		&c.Capital,
		&c.CurrencyName,
		&c.Continent,
		&c.Population,
		&c.CurrencyCode,
		&c.CurrencyName,
		&c.CountryISO2,
		&c.CountryIso3,
		&c.CountryIsoNumeric,
		&c.Latitude,
		&c.Longitude,
	)

	if err != nil {
		return models.Country{}, err
	}

	return c, nil
}

// City

func (r *LocationsRepository) getCityData(ctx context.Context, query string,
	args ...interface{}) ([]models.City, error) {
	var city []models.City

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.City
		err := rows.Scan(&c.ID, &c.CityID, &c.CityName, &c.Timezone, &c.GMT,
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

func (r *LocationsRepository) GetCity(ctx context.Context, page, pageSize int,
	orderBy string, sortBy string, name string) ([]models.City, error) {
	offset := (page - 1) * pageSize
	query := `SELECT
			    ct.id,
			    ct.city_id,
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
			JOIN country cou ON cou.country_iso2 = ct.country_iso2
			WHERE ct.city_name IS NOT NULL AND TRIM(UPPER(ct.city_name)) != ''
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
			OFFSET $4 LIMIT $5;`

	return r.getCityData(ctx, query, name, orderBy, sortBy, offset, pageSize)
}

func (r *LocationsRepository) GetCityLocation(ctx context.Context) ([]models.City, error) {
	query := `select DISTINCT ON(ct.city_name)
								ct.id,
								ct.city_id,
								ct.city_name,
								ct.timezone,
								ct.gmt,
								cou.continent,
								cou.country_name,
								cou.currency_name,
								cou.phone_prefix,
								ct.latitude,
								ct.longitude
            from city ct
            join
              country cou on cou.country_iso2 = ct.country_iso2
            where
              ct.city_name IS NOT NULL AND TRIM(UPPER(ct.city_name)) != ''
            order by
              ct.city_name`

	return r.getCityData(ctx, query)
}

func (r *LocationsRepository) GetCitySum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT COUNT(DISTINCT ct.city_name)
										FROM city ct
										JOIN country cou ON cou.country_iso2 = ct.country_iso2
										WHERE ct.city_name IS NOT NULL AND TRIM(UPPER(ct.city_name)) != '';`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LocationsRepository) GetCityByName(ctx context.Context, page, pageSize int,
	name string, orderBy string, sortBy string) ([]models.City, error) {
	offset := (page - 1) * pageSize
	query := `SELECT ct.id,
              ct.city_name, ct.timezone, ct.gmt,
              cou.continent, cou.country_name,
              cou.currency_name, cou.phone_prefix,
              ct.latitude, ct.longitude
              FROM city ct
              JOIN
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

func (r *LocationsRepository) GetCityByID(ctx context.Context, cityID int) (models.City, error) {
	var c models.City
	query := `
			SELECT
				ct.city_id,
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
			JOIN country cou ON cou.country_iso2 = ct.country_iso2
			WHERE ct.city_name IS NOT NULL
			  AND TRIM(UPPER(ct.city_name)) != ''
			  AND ct.city_id = $1;
	`

	err := r.pgpool.QueryRow(ctx, query, cityID).Scan(
		&c.CityID, &c.CityName, &c.Timezone, &c.GMT,
		&c.Continent, &c.CountryName, &c.CurrencyName,
		&c.PhonePrefix, &c.Latitude, &c.Longitude,
	)

	if err != nil {
		return models.City{}, err
	}

	return c, nil
}
