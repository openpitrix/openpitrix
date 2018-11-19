ALTER TABLE repo
	ADD COLUMN controller TINYINT NOT NULL DEFAULT 0;

CREATE INDEX repo_controller_idx
	ON repo (controller);

UPDATE repo
SET controller = 1
where repo_id in ('repo-vmbased', 'repo-helm');
