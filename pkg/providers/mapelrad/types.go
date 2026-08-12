package mapelrad

type coin string
type currency string
type channel string
type blockchain string

type errorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
