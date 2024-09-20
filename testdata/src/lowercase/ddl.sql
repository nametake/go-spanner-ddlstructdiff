CREATE TABLE Singer (
  singerid   INT64 NOT NULL,
  firstname  STRING(1024),
  lastname   STRING(1024),
  singerinfo BYTES(MAX),
) PRIMARY KEY (SingerId);
