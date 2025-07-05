-- +goose Up
-- +goose StatementBegin
CREATE TABLE show_aliases (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	show_id INTEGER NOT NULL,
	alias TEXT NOT NULL,
	country TEXT,
	source TEXT NOT NULL DEFAULT 'user',

	FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS show_aliases;
-- +goose StatementEnd
