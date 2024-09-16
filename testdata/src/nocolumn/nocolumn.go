package singer

type Singer struct { // want "Singer table does not have a column corresponding to LastName"
	SingerId   int64
	FirstName  string
	LastName   string
	SingerInfo string
}
