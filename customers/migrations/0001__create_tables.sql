DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'customers') THEN
            CREATE DATABASE customers;
        END IF;
    END $$;

create table customers (
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

CREATE INDEX idx_customers_name ON customers(name);
CREATE INDEX idx_customers_surname ON customers(surname);
CREATE INDEX idx_customers_phone_number ON customers(phone_number);

alter table customers owner to postgres;