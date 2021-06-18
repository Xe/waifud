CREATE TABLE IF NOT EXISTS templates
  ( uuid         TEXT UNIQUE NOT NULL PRIMARY KEY
  , "name"       TEXT UNIQUE NOT NULL
  , distro       TEXT        NOT NULL
  , version      TEXT        NOT NULL
  , download_url TEXT        NOT NULL
  , sha256sum    TEXT        NOT NULL
  , local_url    TEXT        NOT NULL
  );
