-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;

ALTER TABLE subscriptions
DROP COLUMN plan;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN is_admin;

ALTER TABLE subscriptions
ADD COLUMN plan VARCHAR(255);
-- +goose StatementEnd
