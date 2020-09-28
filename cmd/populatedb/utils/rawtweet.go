package utils

// RawTweet interface generated with https://mholt.github.io/json-to-go/
type RawTweet struct {
	CreatedAt        string `json:"created_at"`
	ID               int64  `json:"id"`
	IDStr            string `json:"id_str"`
	FullText         string `json:"full_text"`
	Truncated        bool   `json:"truncated"`
	DisplayTextRange []int  `json:"display_text_range"`
	Entities         struct {
		Hashtags []struct {
			Text    string `json:"text"`
			Indices []int  `json:"indices"`
		} `json:"hashtags"`
		Symbols      []interface{} `json:"symbols"`
		UserMentions []struct {
			ScreenName string `json:"screen_name"`
			Name       string `json:"name"`
			ID         int64  `json:"id"`
			IDStr      string `json:"id_str"`
			Indices    []int  `json:"indices"`
		} `json:"user_mentions"`
		Urls []struct {
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
	User                 struct {
		ID          int64  `json:"id"`
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
				Urls []struct {
					URL         string `json:"url"`
					ExpandedURL string `json:"expanded_url"`
					DisplayURL  string `json:"display_url"`
					Indices     []int  `json:"indices"`
				} `json:"urls"`
			} `json:"description"`
		} `json:"entities"`
		Protected                      bool        `json:"protected"`
		FollowersCount                 int         `json:"followers_count"`
		FriendsCount                   int         `json:"friends_count"`
		ListedCount                    int         `json:"listed_count"`
		CreatedAt                      string      `json:"created_at"`
		FavouritesCount                int         `json:"favourites_count"`
		UtcOffset                      interface{} `json:"utc_offset"`
		TimeZone                       interface{} `json:"time_zone"`
		GeoEnabled                     bool        `json:"geo_enabled"`
		Verified                       bool        `json:"verified"`
		StatusesCount                  int         `json:"statuses_count"`
		Lang                           interface{} `json:"lang"`
		ContributorsEnabled            bool        `json:"contributors_enabled"`
		IsTranslator                   bool        `json:"is_translator"`
		IsTranslationEnabled           bool        `json:"is_translation_enabled"`
		ProfileBackgroundColor         string      `json:"profile_background_color"`
		ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
		ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
		ProfileBackgroundTile          bool        `json:"profile_background_tile"`
		ProfileImageURL                string      `json:"profile_image_url"`
		ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
		ProfileBannerURL               string      `json:"profile_banner_url"`
		ProfileImageExtensionsAltText  interface{} `json:"profile_image_extensions_alt_text"`
		ProfileBannerExtensionsAltText interface{} `json:"profile_banner_extensions_alt_text"`
		ProfileLinkColor               string      `json:"profile_link_color"`
		ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
		ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
		ProfileTextColor               string      `json:"profile_text_color"`
		ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
		HasExtendedProfile             bool        `json:"has_extended_profile"`
		DefaultProfile                 bool        `json:"default_profile"`
		DefaultProfileImage            bool        `json:"default_profile_image"`
		Following                      bool        `json:"following"`
		FollowRequestSent              bool        `json:"follow_request_sent"`
		Notifications                  bool        `json:"notifications"`
		TranslatorType                 string      `json:"translator_type"`
	} `json:"user"`
	Geo               *Geo         `json:"geo"`
	Coordinates       *Coordinates `json:"coordinates"`
	Place             *Place       `json:"place"`
	Contributors      interface{}  `json:"contributors"`
	IsQuoteStatus     bool         `json:"is_quote_status"`
	FavoriteCount     int64        `json:"favorite_count"`
	Favorited         bool         `json:"favorited"`
	Retweeted         bool         `json:"retweeted"`
	RetweetedStatus   *RawTweet    `json:"retweeted_status"`
	RetweetCount      int64        `json:"retweet_count"`
	PossiblySensitive bool         `json:"possibly_sensitive"`
	Lang              string       `json:"lang"`
}

type Place struct {
	ID              string        `json:"id"`
	URL             string        `json:"url"`
	PlaceType       string        `json:"place_type"`
	Name            string        `json:"name"`
	FullName        string        `json:"full_name"`
	CountryCode     string        `json:"country_code"`
	Country         string        `json:"country"`
	ContainedWithin []interface{} `json:"contained_within"`
	BoundingBox     struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"bounding_box"`
	Attributes struct {
	} `json:"attributes"`
}

type Geo struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Coordinates struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
