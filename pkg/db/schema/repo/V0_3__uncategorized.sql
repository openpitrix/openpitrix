INSERT INTO category_resource (category_id, resource_id, status)
	SELECT
		"ctg-uncategorized",
		repo_id,
		"enabled"
	FROM repo
	WHERE repo_id NOT IN (
		SELECT resource_id
		FROM category_resource
	);
