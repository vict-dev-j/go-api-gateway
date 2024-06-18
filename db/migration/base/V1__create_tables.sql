-- todo: add migration
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'customers') THEN
            CREATE DATABASE customers;
        END IF;
    END $$;


DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'invest_accounts') THEN
            CREATE DATABASE customers;
        END IF;
    END $$;

create table customers
(
    id                                serial primary key,
    name                              text      not null,
    surname                           text      not null,
    age                               integer   not null,
    phone_number                      text      not null,
    debit_card                        text      not null,
    credit_card                       text      not null,
    date_of_birth                     timestamp not null,
    date_of_issue                     timestamp not null,
    issuing_authority                 text      not null,
    has_foreign_country_tax_liability boolean   not null
);

alter table customers
    owner to postgres;

create table invest_accounts
(
    id                       serial primary key,
    owner_id                 integer          not null,
    client_survey_number     integer          not null,
    share                    text             not null,
    invested_amount_of_money double precision not null,
    free_amount_of_money     double precision not null
);

alter table invest_accounts
    owner to postgres;



