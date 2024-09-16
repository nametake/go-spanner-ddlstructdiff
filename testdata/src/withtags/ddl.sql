CREATE TABLE singer (
  singer_id   INT64 NOT NULL,
  first_name  STRING(1024),
  last_name   STRING(1024),
  singer_info BYTES(MAX),
) PRIMARY KEY (singer_id);
