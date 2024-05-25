-- +goose Up
-- +goose StatementBegin
CREATE TABLE Orders (
    ID BIGSERIAL NOT NULL,
    Number BIGSERIAL UNIQUE,
    Uploaded_At timestamp default NULL,
    Status VARCHAR(20),
    User_ID int,
    PRIMARY KEY (ID),
    FOREIGN KEY (User_ID) REFERENCES Users(ID)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Orders;
-- +goose StatementEnd
