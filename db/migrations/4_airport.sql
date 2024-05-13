CREATE TABLE airport
(
    id             UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    gmt            varchar(5),
    airport_id     INT                      DEFAULT 0,
    iata_code      varchar(20),
    city_iata_code varchar(20),
    icao_code      varchar(20),
    country_iso2   varchar(20),
    geoname_id     varchar(20),
    latitude       float8                   DEFAULT 0.0,
    longitude      float8                   DEFAULT 0.0,
    airport_name   varchar(100),
    country_name   varchar(100),
    phone_number   varchar(40),
    timezone       varchar(100),
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

create index on "airport" (id);
create index on "airport" (airport_id);
create index on "airport" (iata_code);


