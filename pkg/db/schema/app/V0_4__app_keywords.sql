ALTER TABLE app
	ADD COLUMN keywords TEXT NOT NULL;

CREATE INDEX app_keywords_idx
	ON app (keywords(767));
