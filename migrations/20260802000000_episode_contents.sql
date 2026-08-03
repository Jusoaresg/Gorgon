-- +goose Up
-- +goose StatementBegin

ALTER TABLE episode_content RENAME TO episode_contents;

ALTER TABLE episode_contents DROP COLUMN is_seed;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE episode_contents ADD COLUMN is_seed BOOLEAN;

ALTER TABLE episode_contents RENAME TO episode_content;

-- +goose StatementEnd
