# go-spanner-ddlstructdiff

## Installation

```console
go install github.com/nametake/go-spanner-ddlstructdiff/cmd/ddlstructdiff@latest
```

## Usage

```console
go vet -vettool=`which ddlstructdiff` -ddlstructdiff.ddl=$(pwd)/ddl.sql .
```
