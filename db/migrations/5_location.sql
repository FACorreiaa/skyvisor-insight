CREATE TABLE city
(
    id           UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    gmt          varchar(5),
    city_id      INT                      DEFAULT 0,
    iata_code    varchar(20),
    country_iso2 varchar(20),
    geoname_id   varchar(100),
    latitude     float8                   DEFAULT 0.0,
    longitude    float8                   DEFAULT 0.0,
    city_name    varchar(100),
    timezone     varchar(100),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE country
(
    id                  UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    country_name        varchar(100),
    country_iso2        varchar(20),
    country_iso3        varchar(20),
    country_iso_numeric INT                      DEFAULT 0,
    population          INT                      DEFAULT 0,
    capital             varchar(100),
    continent           varchar(100),
    currency_name       varchar(100),
    currency_code       varchar(20),
    fips_code           varchar(20),
    phone_prefix        varchar(20),
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

create index on "city" (id);
create index on "city" (city_id);
create index on "country" (id);
create index on "country" (country_iso_numeric);
