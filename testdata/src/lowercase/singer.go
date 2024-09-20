package nofield

type Singer struct { // want "singer struct must contain lastname field corresponding to DDL"
	SingerId   int64
	FirstName  string
	SingerInfo string
}
