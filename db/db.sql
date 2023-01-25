CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS users (
  nickname CITEXT PRIMARY KEY, 
  fullname TEXT NOT NULL, 
  about TEXT, 
  email CITEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS forums (
  title TEXT NOT NULL, 
  forumCreator CITEXT REFERENCES users (nickname), 
  slug CITEXT PRIMARY KEY, 
  posts INT DEFAULT 0, 
  threads INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS threads (
  id SERIAL PRIMARY KEY, 
  title TEXT NOT NULL, 
  author CITEXT REFERENCES users (nickname), 
  forum CITEXT REFERENCES forums (slug), 
  message TEXT NOT NULL, 
  slug CITEXT,
  votes INT DEFAULT 0, 
  created timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS posts (
  id SERIAL PRIMARY KEY, 
  parent INT REFERENCES posts (id), 
  author CITEXT REFERENCES users (nickname), 
  message TEXT NOT NULL, 
  edited BOOLEAN DEFAULT FALSE, 
  forum CITEXT REFERENCES forums (slug), 
  thread INT REFERENCES threads (id), 
  created timestamp with time zone DEFAULT now(), 
  path INT[] DEFAULT ARRAY [] :: INTEGER[]
);

CREATE TABLE IF NOT EXISTS votes (
  nickname CITEXT REFERENCES users (nickname), 
  thread INT REFERENCES threads (id), 
  voice INT NOT NULL,
  UNIQUE (nickname, thread)
);

CREATE TABLE IF NOT EXISTS userForum (
  nickname CITEXT  REFERENCES users (nickname), 
  forum CITEXT REFERENCES forums (slug),
  UNIQUE (nickname, forum)
);

-- Обновляет votes у thread после первоначального создания votes
CREATE 
OR REPLACE FUNCTION insertVotesFunc() RETURNS TRIGGER AS $$ BEGIN 
UPDATE 
  threads 
SET 
  votes = threads.votes + NEW.voice 
WHERE 
  id = NEW.thread;
RETURN NEW;
END;

$$ language plpgsql;
CREATE TRIGGER insertVotes 
AFTER 
  INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE insertVotesFunc();

-- Обновляет votes у thread после изменения votes
CREATE 
OR REPLACE FUNCTION updateVotesFunc() RETURNS TRIGGER AS $$ BEGIN 
UPDATE 
  threads 
SET 
  votes = threads.votes + NEW.voice - OLD.voice 
WHERE 
  id = NEW.thread;
RETURN NEW;
END;

$$ language plpgsql;
CREATE TRIGGER updateVotes 
AFTER 
UPDATE 
  ON votes FOR EACH ROW EXECUTE PROCEDURE updateVotesFunc();

-- Сохраянем связь с предыдущеми постами. 
-- Если пост первый в треде, то первый запрос будет пустой и мы положим
-- в path id текущего поста
CREATE 
OR REPLACE FUNCTION insertPostBeforeFunc() RETURNS TRIGGER AS $$ 
BEGIN NEW.path = (
  SELECT 
    path 
  FROM 
    posts 
  WHERE 
    id = NEW.parent
) || NEW.id;
RETURN NEW;
END;

$$ LANGUAGE plpgsql;
CREATE TRIGGER insertPostBefore 
BEFORE INSERT ON posts 
FOR EACH ROW EXECUTE PROCEDURE insertPostBeforeFunc();

-- Обновляет количество posts в forums на каждой вставке в posts
CREATE 
OR REPLACE FUNCTION insertPostAfterFunc() RETURNS TRIGGER AS $$ BEGIN 
UPDATE 
  forums 
SET 
  posts = forums.posts + 1 
WHERE 
  slug = NEW.forum;
RETURN NEW;
END;

$$ LANGUAGE plpgsql;
CREATE TRIGGER insertPostAfter 
AFTER 
  INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE insertPostAfterFunc();

-- Обновляет количество threads в forums на каждой вставке в threads
CREATE 
OR REPLACE FUNCTION insertThreadsFunc() RETURNS TRIGGER AS $$ BEGIN 
UPDATE 
  forums 
SET 
  threads = forums.threads + 1 
WHERE 
  slug = NEW.forum;
RETURN NEW;
END;

$$ LANGUAGE plpgsql;
CREATE TRIGGER insertThreads 
AFTER 
  INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE insertThreadsFunc();

-- После создания нового thread или post 
-- вставляем автора и форум, который именили в userForum
CREATE 
OR REPLACE FUNCTION addUser() RETURNS TRIGGER AS $$ 
BEGIN 
INSERT INTO userForum (nickname, forum) 
VALUES 
  (NEW.author, NEW.forum) ON CONFLICT do nothing;
RETURN NEW;
END;

$$ LANGUAGE plpgsql;
CREATE TRIGGER insertNewThread 
AFTER 
  INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE addUser();

CREATE TRIGGER insertNewPost 
AFTER 
  INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE addUser();

CREATE INDEX IF NOT EXISTS users_idx_nickname_and_email ON users (nickname, email);
CREATE INDEX IF NOT EXISTS posts_idx_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS threads_idx_forum_and_created ON threads (forum, created);
CREATE INDEX IF NOT EXISTS userForum_idx_nickname ON userForum (nickname);
CREATE INDEX IF NOT EXISTS posts_idx_thread_and_path_and_parent ON posts ((path[1]), thread, parent NULLS FIRST);
CREATE INDEX IF NOT EXISTS posts_idx_thread_and_parent_and_id ON posts (thread, id, parent NULLS FIRST);
CREATE INDEX IF NOT EXISTS posts_thread_thread_and_path ON posts (thread, path);
