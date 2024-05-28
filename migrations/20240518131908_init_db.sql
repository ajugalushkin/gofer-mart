-- +goose Up
-- +goose StatementBegin
CREATE TABLE Users (
    Login VARCHAR(250) UNIQUE,
    Password VARCHAR(250),
    PRIMARY KEY (Login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Users;
-- +goose StatementEnd
