package noaaclient

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// type ProductDataValues struct {
// 	DataType DataProduct
// 	Metadata struct {
// 		ID   string `json:"id"`
// 		Name string `json:"name"`
// 		Lat  string `json:"lat"`
// 		Lon  string `json:"lon"`
// 	} `json:"metadata"`
// 	Data []struct {
// 		T string `json:"t"`
// 		V string `json:"v"`
// 		F string `json:"f"`
// 	} `json:"data"`
// }
