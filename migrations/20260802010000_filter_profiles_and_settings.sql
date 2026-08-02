-- +goose Up
-- +goose StatementBegin

CREATE TABLE filter_profiles (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	created_at INTEGER NOT NULL DEFAULT (unixepoch()),
	updated_at INTEGER NOT NULL DEFAULT (unixepoch())
);

CREATE TABLE filter_patterns (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	profile_id INTEGER NOT NULL,
	kind TEXT NOT NULL CHECK (kind IN ('search', 'required', 'rejected', 'preferred')),
	pattern TEXT NOT NULL,
	score INTEGER NOT NULL DEFAULT 0,
	position INTEGER NOT NULL DEFAULT 0,

	FOREIGN KEY (profile_id) REFERENCES filter_profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_filter_patterns_profile_id ON filter_patterns(profile_id);

CREATE TABLE show_settings (
	show_id INTEGER PRIMARY KEY,
	filter_profile_id INTEGER,
	use_aliases INTEGER NOT NULL DEFAULT 1,
	only_latin INTEGER NOT NULL DEFAULT 1,
	created_at INTEGER NOT NULL DEFAULT (unixepoch()),
	updated_at INTEGER NOT NULL DEFAULT (unixepoch()),

	FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE,
	FOREIGN KEY (filter_profile_id) REFERENCES filter_profiles(id) ON DELETE SET NULL
);

CREATE TABLE app_settings (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS app_settings;
DROP TABLE IF EXISTS show_settings;
DROP TABLE IF EXISTS filter_patterns;
DROP TABLE IF EXISTS filter_profiles;

-- +goose StatementEnd
