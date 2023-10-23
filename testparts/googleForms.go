package testparts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/abaskin/Test-Maker/v2/testparts/resources"

	"github.com/phayes/freeport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

type Oauth2ResponseSt struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

type GoogleFormStatus int32

const (
	Connecting GoogleFormStatus = iota
	Available
	Unavailable
	Unknown
)

type GoogleFormSt struct {
	form            *forms.Form
	Status          GoogleFormStatus
	authWait        *sync.WaitGroup
	authCode        string
	port            int
	service         *forms.FormsService
	nextItem        int64
	formCredentials string
}

var prefStore *Preference

func NewGoogleForm(dsn, formCredentials string) *GoogleFormSt {
	port, err := freeport.GetFreePort()
	if err != nil {
		port = 8321
	}

	prefStore = newPreference(dsn)

	return &GoogleFormSt{
		authWait:        &sync.WaitGroup{},
		port:            port,
		nextItem:        0,
		formCredentials: formCredentials,
	}
}

func (gf *GoogleFormSt) Create(title, documentTitle, desc string) (*GoogleFormSt, error) {
	ctx := context.Background()

	formCredentials := fmt.Sprintf(gf.formCredentials, gf.port)

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(
		[]byte(formCredentials),
		forms.FormsBodyScope,
	)
	if err != nil {
		gf.Status = Unavailable
		return gf, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	gf.Status = Connecting
	gf.startHTTPServer()

	for {
		if err := gf.formService(config, ctx); err != nil {
			return gf, err
		}

		err := gf.formCreate(title, documentTitle, config, ctx)

		if err != nil {
			if err := gf.authError(err); err == nil {
				continue
			}
			return gf, err
		}

		break
	}

	return gf, gf.formInit(desc)
}

func (gf *GoogleFormSt) authError(err error) error {
	response := &Oauth2ResponseSt{}
	_, errJson, haveJson := strings.Cut(err.Error(), "{")
	if haveJson {
		if err := json.Unmarshal([]byte("{"+errJson), response); err != nil {
			gf.Status = Unavailable
			return fmt.Errorf("unable to parse error response %v", err)
		}
		if response.Error == "unauthorized_client" {
			prefStore.Del("Forms.Token")
			return nil
		}
	}
	for {
		nested := errors.Unwrap(err)
		if nested != nil {
			err = nested
			continue
		}
		break
	}
	gf.Status = Unavailable
	return fmt.Errorf(" Google Forms Unavailable, %v", err)
}

func (gf *GoogleFormSt) formService(config *oauth2.Config, ctx context.Context) error {
	client := getClient("Forms.Token", config, gf.authWait, &gf.authCode)

	baseService, err := forms.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		gf.Status = Unavailable
		return fmt.Errorf("unable to create forms Client %v", err)
	}
	gf.service = forms.NewFormsService(baseService)

	return err
}

func (gf *GoogleFormSt) formCreate(title, documentTitle string,
	config *oauth2.Config, ctx context.Context) error {
	var err error
	gf.form, err = gf.service.Create(&forms.Form{
		Info: &forms.Info{
			Title:         title,
			DocumentTitle: documentTitle,
		},
	}).Do()

	return err
}

func (gf *GoogleFormSt) formInit(desc string) error {
	response, err := gf.service.BatchUpdate(gf.form.FormId,
		&forms.BatchUpdateFormRequest{
			IncludeFormInResponse: true,
			Requests: []*forms.Request{
				{
					UpdateFormInfo: &forms.UpdateFormInfoRequest{
						Info: &forms.Info{
							Description: desc,
						},
						UpdateMask: "description",
					},
				},
				{
					UpdateSettings: &forms.UpdateSettingsRequest{
						Settings: &forms.FormSettings{
							QuizSettings: &forms.QuizSettings{
								IsQuiz: true,
							},
						},
						UpdateMask: "quizSettings.isQuiz",
					},
				},
			},
		},
	).Do()
	if err == nil {
		gf.form = response.Form
	}

	return err
}

func (gf *GoogleFormSt) AddSection(title, desc string) error {
	response, err := gf.service.BatchUpdate(gf.form.FormId,
		&forms.BatchUpdateFormRequest{
			IncludeFormInResponse: true,
			Requests: []*forms.Request{
				{
					CreateItem: &forms.CreateItemRequest{
						Item: &forms.Item{
							Title:         title,
							Description:   desc,
							PageBreakItem: &forms.PageBreakItem{},
						},
						Location: &forms.Location{
							Index:           gf.nextItem,
							ForceSendFields: []string{"Index"},
						},
					},
				},
			},
		},
	).Do()

	if err == nil {
		gf.form = response.Form
		gf.nextItem = int64(len(gf.form.Items))
	}
	return err
}

func (gf *GoogleFormSt) AddQuestions(questions QuestionSetSt,
	qPoints int64) error {
	q, _ := questions.Get(0)
	if q.Choices.Size() != 0 {
		return gf.addChoiceQuestions(questions, qPoints)
	}
	return gf.addTextQuestions(questions, qPoints)
}

func (gf *GoogleFormSt) addTextQuestions(questions QuestionSetSt,
	qPoints int64) error {
	request := make([]*forms.Request, questions.Size())
	questions.Each(
		func(i int, q *QuestionsSt) {
			request[i] =
				&forms.Request{
					CreateItem: &forms.CreateItemRequest{
						Location: &forms.Location{
							Index:           gf.nextItem,
							ForceSendFields: []string{"Index"},
						},
						Item: &forms.Item{
							Title: q.Question.CleanString(),
							QuestionItem: &forms.QuestionItem{
								Question: &forms.Question{
									Required: true,
									TextQuestion: &forms.TextQuestion{
										Paragraph: true,
									},
									Grading: &forms.Grading{
										PointValue: qPoints,
									},
								},
							},
						},
					},
				}
			gf.nextItem++
		},
	)

	response, err := gf.service.BatchUpdate(gf.form.FormId,
		&forms.BatchUpdateFormRequest{
			IncludeFormInResponse: true,
			Requests:              request,
		},
	).Do()

	if err == nil {
		gf.form = response.Form
	}
	return err
}

func (gf *GoogleFormSt) addChoiceQuestions(questions QuestionSetSt,
	qPoints int64) error {
	request := make([]*forms.Request, questions.Size())
	questions.Each(
		func(i int, q *QuestionsSt) {
			options := make([]*forms.Option, q.Choices.Size())
			q.Choices.Each(
				func(i int, opt string) {
					options[i] =
						&forms.Option{
							Value: opt,
						}
				},
			)
			answerStr, _ := q.Answers.Get(1)
			request[i] = &forms.Request{
				CreateItem: &forms.CreateItemRequest{
					Location: &forms.Location{
						Index:           gf.nextItem,
						ForceSendFields: []string{"Index"},
					},
					Item: &forms.Item{
						Title: q.Question.CleanString(),
						QuestionItem: &forms.QuestionItem{
							Question: &forms.Question{
								Required: true,
								ChoiceQuestion: &forms.ChoiceQuestion{
									Shuffle: true,
									Type:    "RADIO",
									Options: options,
								},
								Grading: &forms.Grading{
									PointValue: qPoints,
									CorrectAnswers: &forms.CorrectAnswers{
										Answers: []*forms.CorrectAnswer{
											{
												Value: answerStr,
											},
										},
									},
								},
							},
						},
					},
				},
			}
			gf.nextItem++
		},
	)

	response, err := gf.service.BatchUpdate(gf.form.FormId,
		&forms.BatchUpdateFormRequest{
			IncludeFormInResponse: true,
			Requests:              request,
		},
	).Do()

	if err == nil {
		gf.form = response.Form
	}
	return err
}

func (gf *GoogleFormSt) startHTTPServer() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", gf.port),
		Handler: newHandler(gf.authWait, &gf.authCode),
	}

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server closed\n")
		} else if err != nil {
			log.Printf("error listening for http server: %s\n", err)
		}
	}()
}

func newHandler(authWait *sync.WaitGroup, authCode *string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/google/callback",
		func(writer http.ResponseWriter, request *http.Request) {
			u, err := url.Parse(request.RequestURI)
			if err != nil {
				log.Println("HTTP Server ", err)
			}
			*authCode = u.Query()["code"][0]
			authWait.Done()

			writer.Write(resources.GoogleFormAuthSuccess)
		},
	)

	return mux
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(key string, config *oauth2.Config, tokenWait *sync.WaitGroup,
	tokenCode *string) *http.Client {
	tok, err := tokenFromPref(key)
	if err != nil {
		tok = getTokenFromWeb(config, tokenWait, tokenCode)
		saveTokenPref(key, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config, authWait *sync.WaitGroup,
	authCode *string) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	authWait.Add(1)
	url, _ := url.Parse(authURL)
	log.Println("Please open the URL and follow the instructions,", url)
	authWait.Wait()

	tok, err := config.Exchange(context.TODO(), *authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return tok
}

// Retrieves a token from preferences.
func tokenFromPref(key string) (*oauth2.Token, error) {
	tok := &oauth2.Token{}
	data, err := prefStore.Get(key)
	if err != nil {
		return tok, err
	}
	err = json.NewDecoder(strings.NewReader(data)).Decode(tok)
	return tok, err
}

// Saves a token to preferences.
func saveTokenPref(key string, token *oauth2.Token) error {
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(token)
	return prefStore.Set(key, buf.String())
}
