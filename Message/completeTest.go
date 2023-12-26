package message

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"sync"
// )

// // ========================Global Variable=========================
// const ErrorCode = "E"
// const SuccessCode = "S"

// // ========================================================

// type reqMessage struct {
// 	Ev     string `json:"ev"`
// 	Et     string `json:"et"`
// 	ID     string `json:"id"`
// 	UID    string `json:"uid"`
// 	MID    string `json:"mid"`
// 	T      string `json:"t"`
// 	P      string `json:"p"`
// 	L      string `json:"l"`
// 	SC     string `json:"sc"`
// 	ATRK1  string `json:"atrk1"`
// 	ATRV1  string `json:"atrv1"`
// 	ATRT1  string `json:"atrt1"`
// 	ATRK2  string `json:"atrk2"`
// 	ATRV2  string `json:"atrv2"`
// 	ATRT2  string `json:"atrt2"`
// 	UATRK1 string `json:"uatrk1"`
// 	UATRV1 string `json:"uatrv1"`
// 	UATRT1 string `json:"uatrt1"`
// 	UATRK2 string `json:"uatrk2"`
// 	UATRV2 string `json:"uatrv2"`
// 	UATRT2 string `json:"uatrt2"`
// 	UATRK3 string `json:"uatrk3"`
// 	UATRV3 string `json:"uatrv3"`
// 	UATRT3 string `json:"uatrt3"`
// }

// type respMessage struct {
// 	Event           string               `json:"event"`
// 	EventType       string               `json:"event_type"`
// 	AppID           string               `json:"app_id"`
// 	UserID          string               `json:"user_id"`
// 	MessageID       string               `json:"message_id"`
// 	PageTitle       string               `json:"page_title"`
// 	PageURL         string               `json:"page_url"`
// 	BrowserLanguage string               `json:"browser_language"`
// 	ScreenSize      string               `json:"screen_size"`
// 	Attributes      map[string]Attribute `json:"attributes"`
// 	Traits          map[string]Trait     `json:"traits"`
// 	Status          string               `json:"status"`
// 	ErrMsg          string               `json:"errMsg"`
// }

// type Attribute struct {
// 	Value string `json:"value"`
// 	Type  string `json:"type"`
// }

// type Trait struct {
// 	Value string `json:"value"`
// 	Type  string `json:"type"`
// }

// func ContactForm(w http.ResponseWriter, r *http.Request) {
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
// 	log.Println("ContactForm (+)")

// 	if r.Method == "POST" {
// 		var lReqMessage reqMessage
// 		var lRespMessage respMessage
// 		lRespMessage.Status = SuccessCode
// 		//  Read the Request From the Body
// 		lBody, lErr1 := ioutil.ReadAll(r.Body)
// 		if lErr1 != nil {
// 			lRespMessage.Status = ErrorCode
// 			lRespMessage.ErrMsg = lErr1.Error()
// 		} else {
// 			lErr2 := json.Unmarshal(lBody, &lReqMessage)
// 			log.Println("lReqMessage", lReqMessage)
// 			if lErr2 != nil {
// 				lRespMessage.Status = ErrorCode
// 				lRespMessage.ErrMsg = lErr2.Error()
// 			} else {
// 				lRespMessage = ChatBox(lReqMessage)
// 				log.Println("lRespMessage", lRespMessage)
// 			}
// 		}

// 		lData, lErr3 := json.Marshal(lRespMessage)
// 		if lErr3 != nil {
// 			log.Println(lErr3)
// 			fmt.Fprintf(w, "Error found on marshalling Datas!")
// 			return
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 	}

// 	log.Println("ContactForm (-)")
// }

// func ChatBox(lReqMessage reqMessage) respMessage {
// 	log.Println("ChatBox(+)")

// 	// Create channels
// 	inputChannel := make(chan reqMessage)
// 	outputChannel := make(chan respMessage)
// 	var wg sync.WaitGroup

// 	// Increment the WaitGroup before starting the worker
// 	wg.Add(1)
// 	go worker(inputChannel, outputChannel, &wg)

// 	// Send the message to the worker
// 	inputChannel <- lReqMessage

// 	// Receive the converted message from the worker
// 	lconvertMessage := <-outputChannel

// 	// Wait for the worker to complete
// 	wg.Wait()
// 	// Close the inputChannel to signal that no more data will be sent
// 	close(inputChannel)
// 	// Close the outputChannel after receiving the result
// 	close(outputChannel)

// 	log.Println("ChatBox(-)")
// 	return lconvertMessage
// }

// func worker(inputChannel chan reqMessage, outputChannel chan<- respMessage, wg *sync.WaitGroup) {
// 	log.Println("worker(+)")
// 	defer wg.Done()
// 	// Receive the message from the inputChannel
// 	msg := <-inputChannel

// 	// Simulate processing and creating a response
// 	convertedMsg := convertMessage(msg)
// 	// log.Println("Worker convertedMsg", convertedMsg)

// 	// Send the response to the outputChannel
// 	outputChannel <- convertedMsg
// 	log.Println("worker(-)")
// }

// func convertMessage(original reqMessage) respMessage {
// 	// log.Println("convertMessage(+)")
// 	//  Map All the values from orginal Message to Respones Message
// 	converted := respMessage{
// 		Event:           original.Ev,
// 		EventType:       original.Et,
// 		AppID:           original.ID,
// 		UserID:          original.UID,
// 		MessageID:       original.MID,
// 		PageTitle:       original.T,
// 		PageURL:         original.P,
// 		BrowserLanguage: original.L,
// 		ScreenSize:      original.SC,
// 		Attributes: map[string]Attribute{
// 			"form_varient": {Value: original.ATRV1, Type: original.ATRT1},
// 			"ref":          {Value: original.ATRV2, Type: original.ATRT2},
// 		},
// 		Traits: map[string]Trait{
// 			"name":  {Value: original.UATRV1, Type: original.UATRT1},
// 			"email": {Value: original.UATRV2, Type: original.UATRT2},
// 			"age":   {Value: original.UATRV3, Type: original.UATRT3},
// 		},
// 		Status: SuccessCode,
// 		ErrMsg: "",
// 	}

// 	log.Println("convertMessage(-)")
// 	// log.Println("converted", converted)

// 	return converted
// }
