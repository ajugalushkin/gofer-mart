-- +goose Up
-- +goose StatementBegin
CREATE TABLE Users (
    ID BIGSERIAL NOT NULL ,
    Login VARCHAR(250) UNIQUE,
    Password_Hash VARCHAR(250),
    PRIMARY KEY (ID)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Users;
-- +goose StatementEnd
