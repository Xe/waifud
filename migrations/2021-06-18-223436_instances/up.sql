CREATE TABLE IF NOT EXISTS connections
  ( hostname TEXT NOT NULL PRIMARY KEY
  , connection_uri TEXT NOT NULL
  );

CREATE TABLE IF NOT EXISTS instances
  ( id        TEXT           NOT NULL PRIMARY KEY
  , "name"    TEXT    UNIQUE NOT NULL
  , ram       INTEGER        NOT NULL DEFAULT 512
  , cores     INTEGER        NOT NULL DEFAULT 2
  , zvol      TEXT           NOT NULL
  , zvol_size INTEGER        NOT NULL DEFAULT 25
  , use_sata  BOOLEAN                 DEFAULT FALSE
  , owner     TEXT           NOT NULL
  , FOREIGN KEY(owner) REFERENCES connections(hostname)
  );


