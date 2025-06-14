CREATE TABLE threads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    thread_id UUID NOT NULL REFERENCES threads(id) on DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    votes INTEGER DEFAULT 0
);

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) on DELETE CASCADE,
    content TEXT NOT NULL,
    votes INTEGER DEFAULT 0
);