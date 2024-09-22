package restapi

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"filmoteka/internal/app/models"
	"filmoteka/internal/restapi/mock_restapi"
	m "filmoteka/internal/restapi/models"
)

func TestHandler_FilmCreate(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockFilmService, a m.CreateFilm)

	testFilm := models.Film{
		Id:          1,
		Name:        "Test Name",
		Description: "Desc1",
		ReleaseYear: 2002,
		Rating:      7.5,
	}

	tests := []struct {
		name                 string
		inputBody            string
		input                m.CreateFilm
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name":"Test Name","description":"Desc1","release_year":2002,"rating":7.5}`,
			input: m.CreateFilm{
				Name:        "Test Name",
				Description: "Desc1",
				ReleaseYear: 2002,
				Rating:      7.5,
			},
			mockBehavior: func(r *mock_restapi.MockFilmService, a m.CreateFilm) {
				r.EXPECT().Create(a).Return(testFilm, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1,"name":"Test Name","description":"Desc1","release_year":2002,"rating":7.5}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            "",
			input:                m.CreateFilm{},
			mockBehavior:         func(r *mock_restapi.MockFilmService, a m.CreateFilm) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid request json decoder: EOF"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name":"Test Name","description":"Desc1","release_year":2002,"rating":7.5}`,
			input: m.CreateFilm{
				Name:        "Test Name",
				Description: "Desc1",
				ReleaseYear: 2002,
				Rating:      7.5,
			},
			mockBehavior: func(r *mock_restapi.MockFilmService, a m.CreateFilm) {
				r.EXPECT().Create(a).Return(models.Film{}, errors.New(`internal error`))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			r := mux.NewRouter()
			svc := mock_restapi.NewMockFilmService(c)
			tt.mockBehavior(svc, tt.input)
			NewFilmHandler(svc).Register(r)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/films",
				bytes.NewBufferString(tt.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_FilmFind(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockFilmService, id string)

	testFilm := models.Film{
		Id:          1,
		Name:        "Test Name",
		Description: "Desc1",
		ReleaseYear: 2002,
		Rating:      7.5,
	}

	tests := []struct {
		name                 string
		inputBody            string
		input                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: ``,
			input:     "1",
			mockBehavior: func(r *mock_restapi.MockFilmService, id string) {
				r.EXPECT().Find(id).Return(testFilm, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"name":"Test Name","description":"Desc1","release_year":2002,"rating":7.5}`,
		},
		{
			name:      "Service Error",
			inputBody: ``,
			input:     "1",
			mockBehavior: func(r *mock_restapi.MockFilmService, id string) {
				r.EXPECT().Find(id).Return(models.Film{}, errors.New(`internal error`))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			r := mux.NewRouter()
			svc := mock_restapi.NewMockFilmService(c)
			tt.mockBehavior(svc, tt.input)
			NewFilmHandler(svc).Register(r)

			// Create Request
			w := httptest.NewRecorder()
			reqUrl := "/films/" + tt.input
			req := httptest.NewRequest("GET", reqUrl, bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			// fmt.Println("Body :", w.Body.String())
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_FilmSearch(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockFilmService)

	actors := []models.Film{

		{
			Id:          1,
			Name:        "Test Name",
			Description: "Desc1",
			ReleaseYear: 2002,
			Rating:      7.5,
		},
		{
			Id:          2,
			Name:        "Test Name 2",
			Description: "Desc2",
			ReleaseYear: 2002,
			Rating:      7.5,
		},
	}

	tests := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: ``,
			mockBehavior: func(r *mock_restapi.MockFilmService) {
				r.EXPECT().Search().Return(actors, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"name":"Test Name","description":"Desc1","release_year":2002,"rating":7.5},{"id":2,"name":"Test Name 2","description":"Desc2","release_year":2002,"rating":7.5}]`,
		},
		{
			name:      "Service Error",
			inputBody: ``,
			mockBehavior: func(r *mock_restapi.MockFilmService) {
				r.EXPECT().Search().Return(actors, errors.New(`internal error`))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/films", bytes.NewBufferString(""))

			r := mux.NewRouter()
			svc := mock_restapi.NewMockFilmService(c)
			tt.mockBehavior(svc)
			NewFilmHandler(svc).Register(r)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_FilmDelete(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockFilmService, id string)

	tests := []struct {
		name                 string
		input                string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			input:     "1",
			inputBody: ``,
			mockBehavior: func(r *mock_restapi.MockFilmService, id string) {
				r.EXPECT().Delete(id).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{}`,
		},
		{
			name:      "Service Error",
			input:     "1",
			inputBody: ``,
			mockBehavior: func(r *mock_restapi.MockFilmService, id string) {
				r.EXPECT().Delete(id).Return(errors.New(`internal error`))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			r := mux.NewRouter()
			svc := mock_restapi.NewMockFilmService(c)
			tt.mockBehavior(svc, tt.input)
			NewFilmHandler(svc).Register(r)

			// Create Request
			w := httptest.NewRecorder()
			reqUrl := "/films/" + tt.input
			req := httptest.NewRequest("DELETE", reqUrl, bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_FilmUpdate(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockFilmService, a m.UpdateFilm, id string)

	testFilm := m.UpdateFilm{
		Name:        "Test Name",
		Description: "Desc1",
		ReleaseYear: 2002,
		Rating:      7.5,
	}

	tests := []struct {
		name                 string
		input                m.UpdateFilm
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			input:     testFilm,
			inputBody: `{"name":"Test Name","description":"Desc1","release_year":2002,"rating":7.5}`,
			mockBehavior: func(r *mock_restapi.MockFilmService, a m.UpdateFilm, id string) {
				r.EXPECT().Update(id, a).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{}`,
		},
		{
			name:      "Service Error",
			input:     testFilm,
			inputBody: `{"name":"Test Name","description":"Desc1","release_year":2002,"rating":7.5}`,
			mockBehavior: func(r *mock_restapi.MockFilmService, a m.UpdateFilm, id string) {
				r.EXPECT().Update(id, a).Return(errors.New(`internal error`))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/films/1",
				bytes.NewBufferString(tt.inputBody))

			r := mux.NewRouter()
			svc := mock_restapi.NewMockFilmService(c)
			// id := mux.Vars(req)["id"]
			tt.mockBehavior(svc, tt.input, "1")
			NewFilmHandler(svc).Register(r)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
