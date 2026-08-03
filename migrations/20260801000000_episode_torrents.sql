-- +goose Up
-- +goose StatementBegin

CREATE TABLE episode_torrents(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	episode_id INTEGER NOT NULL UNIQUE,

	hash TEXT NOT NULL,
	title TEXT,
	indexer TEXT,
	info_url TEXT,
	publish_date TEXT,
	created_at INTEGER,

	FOREIGN KEY(episode_id) REFERENCES episodes(id) ON DELETE CASCADE
);

CREATE INDEX idx_episode_torrents_hash ON episode_torrents(hash);

INSERT INTO episode_torrents (episode_id, hash, created_at)
SELECT id, torrent_hash, strftime('%s', 'now')
FROM episodes
WHERE torrent_hash IS NOT NULL AND torrent_hash != '';

ALTER TABLE episodes DROP COLUMN torrent_hash;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE episodes ADD COLUMN torrent_hash TEXT;

UPDATE episodes
SET torrent_hash = (
	SELECT et.hash
	FROM episode_torrents et
	WHERE et.episode_id = episodes.id
)
WHERE id IN (SELECT episode_id FROM episode_torrents);

DROP TABLE episode_torrents;

-- +goose StatementEnd
