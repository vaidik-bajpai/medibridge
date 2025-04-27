package handlers

func TestHandleRecordAllergy(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := dto.RegAllergyReq{
		AllergyName: "Peanut",
		Severity:    "High",
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.AllergyStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Patient UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Malformed JSON",
			urlID: validUUID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Validation fails on DTO",
			urlID: validUUID,
			body:  []byte(`{"allergyName":""}`),
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Record", mock.Anything, mock.MatchedBy(func(r *dto.RegAllergyReq) bool {
					return r.PatientID == validUUID && r.AllergyName == "Peanut"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Record", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := mocks.NewAllergyStorer(t)
			tt.mockSetup(as)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Allergy: as},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.urlID)

			req := httptest.NewRequest(http.MethodPost, "/patients/"+tt.urlID+"/allergies", bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleRecordAllergy(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleUpdateAllergy(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := dto.UpdateAllergyReq{
		AllergyName: "Peanut",
		Severity:    "Moderate",
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.AllergyStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Allergy UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Malformed JSON",
			urlID: validUUID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Validation fails on DTO",
			urlID: validUUID,
			body:  []byte(`{"allergyName":""}`),
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Update", mock.Anything, mock.MatchedBy(func(r *dto.UpdateAllergyReq) bool {
					return r.AllergyID == validUUID && r.AllergyName == "Peanut"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := mocks.NewAllergyStorer(t)
			tt.mockSetup(as)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Allergy: as},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("allergyID", tt.urlID)

			req := httptest.NewRequest(http.MethodPut, "/allergies/"+tt.urlID, bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleUpdateAllergy(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleDeleteAllergy(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name               string
		urlID              string
		mockSetup          func(*mocks.AllergyStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Allergy UUID",
			urlID: "invalid-uuid",
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Success",
			urlID: validUUID,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Delete", mock.Anything, validUUID).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Delete", mock.Anything, validUUID).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := mocks.NewAllergyStorer(t)
			tt.mockSetup(as)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Allergy: as},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("allergyID", tt.urlID)

			req := httptest.NewRequest(http.MethodDelete, "/allergies/"+tt.urlID, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleDeleteAllergy(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
