package nofield

type Singer struct { // want "Singer struct must contain LastName field corresponding to DDL" "Singer table does not have a column corresponding to Lastname"
	SingerId   int64
	FirstName  string
	Lastname   string
	SingerInfo string
}
