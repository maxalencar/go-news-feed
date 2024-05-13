package news

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"go-news-feed/pkg/model"
)

type TestSuite struct {
	suite.Suite
	router      *http.ServeMux
	serviceMock *MockService
	article     model.Article
}

func (suite *TestSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	suite.serviceMock = NewMockService(ctrl)
	suite.router = newEndpoint(suite.serviceMock).init()
}

func (suite *TestSuite) SetupTest() {
	suite.article = model.Article{
		ID:           "test id",
		Title:        "test title",
		Descriptiopn: "test description",
		Link:         "test link",
	}
}

func (suite *TestSuite) TestFind() {
	testCases := []struct {
		name         string
		given        any
		mockCalls    func()
		expectedCode int
		expected     any
	}{
		{
			name:  "FindInternalServerError",
			given: model.FindRequest{},
			mockCalls: func() {
				suite.serviceMock.EXPECT().Find(gomock.Any(), gomock.Any()).Return(model.FindResponse{}, errors.New("internal server error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:  "FindSuccess",
			given: model.FindRequest{},
			mockCalls: func() {
				suite.serviceMock.EXPECT().Find(gomock.Any(), gomock.Any()).Return(model.FindResponse{
					Articles: []model.Article{suite.article},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expected:     []model.Article{suite.article},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockCalls()

			b, err := json.Marshal(tc.given)
			suite.NoError(err)

			reqBody := bytes.NewReader(b)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/find", reqBody)

			suite.router.ServeHTTP(w, r)

			suite.Equal(tc.expectedCode, w.Code)

			if tc.expected != nil {
				res := w.Result()
				defer res.Body.Close()

				var response model.FindResponse

				decoder := json.NewDecoder(res.Body)
				err = decoder.Decode(&response)
				suite.NoError(err)

				suite.NotEmpty(response)
				suite.Equal(tc.given, response.Criteria)
				suite.Equal(tc.expected, response.Articles)
			}
		})
	}
}

func (suite *TestSuite) TestLoad() {
	testCases := []struct {
		name         string
		given        string
		mockCalls    func()
		expectedCode int
		expected     any
	}{
		{
			name: "LoadInternalServerError",
			mockCalls: func() {
				suite.serviceMock.EXPECT().Load(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal server error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "LoadSuccess",
			mockCalls: func() {
				suite.serviceMock.EXPECT().Load(gomock.Any(), gomock.Any()).Return([]model.Article{suite.article}, nil)
			},
			expectedCode: http.StatusOK,
			expected:     []model.Article{suite.article},
		},
		{
			name:  "LoadSuccessWithURL",
			given: "test",
			mockCalls: func() {
				suite.serviceMock.EXPECT().Load(gomock.Any(), "test").Return([]model.Article{suite.article}, nil)
			},
			expectedCode: http.StatusOK,
			expected:     []model.Article{suite.article},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockCalls()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?feedUrl=%s", "/load", tc.given), nil)

			suite.router.ServeHTTP(w, r)

			suite.Equal(tc.expectedCode, w.Code)

			if tc.expected != nil {
				res := w.Result()
				defer res.Body.Close()

				var articles []model.Article

				decoder := json.NewDecoder(res.Body)
				err := decoder.Decode(&articles)
				suite.NoError(err)

				suite.NotEmpty(articles)
				suite.Equal(tc.expected, articles)
			}
		})
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
