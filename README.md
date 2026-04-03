# go-spanner-ddlstructdiff

## Installation

```console
go install github.com/nametake/go-spanner-ddlstructdiff/cmd/ddlstructdiff@latest
```

## Usage

```console
go vet -vettool=`which ddlstructdiff` -ddlstructdiff.ddl=$(pwd)/ddl.sql .
```

## Ignoring lint errors

### Go struct

Add `//nolint:ddlstructdiff` to suppress lint errors.

**Struct-level** (suppresses all errors for the struct):

```go
//nolint:ddlstructdiff
type Singer struct {
    SingerId  int64
    FirstName string
    // LastName is intentionally omitted
}
```

**Field-level** (suppresses the "extra field" error for that field):

```go
type Singer struct {
    SingerId   int64
    FirstName  string
    ExtraField string //nolint:ddlstructdiff
}
```

### DDL schema

Add `-- nolint:ddlstructdiff` to suppress lint errors from the DDL side.

**Table-level** (suppresses all errors for the table):

```sql
-- nolint:ddlstructdiff
CREATE TABLE Singer (
  SingerId  INT64 NOT NULL,
  FirstName STRING(1024),
  LastName  STRING(1024),
) PRIMARY KEY (SingerId);
```

**Column-level** (suppresses the "missing field" error for that column):

```sql
CREATE TABLE Singer (
  SingerId  INT64 NOT NULL,
  FirstName STRING(1024),
  -- nolint:ddlstructdiff
  LastName  STRING(1024),
) PRIMARY KEY (SingerId);
```
