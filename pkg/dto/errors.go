package dto

type Error struct {
	ServiceName string
	Err         error
	Code        int
}
