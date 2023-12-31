package domain

type Interactive struct {
	BizId int64  `json:"bizId,omitempty"`
	Biz   string `json:"biz,omitempty"`

	ReadCount    int64 `json:"readCount,omitempty"`
	LikeCount    int64 `json:"likeCount,omitempty"`
	CollectCount int64 `json:"collectCount,omitempty"`
	Liked        bool  `json:"liked,omitempty"`
	Collected    bool  `json:"collected,omitempty"`
}
