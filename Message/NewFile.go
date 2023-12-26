package message

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type respMessage struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppID           string               `json:"app_id"`
	UserID          string               `json:"user_id"`
	MessageID       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageURL         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]Attribute `json:"attributes"`
	Traits          map[string]Trait     `json:"traits"`
}

type Attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Trait struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

// http.HandleFunc("/message", message.ContactForm) -- Endpoint To call method
// 	http.ListenAndServe(":29091", nil) --- port no

func ContactForm(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	// log.Println("ContactForm (+)")
	var lResponse respMessage
	var lDynamic map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close() // Close the body when the function returns

	if err := decoder.Decode(&lDynamic); err != nil {
		fmt.Println("err", err)
	} else {
		Response := Worker(lDynamic)
		lResponse = Response
	}
	lData, lErr3 := json.Marshal(lResponse)
	if lErr3 != nil {
		// log.Println(lErr3)
		fmt.Fprintf(w, "Error found on marshalling Datas!")
		return
	} else {
		fmt.Fprintf(w, string(lData))
	}

}

func Worker(lDynamic interface{}) respMessage {

	// Create channels
	outputChannel := make(chan respMessage)
	var wg sync.WaitGroup
	// Increment the WaitGroup before starting the worker
	wg.Add(1)
	go convertedMsg(lDynamic, outputChannel, &wg)

	// Receive the converted message from the worker
	lconvertMessage := <-outputChannel
	// Wait for the worker to complete
	wg.Wait()
	// fmt.Println("lconvertMessage", lconvertMessage)
	return lconvertMessage
}

func convertedMsg(data interface{}, outputChannel chan<- respMessage, wg *sync.WaitGroup) {
	// log.Println("convertedMsg")
	defer wg.Done()
	var lRespDynamic respMessage
	var lvalue string

	switch jsonData := data.(type) {
	case map[string]interface{}:
		// Process the map keys and values
		for key, value := range jsonData {
			if str, ok := value.(string); ok {
				lvalue = str
			}

			switch key {
			case "ev":
				lRespDynamic.Event = lvalue
			case "t":
				lRespDynamic.PageTitle = lvalue
			case "et":
				lRespDynamic.EventType = lvalue
			case "id":
				lRespDynamic.AppID = lvalue
			case "uid":
				lRespDynamic.UserID = lvalue
			case "mid":
				lRespDynamic.MessageID = lvalue
			case "p":
				lRespDynamic.PageURL = lvalue
			case "l":
				lRespDynamic.ScreenSize = lvalue
			case "sc":
				lRespDynamic.ScreenSize = lvalue
			}
		}

		// Process dynamic attributes
		attributes := make(map[string]Attribute)

		// Iterate through the keys
		for i := 1; ; i++ {
			key := fmt.Sprintf("atrk%d", i)
			valueKey := fmt.Sprintf("atrv%d", i)
			typeKey := fmt.Sprintf("atrt%d", i)

			// Check if keys exist
			if attributeName, ok := jsonData[key].(string); ok {
				value := jsonData[valueKey].(string)
				attrType := jsonData[typeKey].(string)

				// Create Attribute struct
				attr := Attribute{
					Value: value,
					Type:  attrType,
				}

				// Add to attributes map
				attributes[attributeName] = attr
			} else {
				break // Break the loop if the key doesn't exist
			}
		}

		// Assign attributes to lRespDynamic
		lRespDynamic.Attributes = attributes

		// Process dynamic traits
		traits := make(map[string]Trait)

		// Iterate through the keys
		for i := 1; ; i++ {
			key := fmt.Sprintf("uatrk%d", i)
			valueKey := fmt.Sprintf("uatrv%d", i)
			typeKey := fmt.Sprintf("uatrt%d", i)

			// Check if keys exist
			if traitName, ok := jsonData[key].(string); ok {
				value := jsonData[valueKey].(string)
				traitType := jsonData[typeKey].(string)

				// Create Trait struct
				trait := Trait{
					Value: value,
					Type:  traitType,
				}

				// Add to traits map
				traits[traitName] = trait
			} else {
				break // Break the loop if the key doesn't exist
			}
		}

		// Assign traits to lRespDynamic
		lRespDynamic.Traits = traits

		// fmt.Println("lRespDynamic", lRespDynamic)
		// Print the result
		// fmt.Printf("%+v\n", lRespDynamic)

	default:
		fmt.Printf("Unsupported type: %T\n", jsonData)
	}

	outputChannel <- lRespDynamic
}
