-- +goose Up
-- +goose StatementBegin

ALTER TABLE episodes ADD COLUMN airstamp_unix INTEGER;

UPDATE episodes
SET airstamp_unix = strftime('%s', airstamp);

CREATE TABLE episodes_new (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	show_id INTEGER NOT NULL,
	season_id INTEGER NOT NULL,

	name TEXT,
	summary TEXT,
	type TEXT,
	number INTEGER,
	season INTEGER NOT NULL,
	airstamp INTEGER, -- Now is INTEGER

	tracking TEXT DEFAULT 'wanted',
	torrent_hash TEXT,

	FOREIGN KEY(show_id) REFERENCES shows(id) ON DELETE CASCADE,
	FOREIGN KEY(season_id) REFERENCES seasons(id) ON DELETE CASCADE
);

INSERT INTO episodes_new (
	id, show_id, season_id,
	name, summary, type, number, season,
	airstamp, tracking, torrent_hash
)
SELECT
	id, show_id, season_id,
	name, summary, type, number, season,
	airstamp_unix, tracking, torrent_hash
FROM episodes;

DROP TABLE episodes;
ALTER TABLE episodes_new RENAME TO episodes;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

CREATE TABLE episodes_old (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	show_id INTEGER NOT NULL,
	season_id INTEGER NOT NULL,

	name TEXT,
	summary TEXT,
	type TEXT,
	number INTEGER,
	season INTEGER NOT NULL,
	airstamp TEXT, -- revert to TEXT (RFC3339)

	tracking TEXT DEFAULT 'wanted',
	torrent_hash TEXT,

	FOREIGN KEY(show_id) REFERENCES shows(id) ON DELETE CASCADE,
	FOREIGN KEY(season_id) REFERENCES seasons(id) ON DELETE CASCADE
);

INSERT INTO episodes_old (
	id, show_id, season_id,
	name, summary, type, number, season,
	airstamp, tracking, torrent_hash
)
SELECT
	id, show_id, season_id,
	name, summary, type, number, season,
	datetime(airstamp, 'unixepoch') AS airstamp,
	tracking, torrent_hash
FROM episodes;

DROP TABLE episodes;
ALTER TABLE episodes_old RENAME TO episodes;

-- +goose StatementEnd
