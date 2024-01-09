package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/skip2/go-qrcode"
	// "strings"
	// "os"
)

func autoRestart() {
	for {
		now := time.Now()
		//resart the program everyday at 4am
		//at 3am, the program goes for 1 hour sleep and after that it will restart
		if now.Hour() == 3 {
			//sleep for an hour so that the hour changes to 4 and this condition
			//in the loop does not  continue again in next iteration
			time.Sleep(60 * 61 * time.Second)
			fmt.Println(now.Hour(), now.Minute(), now.Second())
			log.Println(now.Hour(), now.Minute(), now.Second())
			// Restart the program
			fmt.Println("Restarting the program...")
			log.Println("Restarting the program...")
			execPath, err := os.Executable()
			if err != nil {
				fmt.Println("Error getting executable path:", err)
				log.Println("Error getting executable path:", err)

				return
			}
			cmd := exec.Command(execPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Start()
			if err != nil {
				fmt.Println("Error restarting program:", err)
				log.Println("Error restarting program:", err)
				return
			}
			os.Exit(0)

		}
		time.Sleep(60 * 30 * time.Second)
	}
}

func main() {
	log.Println("Server Started...")
	f, err := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	go autoRestart()
	log.Println("Server Started (+)")
	// Run()
	// http.HandleFunc("/Handle", LinkConverter.FlattradeAppRequest)
	// newemailpackage.SendMail()

}

func Qrcode(url string) {

	err1 := qrcode.WriteFile(url, qrcode.High, 256, "myFirstLink.png")
	if err1 != nil {
		log.Println("Sorry Couldnt able to print Qr code")
	}

}

// Method with pointer receiver that calculates the area

// func HTTP() {
// 	log.Println("HTTP (+)")
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello Domain")
// 	}

// 	http.HandleFunc("/Domain", handler)
// 	log.Println("HTTP (-)")

// 	http.ListenAndServe(":29091", nil)

// }

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"

// 	qrcode "github.com/skip2/go-qrcode"
// )

// func main() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		// Get the host and port from the request
// 		host := r.Host
// 		// You can customize the path as needed
// 		path := "/your/path"

// 		// Get the user agent from the request
// 		userAgent := r.UserAgent()

// 		// Determine the app store link based on the user agent
// 		appStoreLink := getPlatform(userAgent)
// 		fmt.Println("appStoreLink", appStoreLink)
// 		// Combine host, path, and app store link to create the dynamic link
// 		dynamicLink := fmt.Sprintf("http://%s%s?link=%s", host, path, appStoreLink)

// 		// Generate a QR code for the dynamic link
// 		err := qrcode.WriteFile(dynamicLink, qrcode.Medium, 256, "qrcode.png")
// 		if err != nil {
// 			http.Error(w, "Error generating QR code", http.StatusInternalServerError)
// 			return
// 		}

// 		// Serve the HTML page with the QR code
// 		// renderQRCodePage(w, dynamicLink)
// 	})

// 	port := 8080
// 	fmt.Printf("Server running on :%d\n", port)
// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
// }

// Function to get the platform based on the user agent
// func getPlatform(userAgent string) string {

// 	userAgent = strings.ToLower(userAgent)
// 	fmt.Println("userAgent", userAgent)
// 	if strings.Contains(userAgent, "android") {
// 		return "Android"
// 	} else if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") {
// 		return "iOS"
// 	} else if strings.Contains(userAgent, "postmanruntime") {
// 		return "postmanruntime"
// 	} else {
// 		return "unknown"
// 	}
// }
