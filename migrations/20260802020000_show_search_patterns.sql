-- +goose Up
-- +goose StatementBegin

CREATE TABLE show_search_patterns (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	show_id INTEGER NOT NULL,
	pattern TEXT NOT NULL,
	position INTEGER NOT NULL DEFAULT 0,

	FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE
);

CREATE INDEX idx_show_search_patterns_show_id ON show_search_patterns(show_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS show_search_patterns;

-- +goose StatementEnd
