package httphandlers

type shortenReq struct {
	URL string `json:"url" binding:"required"`
}

type shortenResp struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}
