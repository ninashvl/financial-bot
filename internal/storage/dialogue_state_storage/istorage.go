package dialogue_state_storage

type IStorage interface {
	Set(userID int64, state int)
	Get(userID int64) int
	DeleteState(userId int64)
}
