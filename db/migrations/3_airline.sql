CREATE TABLE "airline"
(
    id                     UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    fleet_average_age      FLOAT                    DEFAULT 0,
    airline_id             INT                      DEFAULT 0,
    callsign               varchar(100),
    hub_code               varchar(10),
    iata_code              varchar(10),
    icao_code              varchar(10),
    country_iso2           varchar(10),
    date_founded           INT                      DEFAULT 0,
    iata_prefix_accounting INT                      DEFAULT 0,
    airline_name           varchar(100),
    country_name           varchar(100),
    fleet_size             INT                      DEFAULT 0,
    status                 varchar(100),
    type                   varchar(100),
    created_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS aircraft
(
    id            UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    iata_code     varchar(100),
    aircraft_name varchar(100),
    plane_type_id INT                      DEFAULT 0,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE "airplane"
(
    id                       UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    iata_type                varchar(40),
    airplane_id              INT                      DEFAULT 0,
    airline_iata_code        varchar(20),
    iata_code_long           varchar(20),
    iata_code_short          varchar(20),
    airline_icao_code        varchar(20),
    construction_number      varchar(40),
    delivery_date            TIMESTAMP WITH TIME ZONE,
    engines_count            INT                      DEFAULT 0,
    engines_type             varchar(10),
    first_flight_date        TIMESTAMP WITH TIME ZONE,
    icao_code_hex            varchar(20),
    line_number              varchar(20),
    model_code               varchar(50),
    registration_number      varchar(50),
    test_registration_number varchar(50),
    plane_age                INT                      DEFAULT 0,
    plane_class              varchar(100),
    model_name               varchar(100),
    plane_owner              varchar(100),
    plane_series             varchar(100),
    plane_status             varchar(100),
    production_line          varchar(100),
    registration_date        TIMESTAMP,
    rollout_date             TIMESTAMP,
    created_at               TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE "tax"
(
    id         UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    tax_id     INT                      DEFAULT 0,
    tax_name   varchar(255),
    iata_code  varchar(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

create index on "tax" (id);
create index on "tax" (tax_id);
create index on "airplane" (id);
create index on "airplane" (airplane_id);
create index on "aircraft" (id);
create index on "aircraft" (plane_type_id);
create index on "airline" (id);
create index on "airline" (airline_id);
