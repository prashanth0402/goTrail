package Authentication

// var (
// 	googleOauthConfig = &oauth2.Config{
// 		RedirectURL:  "http://localhost:29091/googlecallback",
// 		ClientID:     "YOUR_CLIENT_ID",
// 		ClientSecret: "YOUR_CLIENT_SECRET",
// 		Scopes:       []string{"profile", "email"},
// 		Endpoint:     google.Endpoint,
// 	}
// )

// type UserInfo struct {
// 	Email   string `json:"email"`
// 	Name    string `json:"name"`
// 	Picture string `json:"picture"`
// }

// func getUserInfo(token *oauth2.Token) (*UserInfo, error) {
// 	client := googleOauthConfig.Client(oauth2.NoContext, token)
// 	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("Failed to retrieve user information. Status code: %d", resp.StatusCode)
// 	}

// 	userInfo := &UserInfo{}
// 	if err := json.NewDecoder(resp.Body).Decode(userInfo); err != nil {
// 		return nil, err
// 	}

// 	return userInfo, nil
// }

// func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")

// 	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
// 	if err != nil {
// 		log.Printf("Failed to exchange code for token: %v", err)
// 		http.Error(w, "Failed to exchange code for token", http.StatusInternalServeError)
// 		return
// 	}

// 	userInfo, err := getUserInfo(token)
// 	if err != nil {
// 		log.Printf("Failed to get user info: %v", err)
// 		http.Error(w, "Failed to get user info", http.StatusInternalServeError)
// 		return
// 	}

// 	// You can now use the 'userInfo' struct to access user information.
// 	// userInfo.Email, userInfo.Name, userInfo.Picture, etc.
// 	fmt.Printf("User Email: %s\n", userInfo.Email)
// 	fmt.Printf("User Name: %s\n", userInfo.Name)
// 	fmt.Printf("User Picture URL: %s\n", userInfo.Picture)

// 	// You may want to store this information in your application's user database.

// 	// Redirect or render your website as appropriate.
// }
