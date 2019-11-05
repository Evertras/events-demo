package db

type DataFriendList struct {
	Friends []string `json:"friends"`
}

type DataSession struct {
	Host string `json:"host"`
}
