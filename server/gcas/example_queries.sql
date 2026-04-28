-- get the list of content hashes in a parity group in order
SELECT content_hash
FROM ChunkParityGroups
WHERE parity_group_id = ?
ORDER BY index_in_group ASC;

-- get the list of content hashes that lack a parity group
SELECT c.content_hash
FROM Chunks c
LEFT JOIN ChunkParityGroups cpg ON c.content_hash = cpg.content_hash
WHERE cpg.parity_group_id IS NULL;
