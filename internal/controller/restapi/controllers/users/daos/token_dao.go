package daos

type TokenDAO struct {
	Token string `json:"token"`
}

func ToTokenDAO(token string) *TokenDAO {
	return &TokenDAO{
		Token: token,
	}
}
