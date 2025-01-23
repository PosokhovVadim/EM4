INSERT INTO songs (group_name, name, link, release_date, inserted_at)
VALUES ('Imagine Dragons', 'Radioactive', 'https://www.youtube.com/watch?v=ktvTqknDZAU', '2012-07-02', NOW()),
       ('The Neighbourhood', 'Sweater Weather', 'https://www.youtube.com/watch?v=GKwUW3cdq4Y', '2012-12-03', NOW()),
       ('The 1975', 'The Sound', 'https://www.youtube.com/watch?v=-xKKo7CPHqc', '2016-02-25', NOW())
ON CONFLICT (group_name, name) DO NOTHING;

INSERT INTO lyrics (song_id, verse_number, text)
VALUES (1, 1, 'I raise my flags dye my clothes'),
       (1, 2, 'It s a revolution, I suppose'),
       (2, 1, 'And it s been two years I haven t seen her'),
       (2, 2, 'And I ve been living on such thin air'),
       (3, 1, 'Go ahead, go ahead, go ahead, go ahead'),
       (3, 2, 'Give me the taste, give me the taste, give me the taste, give me the taste'),
       (3, 3, 'I wanna give you what you want'),
       (3, 4, 'I wanna give you what you need');
