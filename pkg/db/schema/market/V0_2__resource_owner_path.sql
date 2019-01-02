
ALTER TABLE market
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX market_owner_path_idx
	ON market (owner_path);
UPDATE market
SET owner_path = CONCAT(':', owner);

ALTER TABLE market_user
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX market_user_owner_path_idx
	ON market_user (owner_path);
UPDATE market_user
SET owner_path = CONCAT(':', owner);
