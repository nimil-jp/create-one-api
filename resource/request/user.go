package request

type UserCreate struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Username        string `json:"username"`
}

type UserLogin struct {
	Session  bool   `json:"session"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRefreshToken struct {
	Session      bool   `json:"session"`
	RefreshToken string `json:"refresh_token"`
}

type UserResetPasswordRequest struct {
	Email string `json:"email"`
}

type UserResetPassword struct {
	RecoveryToken   string `json:"recovery_token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// profile

type UserSetCoverImage string

type UserEditProfile struct {
	AvatarImage  string `json:"avatar_image"`
	Name         string `json:"name"`
	About        string `json:"about"`
	Introduction string `json:"introduction"`

	Website   string `json:"website"`
	Youtube   string `json:"youtube"`
	Twitter   string `json:"twitter"`
	Facebook  string `json:"facebook"`
	Instagram string `json:"instagram"`
	Pinterest string `json:"pinterest"`
	Linkedin  string `json:"linkedin"`
	Github    string `json:"github"`
	Qiita     string `json:"qiita"`
	Zenn      string `json:"zenn"`
}
