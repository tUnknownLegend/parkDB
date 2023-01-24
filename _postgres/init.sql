CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE UNLOGGED TABLE IF NOT EXISTS users (
  nickname CITEXT COLLATE "ucs_basic" PRIMARY KEY, 
  fullname VARCHAR(50) NOT NULL, 
  about text, 
  email CITEXT NOT NULL UNIQUE
);
CREATE UNLOGGED TABLE IF NOT EXISTS forums (
  title VARCHAR(50) NOT NULL, 
  user_ CITEXT REFERENCES users (nickname), 
  slug CITEXT PRIMARY KEY, 
  posts INT DEFAULT 0, 
  threads INT DEFAULT 0
);
CREATE UNLOGGED TABLE IF NOT EXISTS threads (
  id SERIAL PRIMARY KEY, 
  title VARCHAR(50) NOT NULL, 
  author CITEXT REFERENCES users (nickname), 
  forum CITEXT REFERENCES forums (slug), 
  message text NOT NULL, 
  votes INT DEFAULT 0, 
  slug CITEXT, 
  created timestamp with time zone DEFAULT now()
);
CREATE UNLOGGED TABLE IF NOT EXISTS posts (
  id SERIAL PRIMARY KEY, 
  parent INT REFERENCES posts (id), 
  author CITEXT REFERENCES users (nickname), 
  message text NOT NULL, 
  is_edited BOOLEAN DEFAULT FALSE, 
  forum CITEXT REFERENCES forums (slug), 
  thread INT REFERENCES threads (id), 
  created timestamp with time zone DEFAULT now(), 
  path INT[] DEFAULT ARRAY [] :: INTEGER[]
);
CREATE UNLOGGED TABLE IF NOT EXISTS votes (
  nickname CITEXT REFERENCES users (nickname), 
  thread INT REFERENCES threads (id), 
  voice INT NOT NULL, 
  constraINT user_thread_key UNIQUE (nickname, thread)
);
CREATE UNLOGGED TABLE IF NOT EXISTS user_forum (
  nickname CITEXT COLLATE "ucs_basic" REFERENCES users (nickname), 
  forum CITEXT REFERENCES forums (slug), 
  constraINT user_forum_key UNIQUE (nickname, forum)
);
CREATE 
OR REPLACE FUNCTION insert_votes_proc() RETURNS TRIGGER AS $$ BEGIN 
UPDATE 
  threads 
SET 
  votes = threads.votes + NEW.voice 
WHERE 
  id = NEW.thread;
RETURN NEW;
END;
$$ language plpgsql;
CREATE TRIGGER insert_votes 
AFTER 
  INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE insert_votes_proc();
CREATE 
OR REPLACE FUNCTION update_votes_proc() RETURNS TRIGGER AS $$ BEGIN 
UPDATE 
  threads 
SET 
  votes = threads.votes + NEW.voice - OLD.voice 
WHERE 
  id = NEW.thread;
RETURN NEW;
END;
$$ language plpgsql;
CREATE TRIGGER update_votes 
AFTER 
UPDATE 
  ON votes FOR EACH ROW EXECUTE PROCEDURE update_votes_proc();
CREATE 
OR REPLACE FUNCTION insert_post_before_proc() RETURNS TRIGGER AS $$ DECLARE parent_post_id posts.id % type := 0;
BEGIN NEW.path = (
  SELECT 
    path 
  FROM 
    posts 
  WHERE 
    id = new.parent
) || NEW.id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER insert_post_before BEFORE INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE insert_post_before_proc();
CREATE 
OR REPLACE FUNCTION insert_post_after_proc() RETURNS TRIGGER AS $$ BEGIN 
UPDATE 
  forums 
SET 
  posts = forums.posts + 1 
WHERE 
  slug = NEW.forum;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER insert_post_after 
AFTER 
  INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE insert_post_after_proc();
CREATE 
OR REPLACE FUNCTION insert_threads_proc() RETURNS TRIGGER AS $$ BEGIN 
UPDATE 
  forums 
SET 
  threads = forums.threads + 1 
WHERE 
  slug = NEW.forum;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER insert_threads 
AFTER 
  INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE insert_threads_proc();
CREATE 
OR REPLACE FUNCTION add_user() RETURNS TRIGGER AS $$ BEGIN INSERT INTO user_forum (nickname, forum) 
VALUES 
  (NEW.author, NEW.forum) ON CONFLICT do nothing;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER insert_new_thread 
AFTER 
  INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE add_user();
CREATE TRIGGER insert_new_post 
AFTER 
  INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE add_user();
