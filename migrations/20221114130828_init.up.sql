DO
$$
BEGIN

-- Таблица трат
CREATE TABLE IF NOT EXISTS expenses (
        id serial           PRIMARY KEY     NOT NULL,
        user_id             BIGINT          NOT NULL,
        amount              NUMERIC         NOT NULL,
        exp_category        VARCHAR         NOT NULL,
        exp_date            TIMESTAMP       NOT NULL,

        created_at          TIMESTAMP       DEFAULT NOW(),
        updated_at          TIMESTAMP,
        deleted_at          TIMESTAMP
);

-- Таблица пользовательских курсов
CREATE TABLE IF NOT EXISTS user_currency (
    user_id                 BIGINT      PRIMARY KEY     NOT NULL,
    currency_value          VARCHAR                     NOT NULL,

    created_at              TIMESTAMP   DEFAULT NOW(),
    updated_at              TIMESTAMP,
    deleted_at              TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_limit (
    user_id             BIGINT      PRIMARY KEY     NOT NULL,
    limit_value         NUMERIC                     NOT NULL,

    created_at          TIMESTAMP  DEFAULT NOW(),
    updated_at          TIMESTAMP,
    deleted_at          TIMESTAMP
);

-- Выбор поля user_id обусловлен условием запроса при генерации отчета
-- Тк в запросе прямое обращение по user_id - можно использовать в качестве индекса hash
CREATE INDEX user_id_index ON expenses USING hash(user_id);

-- На остальные таблицы индексы не нужных, тк обращение происходит на user_id, который является PRIMARY KEY

END

$$;