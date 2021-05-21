CREATE TABLE IF NOT EXISTS suggestions (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    uuid INTEGER NOT NULL,
    weekID INTEGER NOT NULL,
    author VARCHAR(255) NOT NULL,
    movie VARCHAR(255) NOT NULL,
    movieHash VARCHAR(255) NOT NULL,
    dateAdded DATETIME NOT NULL DEFAULT current_timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS ix_suggestions_uuid ON suggestions(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS ix_suggestions_movieHash ON suggestions(movieHash);
CREATE TABLE IF NOT EXISTS votes (
    suggestionID INTEGER NOT NULL,
    weekID INTEGER NOT NULL,
    author VARCHAR(255) NOT NULL,
    preference INTEGER NOT NULL,
    dateAdded DATETIME NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY(weekID, author, suggestionID),
    CONSTRAINT fk_votes_suggestionID FOREIGN KEY (suggestionID) REFERENCES suggestions(id) ON DELETE CASCADE 
);