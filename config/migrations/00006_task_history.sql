-- +goose Up
-- +goose StatementBegin
CREATE TABLE task_history (
    id CHAR(36) NOT NULL DEFAULT (UUID()),
    task_name VARCHAR(128) NOT NULL,
    task_status ENUM('Succeeded', 'Failed') NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    details VARCHAR(255) NOT NULL,

    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task_history;
-- +goose StatementEnd
