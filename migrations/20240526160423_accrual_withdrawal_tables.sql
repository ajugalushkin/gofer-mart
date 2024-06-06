-- +goose Up
-- +goose StatementBegin
CREATE TABLE Withdrawals (
    Number BIGSERIAL UNIQUE,
    Sum float,
    Processed_At timestamp default NULL,
    User_ID VARCHAR(250),
    PRIMARY KEY (Number),
    FOREIGN KEY (User_ID) REFERENCES Users(Login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Withdrawals;
-- +goose StatementEnd
