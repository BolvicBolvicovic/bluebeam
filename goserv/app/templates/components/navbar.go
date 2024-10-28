package components

type Navbar struct {
	WhysButton 	Button
	SettingsButton 	Button
	ApisButton 	Button
	AnalyzerButton	Button
	LoginButton	Button
	LogoutButton	Button
	IsLoggedIn	bool
}

func NewNavbar(isLoggedIn bool) Navbar {
	return Navbar{
		WhysButton: 	Button{
			Text:		"why bluebeam",
			Link:		"/whyBluebeam",
			IsSubmit:	false,
			IsPrimary:	false,
		},
		SettingsButton: Button{
			Text:		"settings",
			Link:		"/settings",
			IsSubmit:	false,
			IsPrimary:	false,
		},
		AnalyzerButton: Button{
			Text:		"analyzer",
			Link:		"/analyzerPage",
			IsSubmit:	false,
			IsPrimary:	false,
		},
		ApisButton: 	Button{
			Text:		"api",
			Link:		"/apiPage",
			IsSubmit:	false,
			IsPrimary:	false,
		},
		LoginButton:	Button{
			Text:		"login",
			Link:		"/loginPage",
			IsSubmit:	false,
			IsPrimary:	false,
		},
		LogoutButton:	Button{
			Text:		"logout",
			Link:		"/logout",
			IsSubmit:	false,
			IsPrimary:	false,
		},
		IsLoggedIn:	isLoggedIn,
	}
}
