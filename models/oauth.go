package models

// OAuth struct
type OAuth struct {
	ID          uint   `json:"id"`
	Service     string `json:"service"`
	OpenID      string `json:"open_id"`
	UserID      uint   `json:"user_id"`
	Token       string `json:"token"`
	TokenSecret string `json:"token_secret"`
}

// GoogleOAuthUserInfo struct							|
// endpoint: https://www.googleapis.com/oauth2/v2/userinfo
type GoogleOAuthUserInfo struct {
	OpenID        string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"local"`
}

// FacebookOAuthUserInfo struct							|
// endpoint: https://graph.facebook.com/v3.2/me?fields=id,name,picture,email,birthday&access_token=
type FacebookOAuthUserInfo struct {
	OpenID   string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	BirthDay string `json:"birthday"`
	Picture  struct {
		Data struct {
			Height       uint   `json:"height"`
			Width        uint   `json:"width"`
			IsSilhouette bool   `json:"is_silhouette"`
			URL          string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

// TwitterOAuthUserInfo struct							|
// endpoint: https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true
type TwitterOAuthUserInfo struct {
	ID          int    `json:"id"`
	IDStr       string `json:"id_str"`
	Name        string `json:"name"`
	ScreenName  string `json:"screen_name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Entities    struct {
		URL struct {
			Urls []struct {
				URL         string `json:"url"`
				ExpandedURL string `json:"expanded_url"`
				DisplayURL  string `json:"display_url"`
				Indices     []int  `json:"indices"`
			} `json:"urls"`
		} `json:"url"`
		Description struct {
			Urls []interface{} `json:"urls"`
		} `json:"description"`
	} `json:"entities"`
	Protected       bool        `json:"protected"`
	FollowersCount  int         `json:"followers_count"`
	FriendsCount    int         `json:"friends_count"`
	ListedCount     int         `json:"listed_count"`
	CreatedAt       string      `json:"created_at"`
	FavouritesCount int         `json:"favourites_count"`
	UtcOffset       interface{} `json:"utc_offset"`
	TimeZone        interface{} `json:"time_zone"`
	GeoEnabled      bool        `json:"geo_enabled"`
	Verified        bool        `json:"verified"`
	StatusesCount   int         `json:"statuses_count"`
	Lang            interface{} `json:"lang"`
	Status          struct {
		CreatedAt string `json:"created_at"`
		ID        int64  `json:"id"`
		IDStr     string `json:"id_str"`
		Text      string `json:"text"`
		Truncated bool   `json:"truncated"`
		Entities  struct {
			Hashtags     []interface{} `json:"hashtags"`
			Symbols      []interface{} `json:"symbols"`
			UserMentions []interface{} `json:"user_mentions"`
			Urls         []struct {
				URL         string `json:"url"`
				ExpandedURL string `json:"expanded_url"`
				DisplayURL  string `json:"display_url"`
				Indices     []int  `json:"indices"`
			} `json:"urls"`
		} `json:"entities"`
		Source               string      `json:"source"`
		InReplyToStatusID    interface{} `json:"in_reply_to_status_id"`
		InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str"`
		InReplyToUserID      interface{} `json:"in_reply_to_user_id"`
		InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str"`
		InReplyToScreenName  interface{} `json:"in_reply_to_screen_name"`
		Geo                  interface{} `json:"geo"`
		Coordinates          interface{} `json:"coordinates"`
		Place                interface{} `json:"place"`
		Contributors         interface{} `json:"contributors"`
		IsQuoteStatus        bool        `json:"is_quote_status"`
		RetweetCount         int         `json:"retweet_count"`
		FavoriteCount        int         `json:"favorite_count"`
		Favorited            bool        `json:"favorited"`
		Retweeted            bool        `json:"retweeted"`
		PossiblySensitive    bool        `json:"possibly_sensitive"`
		Lang                 string      `json:"lang"`
	} `json:"status"`
	ContributorsEnabled            bool   `json:"contributors_enabled"`
	IsTranslator                   bool   `json:"is_translator"`
	IsTranslationEnabled           bool   `json:"is_translation_enabled"`
	ProfileBackgroundColor         string `json:"profile_background_color"`
	ProfileBackgroundImageURL      string `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHTTPS string `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool   `json:"profile_background_tile"`
	ProfileImageURL                string `json:"profile_image_url"`
	ProfileImageURLHTTPS           string `json:"profile_image_url_https"`
	ProfileBannerURL               string `json:"profile_banner_url"`
	ProfileLinkColor               string `json:"profile_link_color"`
	ProfileSidebarBorderColor      string `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool   `json:"profile_use_background_image"`
	HasExtendedProfile             bool   `json:"has_extended_profile"`
	DefaultProfile                 bool   `json:"default_profile"`
	DefaultProfileImage            bool   `json:"default_profile_image"`
	CanMediaTag                    bool   `json:"can_media_tag"`
	FollowedBy                     bool   `json:"followed_by"`
	Following                      bool   `json:"following"`
	FollowRequestSent              bool   `json:"follow_request_sent"`
	Notifications                  bool   `json:"notifications"`
	TranslatorType                 string `json:"translator_type"`
	Suspended                      bool   `json:"suspended"`
	NeedsPhoneVerification         bool   `json:"needs_phone_verification"`
	Email                          string `json:"email"`
}
