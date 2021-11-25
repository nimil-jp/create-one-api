package request

type UserCreate struct {
	FirebaseUID string `json:"firebase_uid"`
	Email       string `json:"email"`
}

// profile

type UserSetCoverImage string

type UserEditProfile struct {
	AvatarImage  string `json:"avatar_image"`
	Name         string `json:"name"`
	About        string `json:"about"`
	Introduction string `json:"introduction"`

	Website   string `json:"website" validate:"omitempty,url"`
	Youtube   string `json:"youtube" validate:"omitempty,social_link"`
	Twitter   string `json:"twitter"`
	Facebook  string `json:"facebook"`
	Instagram string `json:"instagram"`
	Pinterest string `json:"pinterest"`
	Linkedin  string `json:"linkedin" validate:"omitempty,social_link"`
	Github    string `json:"github"`
	Qiita     string `json:"qiita"`
	Zenn      string `json:"zenn"`
}
