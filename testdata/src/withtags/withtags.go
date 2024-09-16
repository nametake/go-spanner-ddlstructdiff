package withtags

type Singer struct { // want "Singer struct must contain last_name field corresponding to DDL"
	SingerId   int64  `spanner:"singer_id"`
	FirstName  string `spanner:"first_name"`
	SingerInfo string `spanner:"singer_info"`
}
