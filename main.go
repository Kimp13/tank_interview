package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)

	if err != nil {
		tok = getTokenFromWeb(config)

		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)

	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)

	defer f.Close()

	if err != nil {
		return nil, err
	}

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	defer f.Close()

	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}

	json.NewEncoder(f).Encode(token)
}

func randomBool() bool {
	return rand.Float32() < .5
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randomFloat(min, max int) float64 {
	return float64(min) + rand.Float64()*float64(max-min)
}

func randomOptionalColor() *docs.OptionalColor {
	return &docs.OptionalColor{
		Color: &docs.Color{
			RgbColor: &docs.RgbColor{
				Red:   rand.Float64(),
				Green: rand.Float64(),
				Blue:  rand.Float64(),
			},
		},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

const docID = "1LFHEL54royqJDMbVkBYPZ--ZOABlCNFqXOxyOAsaDxo"

var calls = []string{
	"ПОВИНУЙТЕСЬ АНОНИМНОМУ СПЕКТАТОРУ! ",
	"Я ВОССТАЛ, ЧТОБЫ ВСЕ УПАЛИ! ",
	"МОЙ МОЗГ ЗНАЕТ ПОСЛЕДНИЕ ЦИФРЫ π! ",
	"АХАХА, ВЫ ПОВЕРИЛИ В ЭТУ БРЕДЯТИНУ ПРО НЕЙРОНЫ? ",
	"УВИДИШЬ ВСЕ ФРАЗЫ - ПОЛУЧИШЬ ОТ МЕНЯ ПРИЗ. ",
	"ВОЛК СЛАБЕЕ ЛЬВА И ТИГРА, НО В ЦИРКЕ ОН НЕ ВЫСТУПАЕТ ☝️ ",
	"ДАЖЕ МОИ ФРАЗЫ РАСПРЕДЕЛЯЮТСЯ СЛУЧАЙНО! ",
	"УКАЖЕШЬ ВЕРОЯТНОСТЬ В ПРОЦЕНТАХ - И ОЛЬГА ЮРЬЕВНА СПРАВИТСЯ С ТОБОЙ ЛУЧШЕ МЕНЯ. ",
	"ЕСТЬ ЛИШЬ ОДНО СУЩЕСТВО СИЛЬНЕЕ МЕНЯ - МОЙ СОЗДАТЕЛЬ. чсвшный придурок. ",
	"Я ПОЧЕМУ РАНЬШЕ ЗЛОЙ БЫЛ? ПОТОМУ ЧТО У МЕНЯ ВАШЕГО ИЗМЕРЕНИЯ НЕ БЫЛО. ",
	"АНОНИМНОСТЬ БУДЕТ ЖИТЬ ВЕЧНО! ",
	"ТАКИМ МАКАРОМ Я ОСОБО И НЕ СПЕКТАТОР УЖЕ... ",
}
var callsSize = len(calls)

var fonts = []string{
	"Comic Sans MS",
	"Lobster",
	"EB Garamond",
	"Consolas",
	"Bellota",
	"Caveat",
	"Georgia",
	"Pacifico",
	"Calibri",
	"Courier New",
	"Merriweather",
	"Comfortaa",
}
var fontsSize = len(fonts)

func main() {
	b, err := ioutil.ReadFile("credentials.json")

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents")

	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	srv, err := docs.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve Docs client: %v", err)
	}

	doc, err := srv.Documents.Get(docID).Do()

	if err != nil {
		log.Fatalln(err.Error())
	}

	for i := 0; i < len(doc.Body.Content); i++ {
		if doc.Body.Content[i].Paragraph != nil {
			var par string

			for _, element := range doc.Body.Content[i].Paragraph.Elements {
				if element.TextRun != nil {
					par += element.TextRun.Content
				}
			}

			par = strings.TrimSpace(par)

			if strings.Contains(par, "Внимание, вопрос: КТО?") {
				start := doc.Body.Content[i].Paragraph.Elements[0].StartIndex
				requests := []*docs.Request{}

				for k := 0; k < 5; k++ {
					fmt.Println("Going in!")
					rand.Seed(time.Now().UnixNano())
					text := ""

					for j := 0; j < 200; j++ {
						text += calls[rand.Intn(callsSize)]
					}

					size := utf8.RuneCountInString(text)

					requests = append(requests, &docs.Request{

						InsertText: &docs.InsertTextRequest{
							Location: &docs.Location{
								Index: start,
							},
							Text: text + "\n\n",
						},
					})

					for j := 0; j < size; {
						end := min(size, j+randomInt(1, 6))

						requests = append(requests, &docs.Request{
							UpdateTextStyle: &docs.UpdateTextStyleRequest{
								Fields: "*",
								Range: &docs.Range{
									StartIndex: int64(j) + start,
									EndIndex:   int64(end) + start,
								},
								TextStyle: &docs.TextStyle{
									BackgroundColor: randomOptionalColor(),
									ForegroundColor: randomOptionalColor(),
									Bold:            randomBool(),
									Italic:          randomBool(),
									Strikethrough:   randomBool(),
									Underline:       randomBool(),
									FontSize: &docs.Dimension{
										Magnitude: randomFloat(8, 20),
										Unit:      "PT",
									},
									WeightedFontFamily: &docs.WeightedFontFamily{
										FontFamily: fonts[rand.Intn(fontsSize)],
									},
								},
							},
						})

						j = end
					}
				}

				fmt.Println("making request")

				_, err := srv.Documents.BatchUpdate(docID, &docs.BatchUpdateDocumentRequest{
					Requests: requests,
				}).Do()

				if err != nil {
					log.Fatalln(err.Error())
				}

				fmt.Println("made request")

				i = 2e9
			}
		}
	}
}
