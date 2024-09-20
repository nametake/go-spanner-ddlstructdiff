package nofield

type Singer struct { // want "Singer struct must contain lastname field corresponding to DDL"
	SingerId   int64
	FirstName  string
	SingerInfo string
}
