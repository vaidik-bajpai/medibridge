package handlers

func TestHandleAddCondition(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := dto.ConditionReq{
		ConditionName: "Asthma",
		Severity:      "Moderate",
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.ConditionStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Patient UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(cs *mocks.ConditionStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Malformed JSON",
			urlID: validUUID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(cs *mocks.ConditionStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Validation fails on DTO",
			urlID: validUUID,
			body:  []byte(`{"conditionName":""}`),
			mockSetup: func(cs *mocks.ConditionStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(cs *mocks.ConditionStorer) {
				cs.On("Add", mock.Anything, mock.MatchedBy(func(r *dto.ConditionReq) bool {
					return r.PatientID == validUUID && r.ConditionName == "Asthma"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(cs *mocks.ConditionStorer) {
				cs.On("Add", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := mocks.NewConditionStorer(t)
			tt.mockSetup(cs)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Condition: cs},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.urlID)

			req := httptest.NewRequest(http.MethodPost, "/patients/"+tt.urlID+"/conditions", bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleAddCondition(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleInactiveCondition(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name               string
		urlID              string
		mockSetup          func(*mocks.ConditionStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Condition UUID",
			urlID: "invalid-uuid",
			mockSetup: func(cs *mocks.ConditionStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Success",
			urlID: validUUID,
			mockSetup: func(cs *mocks.ConditionStorer) {
				cs.On("Deactivate", mock.Anything, validUUID).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			mockSetup: func(cs *mocks.ConditionStorer) {
				cs.On("Deactivate", mock.Anything, validUUID).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := mocks.NewConditionStorer(t)
			tt.mockSetup(cs)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Condition: cs},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("conditionID", tt.urlID)

			req := httptest.NewRequest(http.MethodPost, "/conditions/"+tt.urlID+"/inactive", nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleInactiveCondition(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
