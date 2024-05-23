-- +goose Up
-- +goose StatementBegin
CREATE TABLE Orders (
    Order_ID BIGSERIAL NOT NULL,
    Order_Number BIGSERIAL UNIQUE,
    User_ID int,
    PRIMARY KEY (Order_ID),
    FOREIGN KEY (User_ID) REFERENCES Users(User_ID)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Orders;
-- +goose StatementEnd
