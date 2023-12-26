package OTP

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pquerna/otp/totp"
	qrcode "github.com/skip2/go-qrcode"
)

func Otp() {
	log.Println("otp (+)")
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "alice@example.com",
	})

	if err != nil {
		log.Fatal(err)
	}

	// Convert TOTP key into a QR code encoded as a PNG image.
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(&buf, img)
	if err != nil {
		log.Fatal(err)
	}
	// Display the QR code to the user (not provided in this example).
	displayQrCode(buf.Bytes())
	// You should have your own implementation to show the QR code.

	// Now validate that the user has successfully added the passcode.
	passcode := promptForPasscode()
	log.Println("passcode", passcode)

	valid := totp.Validate(passcode, key.Secret())

	if valid {
		// User successfully used their TOTP, save it to your backend!
		storeSecret("alice@example.com", key.Secret())
	}

	log.Println("otp (-)")
}

// func VerifyOtp() {
// 	log.Println("VerifyOtp(+)")
// 	passcode := promptForPasscode()
// 	secret := getSecret("alice@example.com")

// 	valid := totp.Validate(passcode, secret)

// 	if valid {
// 		log.Println("Success To Login")
// 		// Success! continue login process.
// 	}
// 	log.Println("VerifyOtp(-)")

// }

func Qrcode() {
	content := "https://www.example.com" // Replace with your desired content

	// Generate a QR code
	code, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

	// Print the QR code as a string (for demonstration purposes)
	fmt.Println(string(code))

	err1 := qrcode.WriteFile("http://onelink.to/amoflattrade", qrcode.Medium, 256, "myFirstQr.png")
	if err1 != nil {
		log.Println("Sorry Couldnt able to print Qr code")
	}

	// You can save the QR code to a file or display it as an image in your application
}

// }

// Replace these functions with your actual implementations.
func displayQrCode(data []byte) {
	// Your implementation to display the QR code.
	// fmt.Printf("QR code image data: %v\n", data)
	err1 := qrcode.WriteFile(string(data), qrcode.Low, 256, "myQrCode.png")
	if err1 != nil {
		log.Println("Sorry Couldnt able to print Qr code")
	}
}

func promptForPasscode() string {
	// Your implementation to get the passcode from the user.
	var passcode string

	fmt.Print("Enter Passcode: ")
	_, err := fmt.Scanln(&passcode)
	if err != nil {
		log.Fatal(err)
	}

	// Trim leading and trailing whitespace
	passcode = strings.TrimSpace(passcode)

	return passcode // Replace with user input.
}

func storeSecret(accountName string, secretKey string) error {
	// Your implementation to store the secret key.
	// Write the secret to a file
	err := ioutil.WriteFile(accountName, []byte(secretKey), 0600)
	if err != nil {
		return err
	}

	return nil
}
