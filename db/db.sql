/*
 * Extensions
 */

CREATE EXTENSION IF NOT EXISTS citext;

/*
 * End of extensions
 */

/*
 * Tables creation
 */

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE,
    fullname VARCHAR(128) NOT NULL,
    about TEXT,
    email CITEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS forums (
    title VARCHAR(128) NOT NULL,
    "user" CITEXT NOT NULL REFERENCES users(nickname),
    slug CITEXT PRIMARY KEY,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS threads (
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    author CITEXT NOT NULL REFERENCES users(nickname),
    forum CITEXT NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    message TEXT NOT NULL,
    votes INT DEFAULT 0,
    slug CITEXT UNIQUE,
    created TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_forum (
    "user" CITEXT NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    forum CITEXT NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    PRIMARY KEY ("user", forum)
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    parent INT,
    author CITEXT NOT NULL REFERENCES users(nickname),
    message TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE,
    thread INT REFERENCES threads(id) ON DELETE CASCADE,
    created TIMESTAMP,
    post_tree INT[]
);

CREATE TABLE IF NOT EXISTS votes (
    thread_id INT NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    nickname CITEXT NOT NULL REFERENCES users(nickname),
    voice INT NOT NULL,
    PRIMARY KEY (thread_id, nickname)
);

/*
 * End of tables creation
 */


/*
 * Trigger procedures
 */

CREATE OR REPLACE FUNCTION on_insert_vote_trigger_proc()
RETURNS TRIGGER AS $$
BEGIN
UPDATE threads SET votes = votes + NEW.voice WHERE id = NEW.thread_id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION on_update_vote_trigger_proc()
RETURNS TRIGGER AS $$
BEGIN
UPDATE threads SET votes = votes + NEW.voice - OLD.voice WHERE id = NEW.thread_id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION on_insert_post_trigger_proc()
RETURNS TRIGGER AS $$
BEGIN
NEW.post_tree = (SELECT post_tree FROM posts WHERE id = NEW.parent) || NEW.id;
UPDATE forums SET posts = posts + 1 WHERE slug = NEW.forum;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION on_insert_thread_trigger_proc()
RETURNS TRIGGER AS $$
BEGIN
UPDATE forums SET threads = threads + 1 WHERE slug = NEW.forum;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

/*
 * End of trigger procedures
 */


/*
 * Triggers
 */

CREATE TRIGGER on_insert_vote_trigger
AFTER INSERT ON votes
FOR EACH ROW
EXECUTE PROCEDURE on_insert_vote_trigger_proc();

CREATE TRIGGER on_update_vote_trigger
AFTER UPDATE ON votes
FOR EACH ROW
EXECUTE PROCEDURE on_update_vote_trigger_proc();

CREATE TRIGGER on_insert_post_trigger
BEFORE INSERT ON posts
FOR EACH ROW
EXECUTE PROCEDURE on_insert_post_trigger_proc();

CREATE TRIGGER on_insert_thread_trigger
AFTER INSERT ON threads
FOR EACH ROW
EXECUTE PROCEDURE on_insert_thread_trigger_proc();

/*
 * End of functions and triggers
 */


/*
 * Indexes
 */

-- posts
CREATE INDEX IF NOT EXISTS posts_id_index on posts (id);
CREATE INDEX IF NOT EXISTS posts_thread_id_index on posts (thread, id);
CREATE INDEX IF NOT EXISTS posts_parent_thread_id_index on posts (parent, thread, id);
CREATE INDEX IF NOT EXISTS posts_post_tree_one_post_tree_index on posts ((post_tree[1]), post_tree);

-- forums
CREATE INDEX IF NOT EXISTS forums_user_nickname_index ON forums ("user");

-- users_forum
CREATE INDEX IF NOT EXISTS users_forum_forum_user_nickname_index ON users_forum (forum, "user");

-- threads
CREATE INDEX IF NOT EXISTS threads_slug_index ON threads (slug);
CREATE INDEX IF NOT EXISTS threads_forum_created ON threads (forum, created);

-- users
CREATE INDEX IF NOT EXISTS index_users_email_nickname ON users (email, nickname);

/*
 * End of indexes
 */


/*
 * Vacuum for updated. All tables gets updated
 */
VACUUM ANALYZE;
