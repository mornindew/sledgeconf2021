package noaaclient

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
