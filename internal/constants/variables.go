package constants

const (
	LevelManageAdmin int8 = iota + 1 // 1
	LevelManageUser                  // 2
)
const (
	StatusInactive = iota // 0
	StatusActive          // 1
)

const (
	OptionAnswerType int8 = iota + 1 // 1
	TextAnswerType                   // 2
)

const (
	ArticleInactive uint = iota // 0
	ArticleActive               // 1
	ArticleTop                  // 2
)

const (
	FileProcessFailed     uint = iota //0
	FileProcessReady                  //1
	FileProcessProcessing             //2
)
