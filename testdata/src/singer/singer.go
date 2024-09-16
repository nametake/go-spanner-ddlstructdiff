package singer

type Singer struct { // want "Singer struct must contain LastName field corresponding to DDL"
	SingerId   int64
	FirstName  string
	SingerInfo string
}
