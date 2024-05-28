-- +goose Up
-- +goose StatementBegin
CREATE TABLE Withdrawals (
    Number BIGSERIAL not null REFERENCES Orders(Number) UNIQUE,
    Sum float,
    Processed_At timestamp default NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Withdrawals;
-- +goose StatementEnd
