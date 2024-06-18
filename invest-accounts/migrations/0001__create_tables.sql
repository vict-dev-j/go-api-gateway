DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'invest_accounts') THEN
            CREATE DATABASE customers;
        END IF;
    END $$;

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



