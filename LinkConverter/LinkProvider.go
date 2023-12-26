package LinkConverter

import (
	"log"
	"net/http"
	"strings"

	"github.com/skip2/go-qrcode"
)

//  create an live end point and call this Function it will return url according to the device
func FlattradeAppRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("FlattradeAppRequest (+)")
	// user Agent will Get Request from Header  as string
	userAgent := r.UserAgent()
	//  Get Platform Method return url as per Device
	destinationURL := getPlatform(userAgent)

	// Write the destination URL as the response
	w.Write([]byte(destinationURL))

}

// Function to get the platform based on the user agent
func getPlatform(userAgent string) string {

	// userAgent will Get an string like Coming Request From Browser and with Device
	userAgent = strings.ToLower(userAgent)
	//  string.contains()   method is used Find the Partictuclar data present in large string  or not
	if strings.Contains(userAgent, "android") {
		return "https://play.google.com/store/apps/details?id=com.noren.ftconline"
	} else if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") {
		return "https://apps.apple.com/in/app/flattrade/id1543677587"
	} else if strings.Contains(userAgent, "Windows") {
		return "https://web.flattrade.in/#/"
	} else {
		return "Invalid Platform"
	}
}

//  if data comes in Query parameter we can use this
func LinkinQuery(w http.ResponseWriter, r *http.Request) {
	// Parse parameters from the request, you can use URL parameters or headers
	platform := r.URL.Query().Get("platform")
	// fmt.Println("platform", platform)
	// Determine the destination URL based on the platform
	var destinationURL string
	switch platform {
	case "android":
		destinationURL = "https://play.google.com/store/apps/details?id=com.noren.ftconline"
	case "ios":
		destinationURL = "https://apps.apple.com/in/app/flattrade/id1543677587"
	default:
		http.Error(w, "Invalid platform", http.StatusBadRequest)
		return
	}
	w.Write([]byte(destinationURL))
}

// Create an qr code using this method
func Qrcode(url string) {

	err1 := qrcode.WriteFile(url, qrcode.High, 256, "myFirstLink.png")
	if err1 != nil {
		log.Println("Sorry Couldnt able to print Qr code")
	}

}
