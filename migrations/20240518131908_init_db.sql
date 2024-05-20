-- +goose Up
-- +goose StatementBegin
CREATE TABLE Users (
    User_ID BIGSERIAL NOT NULL ,
    Login VARCHAR(250) UNIQUE,
    Password_Hash VARCHAR(250),
    PRIMARY KEY (User_ID)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Users;
-- +goose StatementEnd
