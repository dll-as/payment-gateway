package auth

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type loginResponse struct {
// 	Token string `json:"token"`
// }

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
