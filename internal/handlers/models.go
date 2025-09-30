package httphandlers

type shortenReq struct {
	URL string `json:"url" binding:"required"`
}

type shortenResp struct {
	ShortURL string `json:"short_url"`
}