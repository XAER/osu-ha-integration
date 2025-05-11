package domain

type OsuUser struct {
	Username    string  `json:"username"`
	GlobalRank  int     `json:"global_rank"`
	CountryRank int     `json:"country_rank"`
	PP          float64 `json:"pp"`
	Accuracy    float64 `json:"accuracy"`
	PlayCount   int     `json:"play_count"`
}
