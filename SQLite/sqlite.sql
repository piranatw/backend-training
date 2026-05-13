-- Table: movies
CREATE TABLE movies (
    id integer NOT NULL CONSTRAINT movies_pk PRIMARY KEY AUTOINCREMENT,
    title varchar(255) NOT NULL
);

-- Table: reviews
CREATE TABLE reviews (
    id integer NOT NULL CONSTRAINT reviews_pk PRIMARY KEY,
    stars int NOT NULL,
    comment varchar(255) NOT NULL,
    movie_id int NOT NULL,
    CONSTRAINT reviews_movies FOREIGN KEY (movie_id)
    REFERENCES movies (id)
);

-- Insert a single movie
INSERT INTO movies ("title") VALUES ('Spider-Man: Into the Spider-Verse');

-- Insert multiple movies at once
INSERT INTO movies ("title") VALUES
    ('The Amazing Spider-Man'),
    ('Spider-Man: Homecoming'),
    ('Spider-Man (Sam Raimi)');


-- Select all movies
SELECT * FROM movies;
-- 1|Spider-Man: Into the Spider-Verse
-- 2|The Amazing Spider-Man
-- 3|Spider-Man: Homecoming
-- 4|Spider-Man (Sam Raimi)

-- Select with a condition (WHERE clause)
SELECT * FROM movies WHERE id < 3;
-- 1|Spider-Man: Into the Spider-Verse
-- 2|The Amazing Spider-Man


UPDATE movies SET "title" = 'Spider-Man (Tobey Maguire)' WHERE id = 4;

-- Verify
SELECT * FROM movies WHERE id = 4;
-- 4|Spider-Man (Tobey Maguire)


INSERT INTO reviews ("stars", "comment", "movie_id") VALUES
    (5, 'Incredible animation and storytelling!', 1),
    (4, 'Great movie, solid pacing.', 1);

-- Verify
SELECT * FROM reviews;
-- 1|5|Incredible animation and storytelling!|1
-- 2|4|Great movie, solid pacing.|1


