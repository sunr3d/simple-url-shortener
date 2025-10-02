package httphandlers

type shortenReq struct {
	URL string `json:"url" binding:"required"`
}

type shortenResp struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

type analyticsResp struct {
	Total   int64              `json:"total"`
	ByDay   []analyticsByDay   `json:"by_day,omitempty"`
	ByMonth []analyticsByMonth `json:"by_month,omitempty"`
	ByUA    []analyticsByUA    `json:"by_user_agent,omitempty"`
}

type analyticsByDay struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type analyticsByMonth struct {
	Year  int   `json:"year"`
	Month int   `json:"month"`
	Count int64 `json:"count"`
}

type analyticsByUA struct {
	UA    string `json:"ua"`
	Count int64  `json:"count"`
}
