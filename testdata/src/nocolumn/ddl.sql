CREATE TABLE Singer (
  SingerId   INT64 NOT NULL,
  FirstName  STRING(1024),
  SingerInfo BYTES(MAX),
) PRIMARY KEY (SingerId);
