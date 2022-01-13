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
CREATE VIEW IF NOT EXISTS vw_leaderboard
AS
SELECT
    s.id AS suggestionID
    , s.weekID
    , s.movie
    , v.preference
    , COUNT(s.id) AS votes
FROM suggestions s
INNER JOIN votes v
    ON v.suggestionID = s.id
    AND v.weekID = s.weekID
GROUP BY
    s.id
    , s.weekID
    , s.movie
    , v.preference
ORDER BY v.preference ASC, movie ASC;