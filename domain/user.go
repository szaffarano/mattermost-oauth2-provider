package domain

// User is an oauth user
type User struct {
	AvatarURL        string      `json:"avatar_url"`
	Bio              interface{} `json:"bio"`
	CanCreateGroup   bool        `json:"can_create_group"`
	CanCreateProject bool        `json:"can_create_project"`
	ColorSchemeID    int64       `json:"color_scheme_id"`
	ConfirmedAt      string      `json:"confirmed_at"`
	CreatedAt        string      `json:"created_at"`
	CurrentSignInAt  string      `json:"current_sign_in_at"`
	Email            string      `json:"email"`
	External         bool        `json:"external"`
	ID               int64       `json:"id"`
	Identities       []struct {
		ExternUID string `json:"extern_uid"`
		Provider  string `json:"provider"`
	} `json:"identities"`
	LastActivityOn            string      `json:"last_activity_on"`
	LastSignInAt              string      `json:"last_sign_in_at"`
	Linkedin                  string      `json:"linkedin"`
	Location                  interface{} `json:"location"`
	Name                      string      `json:"name"`
	Organization              interface{} `json:"organization"`
	PrivateProfile            interface{} `json:"private_profile"`
	ProjectsLimit             int64       `json:"projects_limit"`
	PublicEmail               string      `json:"public_email"`
	SharedRunnersMinutesLimit int64       `json:"shared_runners_minutes_limit"`
	Skype                     string      `json:"skype"`
	State                     string      `json:"state"`
	ThemeID                   interface{} `json:"theme_id"`
	Twitter                   string      `json:"twitter"`
	TwoFactorEnabled          bool        `json:"two_factor_enabled"`
	Username                  string      `json:"username"`
	WebURL                    string      `json:"web_url"`
	WebsiteURL                string      `json:"website_url"`
}
