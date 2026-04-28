-- chunks stored in the CAS
CREATE TABLE Chunks (
    content_hash TEXT NOT NULL,
    size INTEGER NOT NULL,
    PRIMARY KEY (content_hash)
);

-- parity group ids
CREATE TABLE ParityGroups (
    parity_group_id INTEGER NOT NULL AUTOINCREMENT,
    PRIMARY KEY (parity_group_id)
);

-- mapping of content hashes to the nodes that have the chunk
CREATE TABLE ChunkNodes (
    content_hash TEXT NOT NULL,
    node_id TEXT NOT NULL,
    PRIMARY KEY (content_hash, node_id)
);

-- lists of content hashes for each parity group
CREATE TABLE ChunkParityGroups (
    parity_group_id INTEGER NOT NULL,
    index_in_group INTEGER NOT NULL,
    content_hash TEXT NOT NULL,
    PRIMARY KEY (content_hash, parity_group_id),
    FOREIGN KEY (parity_group_id) REFERENCES ParityGroups(parity_group_id),
    FOREIGN KEY (content_hash) REFERENCES Chunks(content_hash),
    UNIQUE (parity_group_id, index_in_group)
);