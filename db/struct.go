package db

type Dingdingbot struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

type Miner struct {
	MinerID string `json:"minerID"`
	Cname   string `json:"cname"`
	Size    int    `json:"size"`
}
