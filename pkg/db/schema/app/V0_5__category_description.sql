ALTER TABLE category
	ADD COLUMN description TEXT NOT NULL;

CREATE INDEX category_description_idx
	ON category (description(767));
