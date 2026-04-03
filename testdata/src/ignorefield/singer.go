package ignorefield

type Singer struct {
	SingerId   int64
	FirstName  string
	ExtraField string //nolint:ddlstructdiff
	SingerInfo string
}
