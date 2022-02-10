CREATE TABLE IF NOT EXISTS instances
  ( uuid TEXT PRIMARY KEY NOT NULL
  , name TEXT NOT NULL UNIQUE
  , host TEXT NOT NULL
  , mac_address TEXT NOT NULL
  , memory INTEGER NOT NULL
  , disk_size INTEGER NOT NULL
  , zvol_name TEXT NOT NULL
  , status TEXT NOT NULL DEFAULT 'unknown'
  , distro TEXT NOT NULL
  );

CREATE TABLE IF NOT EXISTS audit_logs
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  , ts INTEGER NOT NULL DEFAULT (STRFTIME('%s', 'now'))
  , kind TEXT NOT NULL
  , op TEXT NOT NULL
  , data TEXT
  , uuid TEXT GENERATED ALWAYS AS (json_extract(data, '$.uuid'))
  , name TEXT GENERATED ALWAYS AS (json_extract(data, '$.name'))
  );

CREATE INDEX IF NOT EXISTS audit_logs_uuid
  ON audit_logs(uuid);

CREATE INDEX IF NOT EXISTS audit_logs_name
  ON audit_logs(name);

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
