package components

type Navbar struct {
	WhysButton 	Button
	DashboardButton	Button
	ApisButton 	Button
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
		DashboardButton: Button{
			Text:		"dashboard",
			Link:		"/dashboard",
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
