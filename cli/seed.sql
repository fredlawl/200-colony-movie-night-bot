PRAGMA foreign_keys = ON;
insert or replace into suggestions (id, uuid, weekID, author, movie, movieHash) values
(1, "0482d3ff-6f1b-4629-9179-d8ba77f38c6a", "202121", "liam", "test", "tst")
,(2, "42e4b7ea-04cc-467b-832d-4f46c701189e", "202121", "liam", "shreck", "shrck")
,(3, "13fecbe2-18a2-4ba3-97e3-fc6d6dd73103", "202121", "liam", "pooh", "ph");

insert or replace into votes (suggestionID, weekID, author, preference) VALUES
(1, "202121", "liam", 1)
,(2, "202121", "liam", 2)
,(3, "202121", "liam", 3)
,(1, "202121", "noah", 3)
,(2, "202121", "noah", 2)
,(3, "202121", "noah", 1)
,(1, "202121", "oliver", 1)
,(2, "202121", "oliver", 3)
,(3, "202121", "oliver", 2)
,(1, "202121", "william", 2)
,(2, "202121", "william", 1)
,(3, "202121", "william", 3)
,(1, "202121", "james", 1)
,(2, "202121", "james", 2)
,(3, "202121", "james", 3)
,(1, "202121", "sneaky", 1)
,(2, "202121", "sneaky", 2);

select * from suggestions;
select * from votes;
select * from vw_leaderboard;