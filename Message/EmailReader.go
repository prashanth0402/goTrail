package message

import (
	"database/sql"
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type Bond struct {
	SecurityCategory   string `json:"securityCategory"`
	CoupType           string `json:"coupType"`
	AuctionDescription string `json:"auctionDescription"`
	ISIN               string `json:"iSIN"`
	MaturityDate       string `json:"maturityDate"`
	BiddingStartTime   string `json:"biddingStartTime"`
	BiddingEndTime     string `json:"biddingEndTime"`
	AllotmentDate      string `json:"allotmentDate"`
	SettlementDate     string `json:"settlementDate"`
	UPIEndTime         string `json:"uPIEndTime"`
	NetBankingEndTime  string `json:"netBankingEndTime"`
	NACHEndTime        string `json:"nACHEndTime"`
	IndicativeYield    string `json:"indicativeYield"`
}

type GmailRespBond struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func GmailRetriever(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		log.Println("GmailRetriever (+)")
		lConfigFile := common.ReadTomlConfig("toml/NovoRbiMail.toml")
		loginEmail := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["loginEmail"])
		KeyWords := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["KeyWords"])
		loginAppKey := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["loginAppKey"])
		From := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["From"])
		var lGmailResp GmailRespBond
		lGmailResp.Status = common.SuccessCode
		// Connect to the mail server
		c, lErr1 := connectToMailServer()
		if lErr1 != nil {
			log.Println("GMGM01", lErr1)
			lGmailResp.Status = common.ErrorCode
			lGmailResp.ErrMsg = lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("GMGM01", "Unable To Connect an Mail Server"))
			return
		} else {
			defer c.Logout()
			lErr2 := authenticate(c, loginEmail, loginAppKey)
			if lErr2 != nil {
				log.Println("GMGM02", lErr2)
				lGmailResp.Status = common.ErrorCode
				lGmailResp.ErrMsg = lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("GMGM02", "Error on Login UserId and Password Authenticate"))
				return
			} else {
				// Select INBOX
				lErr3 := selectInbox(c)
				if lErr3 != nil {
					log.Println("GMGM03", lErr3)
					lGmailResp.Status = common.ErrorCode
					lGmailResp.ErrMsg = lErr3.Error()
					fmt.Fprintf(w, helpers.GetErrorString("GMGM03", "Error on Selecting Inbox"))
					return
				} else {
					IndividualKeyword := strings.Split(KeyWords, ",")
					uids, lErr4 := searchForEmails(c, From, IndividualKeyword)
					if lErr4 != nil {
						log.Println("GMGM04", lErr4)
						lGmailResp.Status = common.ErrorCode
						lGmailResp.ErrMsg = lErr4.Error()
						fmt.Fprintf(w, helpers.GetErrorString("GMGM04", "Error on Search for Emails"))
						return
					} else if len(uids) != 0 {
						// Find the most recent email
						// recentEmailUID, lErr5 := fetchRecentEmailUID(c, uids)
						recentEmailUID, lErr5 := fetchCurrentDayEmailUIDs(c, uids)
						if lErr5 != nil {
							log.Println("GMGM05", lErr5)
							lGmailResp.Status = common.ErrorCode
							lGmailResp.ErrMsg = lErr5.Error()
							fmt.Fprintf(w, helpers.GetErrorString("GMGM05", "Error on Fetching Recent Email Records"))
							return
						} else {
							if recentEmailUID != nil {

								messages, done, lErr6 := fetchEmailBody(c, recentEmailUID)
								if lErr6 != nil {
									log.Println("GMGM06", lErr6)
									lGmailResp.Status = common.ErrorCode
									lGmailResp.ErrMsg = lErr6.Error()
									fmt.Fprintf(w, helpers.GetErrorString("GMGM06", "Unable To Fetch Email Body"))
									return
								} else {
									// Process email messages
									bondDataCh := make(chan Bond)
									go func() {
										lErr7 := processEmailMessages(messages, &imap.BodySectionName{}, bondDataCh)
										if lErr7 != nil {
											log.Println("GMGM07", lErr7)
											lGmailResp.Status = common.ErrorCode
											lGmailResp.ErrMsg = lErr7.Error()
											fmt.Fprintf(w, helpers.GetErrorString("GMGM06", "Unable To Fetch Email Body"))
											return
										}
									}()

									go func() {
										for bond := range bondDataCh {
											log.Println("bond Data Received")

											// Print all fields for debugging purposes
											log.Printf("Security Category: %s\n", bond.SecurityCategory)
											log.Printf("Coup Type: %s\n", bond.CoupType)
											log.Printf("Auction Description: %s\n", bond.AuctionDescription)
											log.Printf("ISIN: %s\n", bond.ISIN)
											log.Printf("Maturity Date: %s\n", bond.MaturityDate)
											log.Printf("Bidding Start Time: %s\n", bond.BiddingStartTime)
											log.Printf("Bidding End Time: %s\n", bond.BiddingEndTime)
											log.Printf("Allotment Date: %s\n", bond.AllotmentDate)
											log.Printf("Settlement Date: %s\n", bond.SettlementDate)
											log.Printf("UPI End Time: %s\n", bond.UPIEndTime)
											log.Printf("Net Banking End Time: %s\n", bond.NetBankingEndTime)
											log.Printf("NACH End Time: %s\n", bond.NACHEndTime)
											log.Printf("Indicative Yield: %s\n", bond.IndicativeYield)
											log.Println("-----------------------------------------------------------")

											// Print the entire 'bond' struct
											log.Printf("bond: %+v\n", bond)

											lErr8 := setRbiMail(bond)
											if lErr8 != nil {
												log.Println("GMGM08", lErr8)
												lGmailResp.Status = common.ErrorCode
												lGmailResp.ErrMsg = lErr8.Error()
												fmt.Fprintf(w, helpers.GetErrorString("GMGM08", "Unable To Insert and Update an Email"))
												return
											}

											// Add your logic here for further processing or database insertion
										}

									}()

								}
								<-done
							} else {
								lGmailResp.ErrMsg = "No Emails on Current Day to Read"
							}
						}

					} else {
						lGmailResp.ErrMsg = "No unread Emails"
					}

				}

			}

		}

		lData, lErr := json.Marshal(lGmailResp)
		if lErr != nil {
			log.Println("VGAV05", lErr)
			fmt.Fprintf(w, helpers.GetErrorString("VGAV05", "Error Occur in Marshalling Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}

	}
}

func connectToMailServer() (*client.Client, error) {
	log.Println("Connecting to mail server (+)")
	c, lErr1 := client.DialTLS("imap.gmail.com:993", nil)
	if lErr1 != nil {
		log.Println("Error connecting to mail server:", lErr1)
		return nil, lErr1
	}
	log.Println("Connected to mail server (-)")
	return c, nil
}

func authenticate(c *client.Client, email, password string) error {
	log.Println("Authenticating (+)")
	if err := c.Login(email, password); err != nil {
		log.Fatal("Error authenticating:", err)
		return err
	}
	log.Println("Authenticated (-)")
	return nil
}

func selectInbox(c *client.Client) error {
	log.Println("Selecting INBOX (+)")
	_, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal("Error selecting INBOX:", err)
		return err
	}
	log.Println("Selected INBOX (-)")
	return nil
}

func searchForEmails(c *client.Client, from string, keywords []string) ([]uint32, error) {
	log.Println("Searching for emails (+)")
	criteria := imap.NewSearchCriteria()
	criteria.Header = textproto.MIMEHeader{"From": {from}}
	criteria.Text = append(criteria.Text, keywords...)
	uids, err := c.Search(criteria)
	if err != nil {
		log.Fatal("Error searching for emails:", err)
		return nil, err
	} else {
		if len(uids) == 0 {
			log.Println("No unread mails")
			return nil, nil
		}

	}
	log.Println("Searching for emails (-)")
	log.Println("Found emails:", uids)
	return uids, nil
}

func fetchCurrentDayEmailUIDs(c *client.Client, uids []uint32) ([]uint32, error) {
	log.Println("Fetching current day email UIDs...")

	var currentDayEmailUIDs []uint32
	currentDay := time.Now() // Truncate to the beginning of the current day

	for _, uid := range uids {
		seqset := new(imap.SeqSet)
		seqset.AddNum(uid)
		items := []imap.FetchItem{imap.FetchEnvelope}

		messages := make(chan *imap.Message, 1)
		done := make(chan struct{})

		go func() {
			defer close(done)
			if err := c.Fetch(seqset, items, messages); err != nil {
				log.Println("Error fetching email:", err)
			}
		}()

		msg := <-messages
		if msg == nil {
			log.Println("Error fetching email: nil message")
			continue
		}

		emailDate := msg.Envelope.Date
		// emailDate = emailDate
		if emailDate.Year() == currentDay.Year() && emailDate.Month() == currentDay.Month() && emailDate.Day() == currentDay.Day() {
			currentDayEmailUIDs = append(currentDayEmailUIDs, uid)
		}
	}

	return currentDayEmailUIDs, nil
}

// func fetchRecentEmailUID(c *client.Client, uids []uint32) ([]uint32, error) {
// 	log.Println("Fetching recent email UID...")
// 	var recentEmailUID uint32
// 	var recentEmailDate time.Time

// 	for _, uid := range uids {
// 		seqset := new(imap.SeqSet)
// 		seqset.AddNum(uid)
// 		items := []imap.FetchItem{imap.FetchEnvelope}

// 		messages := make(chan *imap.Message, 1)
// 		done := make(chan struct{})
// 		var lErr1 error

// 		go func() {
// 			defer close(done)
// 			lErr1 = c.Fetch(seqset, items, messages)
// 		}()
// 		if lErr1 != nil {
// 			return recentEmailUID, lErr1
// 		}
// 		msg := <-messages
// 		emailDate := msg.Envelope.Date

// 		if emailDate.After(recentEmailDate) {
// 			recentEmailDate = emailDate
// 			recentEmailUID = uid
// 		}
// 	}

// 	log.Println("Recent email UID:", recentEmailUID)
// 	return recentEmailUID, nil
// }

func fetchEmailBody(c *client.Client, recentEmailUID []uint32) (chan *imap.Message, chan struct{}, error) {
	log.Println("Fetching email body...")
	seqset := new(imap.SeqSet)
	seqset.AddNum(recentEmailUID...)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}
	messages := make(chan *imap.Message, 10)
	done := make(chan struct{})

	go func() {
		defer close(done)
		if err := c.Fetch(seqset, items, messages); err != nil {
			log.Println("GMFE01 Error fetching email body:", err)
		}
	}()

	return messages, done, nil
}

func processEmailMessages(messages chan *imap.Message, section *imap.BodySectionName, bondDataCh chan<- Bond) error {
	log.Println("Processing email messages...")
	for msg := range messages {
		htmlcontent := msg.GetBody(section)
		if htmlcontent == nil {
			log.Println("Server didn't return message body r ==", htmlcontent)
			break
		} else {
			mailReader, lErr1 := mail.CreateReader(htmlcontent)
			if lErr1 != nil {
				log.Println("GPEM01 Error creating mail reader:", lErr1)
				return lErr1
			} else {
				header := mailReader.Header

				_, lErr2 := header.Date()
				if lErr2 != nil {
					log.Println("GPEM02", lErr2)
					return lErr2
				} else {
					_, lErr3 := header.AddressList("From")
					if lErr3 != nil {
						log.Println("GPEM03", lErr3)
						return lErr3

					} else {
						_, lErr4 := header.Subject()
						if lErr4 != nil {
							log.Println("GPEM04", lErr4)
							return lErr4
						} else {
							htmlBody, lErr6 := io.ReadAll(htmlcontent)
							if lErr6 != nil {
								log.Println("GPEM06 Error reading HTML body:", lErr6)
								return lErr6
							} else {
								// log.Println("htmlBody", string(htmlBody))

								lErr7 := processHTML(string(htmlBody), bondDataCh)
								if lErr7 != nil {
									log.Println("GPEM07 processHTML", lErr7)
									return lErr7
								}
							}

						}

						// processHTML function is used to extract required content from the body in html
					}
				}
			}
		}
	}

	// Close the channel outside the loop
	close(bondDataCh)
	return nil
}

func processHTML(html string, ch chan<- Bond) error {
	log.Println("processHTML (+)")
	// Implement your logic to extract information from the HTML part
	doc, lErr1 := goquery.NewDocumentFromReader(strings.NewReader(html))
	if lErr1 != nil {
		log.Println("GMPH01", lErr1)
		return lErr1
	}
	// Navigate through the HTML structure to find the table
	tableSelector := "table"
	table := doc.Find(tableSelector)

	// Extract information using goquery selectors
	table.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		// Implement your logic to extract information from each row of the table

		var bond Bond

		// Iterate over td elements in the current row
		s.Find("td").Each(func(index int, td *goquery.Selection) {
			switch index {
			case 0:
				bond.SecurityCategory = td.Text()
			case 1:
				bond.CoupType = td.Text()
			case 2:
				bond.AuctionDescription = td.Text()
			case 3:
				bond.ISIN = td.Text()
			case 4:
				bond.MaturityDate = td.Text()
			case 5:
				bond.BiddingStartTime = td.Text()
			case 6:
				bond.BiddingEndTime = td.Text()
			case 7:
				bond.AllotmentDate = td.Text()
			case 8:
				bond.SettlementDate = td.Text()
			case 9:
				bond.UPIEndTime = td.Text()
			case 10:
				bond.NetBankingEndTime = td.Text()
			case 11:
				bond.NACHEndTime = td.Text()
			case 12:
				bond.IndicativeYield = td.Text()
			}
		})

		ch <- bond

	})

	log.Println("processHTML (-)")
	return nil
}

func setRbiMail(pbond Bond) error {
	log.Println("setRbiMail(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GRSRI01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		lFlag, lErr2 := checkMailExist(pbond, lDb)
		if lErr2 != nil {
			log.Println("GRSRI02", lErr2)
			return lErr2

		} else {
			if lFlag == "Y" {
				lErr3 := InsertRBiMail(pbond, lDb)
				if lErr3 != nil {
					log.Println("GRSRI03", lErr3)
					return lErr3
				}

			} else {
				lErr4 := updateRbiMail(pbond, lDb)
				if lErr4 != nil {
					log.Println("GRSRI04", lErr4)
					return lErr4
				}

			}
		}
	}

	log.Println("setRbiMail(-)")
	return nil

}

func InsertRBiMail(pbond Bond, lDb *sql.DB) error {
	log.Println("InsertRBiMail(+)")

	lCoreString := `INSERT INTO novo_rbi_mail
	( SecurityCategory, coupType, ISIN, auctionDescription, MaturityDate, allotmentDate, SettlementDate, biddingStartTime, biddingEndTime, UPIEndTime, NetBankingEndTime, NACHEndTime, IndicativeYield, createdBy, createdDate,updatedBy,updatedDate)
	VALUES(?, ?,?, ?, STR_TO_DATE(?, '%d-%b-%Y'),STR_TO_DATE(?, '%d-%b-%Y'),STR_TO_DATE(?, '%d-%b-%Y'), STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'),STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'),STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'), STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'), STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'), ?, ?, now(), ?, now());`

	_, lErr2 := lDb.Exec(lCoreString, pbond.SecurityCategory, pbond.CoupType, pbond.ISIN, pbond.AuctionDescription, pbond.MaturityDate, pbond.AllotmentDate, pbond.SettlementDate, pbond.BiddingStartTime, pbond.BiddingEndTime, pbond.UPIEndTime, pbond.NetBankingEndTime, pbond.NACHEndTime, pbond.IndicativeYield, common.AUTOBOT, common.AUTOBOT)

	if lErr2 != nil {
		log.Println("GMIRM02", lErr2)
		return lErr2
	}
	log.Println("InsertRBiMail(-)")
	return nil

}

func updateRbiMail(pbond Bond, lDb *sql.DB) error {
	log.Println("updateRbiMail(+)")

	lCoreString := `UPDATE novo_rbi_mail
	SET  coupType=?, auctionDescription=?, MaturityDate= STR_TO_DATE(?, '%d-%b-%Y'), 
	allotmentDate= STR_TO_DATE(?, '%d-%b-%Y'), SettlementDate= STR_TO_DATE(?, '%d-%b-%Y'), biddingStartTime= STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'), 
	biddingEndTime=STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'), UPIEndTime=STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'), NetBankingEndTime=STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'), 
	NACHEndTime=STR_TO_DATE(?, '%d-%b-%Y %h:%i:%s %p'), IndicativeYield=?, updatedDate= now(), updatedBy=?
	WHERE SecurityCategory=? and ISIN=?; `
	_, lErr1 := lDb.Exec(lCoreString, pbond.CoupType, pbond.AuctionDescription, pbond.MaturityDate, pbond.AllotmentDate, pbond.SettlementDate, pbond.BiddingStartTime, pbond.BiddingEndTime, pbond.UPIEndTime, pbond.NetBankingEndTime, pbond.NACHEndTime, pbond.IndicativeYield, common.AUTOBOT, pbond.SecurityCategory, pbond.ISIN)
	if lErr1 != nil {
		log.Println("GMURM01", lErr1)
		return lErr1
	}
	log.Println("updateRbiMail(-)")

	return nil

}

func checkMailExist(pBond Bond, lDb *sql.DB) (string, error) {
	log.Println("checkMailExist (+)")
	var lFlag string
	Count := 0

	lCoreString := `select count(SecurityCategory)
		from novo_rbi_mail 
		where SecurityCategory = ? and isin = ? `
	lRows, lErr1 := lDb.Query(lCoreString, pBond.SecurityCategory, pBond.ISIN)
	if lErr1 != nil {
		log.Println("GCME01", lErr1)
		return lFlag, lErr1
	} else {
		for lRows.Next() {
			lErr2 := lRows.Scan(&Count)
			if lErr2 != nil {
				log.Println("GCME02", lErr2)
				return lFlag, lErr2
			} else {
				if Count == 0 {
					lFlag = "Y"

				} else {
					lFlag = "N"
				}

			}
		}

	}
	log.Println("checkMailExist(-)")

	return lFlag, nil

}
