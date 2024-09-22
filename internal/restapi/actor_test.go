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

func TestHandler_ActorCreate(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockActorService, a m.CreateActor)

	testActor := models.Actor{
		Id:        1,
		Name:      "Name 1",
		Gender:    "M",
		BirthDate: "1995-01-12",
	}

	tests := []struct {
		name                 string
		inputBody            string
		input                m.CreateActor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name":"Name 1","gender":"M","birth_date":"1995-01-12"}`,
			input: m.CreateActor{
				Name:      "Name 1",
				Gender:    "M",
				BirthDate: "1995-01-12",
			},
			mockBehavior: func(r *mock_restapi.MockActorService, a m.CreateActor) {
				r.EXPECT().Create(a).Return(testActor, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1,"name":"Name 1","gender":"M","birth_date":"1995-01-12"}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            "",
			input:                m.CreateActor{},
			mockBehavior:         func(r *mock_restapi.MockActorService, a m.CreateActor) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid request json decoder: EOF"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name":"Name 1","gender":"M","birth_date":"1995-01-12"}`,
			input: m.CreateActor{
				Name:      "Name 1",
				Gender:    "M",
				BirthDate: "1995-01-12",
			},
			mockBehavior: func(r *mock_restapi.MockActorService, a m.CreateActor) {
				r.EXPECT().Create(a).Return(models.Actor{}, errors.New(`internal error`))
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
			svc := mock_restapi.NewMockActorService(c)
			tt.mockBehavior(svc, tt.input)
			NewActorHandler(svc).Register(r)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/actors",
				bytes.NewBufferString(tt.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_ActorFind(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockActorService, id string)

	testActor := models.Actor{
		Id:        1,
		Name:      "Name 1",
		Gender:    "M",
		BirthDate: "1995-01-12",
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
			mockBehavior: func(r *mock_restapi.MockActorService, id string) {
				r.EXPECT().Find(id).Return(testActor, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"name":"Name 1","gender":"M","birth_date":"1995-01-12"}`,
		},
		{
			name:      "Service Error",
			inputBody: ``,
			input:     "1",
			mockBehavior: func(r *mock_restapi.MockActorService, id string) {
				r.EXPECT().Find(id).Return(models.Actor{}, errors.New(`internal error`))
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
			svc := mock_restapi.NewMockActorService(c)
			tt.mockBehavior(svc, tt.input)
			NewActorHandler(svc).Register(r)

			// Create Request
			w := httptest.NewRecorder()
			reqUrl := "/actors/" + tt.input
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

func TestHandler_ActorSearch(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockActorService)

	actors := []models.Actor{

		{
			Id:        1,
			Name:      "Name 1",
			Gender:    "M",
			BirthDate: "1995-01-12",
		},
		{
			Id:        2,
			Name:      "Name 2",
			Gender:    "F",
			BirthDate: "1995-02-12",
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
			mockBehavior: func(r *mock_restapi.MockActorService) {
				r.EXPECT().Search().Return(actors, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"name":"Name 1","gender":"M","birth_date":"1995-01-12"},{"id":2,"name":"Name 2","gender":"F","birth_date":"1995-02-12"}]`,
		},
		{
			name:      "Service Error",
			inputBody: ``,
			mockBehavior: func(r *mock_restapi.MockActorService) {
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
			req := httptest.NewRequest("GET", "/actors", bytes.NewBufferString(""))

			r := mux.NewRouter()
			svc := mock_restapi.NewMockActorService(c)
			tt.mockBehavior(svc)
			NewActorHandler(svc).Register(r)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_ActorDelete(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockActorService, id string)

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
			mockBehavior: func(r *mock_restapi.MockActorService, id string) {
				r.EXPECT().Delete(id).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{}`,
		},
		{
			name:      "Service Error",
			input:     "1",
			inputBody: ``,
			mockBehavior: func(r *mock_restapi.MockActorService, id string) {
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
			svc := mock_restapi.NewMockActorService(c)
			tt.mockBehavior(svc, tt.input)
			NewActorHandler(svc).Register(r)

			// Create Request
			w := httptest.NewRecorder()
			reqUrl := "/actors/" + tt.input
			req := httptest.NewRequest("DELETE", reqUrl, bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_ActorUpdate(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_restapi.MockActorService, a m.UpdateActor, id string)

	testActor := m.UpdateActor{
		Name:      "Name 1",
		Gender:    "M",
		BirthDate: "1995-01-12",
	}

	tests := []struct {
		name                 string
		input                m.UpdateActor
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			input:     testActor,
			inputBody: `{"name":"Name 1","gender":"M","birth_date":"1995-01-12"}`,
			mockBehavior: func(r *mock_restapi.MockActorService, a m.UpdateActor, id string) {
				r.EXPECT().Update(id, a).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{}`,
		},
		{
			name:      "Service Error",
			input:     testActor,
			inputBody: `{"name":"Name 1","gender":"M","birth_date":"1995-01-12"}`,
			mockBehavior: func(r *mock_restapi.MockActorService, a m.UpdateActor, id string) {
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
			req := httptest.NewRequest("PUT", "/actors/1",
				bytes.NewBufferString(tt.inputBody))

			r := mux.NewRouter()
			svc := mock_restapi.NewMockActorService(c)
			// id := mux.Vars(req)["id"]
			tt.mockBehavior(svc, tt.input, "1")
			NewActorHandler(svc).Register(r)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
