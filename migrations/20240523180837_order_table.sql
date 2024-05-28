-- +goose Up
-- +goose StatementBegin
CREATE TABLE Orders (
    Number BIGSERIAL UNIQUE,
    Uploaded_At timestamp default NULL,
    Status VARCHAR(20),
    Accrual float,
    User_ID VARCHAR(250),
    PRIMARY KEY (Number),
    FOREIGN KEY (User_ID) REFERENCES Users(Login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Orders;
-- +goose StatementEnd
