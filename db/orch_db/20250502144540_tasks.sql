-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    arg1 REAL NOT NULL,
    arg2 REAL NOT NULL,
    operation TEXT NOT NULL,
    operation_time_ms INTEGER NOT NULL,
    expr_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expr_id) REFERENCES expressions(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
