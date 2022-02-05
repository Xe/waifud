CREATE TABLE IF NOT EXISTS instances
  ( uuid TEXT PRIMARY KEY NOT NULL
  , name TEXT NOT NULL UNIQUE
  , host TEXT NOT NULL
  , memory INTEGER NOT NULL
  , disk_size INTEGER NOT NULL
  );

CREATE TABLE IF NOT EXISTS cloudconfig_seeds
  ( uuid TEXT PRIMARY KEY NOT NULL
  , user_data TEXT NOT NULL
  );

CREATE TABLE IF NOT EXISTS distros
  ( name TEXT PRIMARY KEY NOT NULL
  , download_url TEXT NOT NULL
  , sha256sum TEXT NOT NULL
  , min_size INTEGER NOT NULL
  , format TEXT NOT NULL
  );
