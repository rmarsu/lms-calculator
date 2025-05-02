-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS expressions (
     id INTEGER PRIMARY KEY AUTOINCREMENT,
     expression TEXT NOT NULL,
     status TEXT NOT NULL,
     result REAL NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS expressions;
-- +goose StatementEnd
