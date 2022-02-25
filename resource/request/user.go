package request

type UserCreate struct {
	FirebaseUID string `json:"firebase_uid"`
	Email       string `json:"email"`
}

// profile

type UserEditRequest struct {
	UnitPrice uint `json:"unit_price" validate:"min=200,max=10000"`

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

type UserPatchRequest struct {
	Username string `json:"username" validate:"omitempty,username"`

	UnitPrice uint `json:"unit_price" validate:"omitempty,min=200,max=10000"`

	CoverImage *string `json:"cover_image"`

	AvatarImage  *string `json:"avatar_image"`
	Name         *string `json:"name"`
	About        *string `json:"about"`
	Introduction *string `json:"introduction"`

	Website   *string `json:"website" validate:"omitempty,url"`
	Youtube   *string `json:"youtube" validate:"omitempty,social_link"`
	Twitter   *string `json:"twitter"`
	Facebook  *string `json:"facebook"`
	Instagram *string `json:"instagram"`
	Pinterest *string `json:"pinterest"`
	Linkedin  *string `json:"linkedin" validate:"omitempty,social_link"`
	Github    *string `json:"github"`
	Qiita     *string `json:"qiita"`
	Zenn      *string `json:"zenn"`
}
