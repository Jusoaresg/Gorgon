-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE shows( 
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	tv_maze_id INTEGER NOT NULL,

	name TEXT,
	type TEXT,
	language TEXT,
	status TEXT,
	premiered TEXT,
	ended TEXT,
	rating REAL,
	summary TEXT,
	updated INTEGER,

	tv_rage INTEGER,
	the_tvdbd INTEGER,
	imdb INTEGER,

	image_original TEXT,
	image_medium TEXT,

	genres TEXT
);

CREATE TABLE schedules(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	show_id INTEGER NOT NULL UNIQUE,
	time TEXT,
	days TEXT,

	FOREIGN KEY(show_id) REFERENCES shows(id) ON DELETE CASCADE
);

CREATE TABLE seasons(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	show_id INTEGER NOT NULL,
	season_number INTEGER,

	FOREIGN KEY(show_id) REFERENCES shows(id) ON DELETE CASCADE
	UNIQUE(show_id, season_number)
);

CREATE TABLE episodes(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	show_id INTEGER NOT NULL,
	season_id INTEGER NOT NULL,

	name TEXT,
	summary TEXT,
	type TEXT,
	number INTEGER,
	season INTEGER NOT NULL,
	airstamp TEXT,

	file_path TEXT,
	tracking TEXT DEFAULT 'wanted',
	torrent_hash TEXT,

	FOREIGN KEY(show_id) REFERENCES shows(id) ON DELETE CASCADE,
	FOREIGN KEY(season_id) REFERENCES seasons(id) ON DELETE CASCADE
);

CREATE TABLE episode_content(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	episode_id INTEGER NOT NULL,

	name TEXT,
	size REAL,
	is_seed BOOLEAN,

	FOREIGN KEY(episode_id) REFERENCES episodes(id) ON DELETE CASCADE
);

CREATE TABLE indexers(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	indexer_id INTEGER NOT NULL,
	name TEXT,
	enabled BOOLEAN,
	definition_name TEXT,
	indexers_urls TEXT,
	language TEXT
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE shows;
DROP TABLE schedules;
DROP TABLE seasons;
DROP TABLE episodes;
DROP TABLE episode_content;
DROP TABLE indexers;
