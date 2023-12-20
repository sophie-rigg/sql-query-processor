-- name: test-query-1
SELECT * FROM test WHERE id = $1;

-- name: test-query-2
UPDATE test SET name = $1 WHERE id = $2;