package dialogue_state_storage

//go:generate mockgen -source=istorage.go -destination=./mocks/mocks.go -package=mocks IStorage

type IStorage interface {
	Set(userID int64, state int)
	Get(userID int64) int
	DeleteState(userId int64)
}
