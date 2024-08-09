package integration

import (
	"avito_bootcamp/internal/entity"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetFlats(t *testing.T) {
	setup()
	defer teardown()

	house := entity.House{
		ID:        1,
		Address:   "123 Test St",
		Year:      2024,
		Developer: "Test Developer",
	}

	houseBody, err := json.Marshal(house)
	if err != nil {
		t.Fatalf("failed to marshal house request body: %v", err)
	}

	moderatorToken, err := getJWTToken("moderator")
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	clientToken, err := getJWTToken("client")
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// создаем дом
	req, err := authRequest("POST", fmt.Sprintf("%s/house/create", baseURL), houseBody, moderatorToken)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}

	// добавляем квартиры к дому, но они добавляются со статусом created и автоинкрементным id
	flats := []entity.Flat{
		{ID: 1, HouseID: 1, Number: 101, Status: "approved"},
		{ID: 2, HouseID: 1, Number: 102, Status: "pending"},
		{ID: 3, HouseID: 1, Number: 103, Status: "approved"},
	}

	for _, flat := range flats {
		flatBody, err := json.Marshal(flat)
		if err != nil {
			t.Fatalf("failed to marshal flat request body: %v", err)
		}

		req, err := authRequest("POST", fmt.Sprintf("%s/flat/create", baseURL), flatBody, moderatorToken)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		_, err = client.Do(req)
		if err != nil {
			t.Fatalf("failed to send POST request: %v", err)
		}
	}

	// проверяем получение квартир как модератор (должны увидеть все квартиры)
	req, err = authRequest("GET", fmt.Sprintf("%s/house/%d", baseURL, house.ID), nil, moderatorToken)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseMod map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMod); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	flatsMod, ok := responseMod["flats"].([]interface{})
	assert.True(t, ok, "Expected response to contain 'flats' field for moderator")
	assert.Equal(t, len(flats), len(flatsMod), "Moderator should see all flats")

	// проверяем получение квартир как клиент (должен увидеть только квартиры со статусом approved)

	// у одной квартиры меняем статус на approved
	flatApproved := entity.Flat{ID: 1, HouseID: 1, Status: "approved"}
	flatApprovedBody, err := json.Marshal(flatApproved)
	if err != nil {
		t.Fatalf("failed to marshal flat request body: %v", err)
	}
	req, err = authRequest("POST", fmt.Sprintf("%s/flat/update", baseURL), flatApprovedBody, moderatorToken)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}

	req, err = authRequest("GET", fmt.Sprintf("%s/house/%d", baseURL, house.ID), nil, clientToken)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseClient map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseClient); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	flatsClient, ok := responseClient["flats"].([]interface{})
	assert.True(t, ok, "Expected response to contain 'flats' field for client")

	// клиент должен видеть только approved квартиры
	expectedClientFlats := []entity.Flat{
		{ID: 1, HouseID: 1, Number: 101, Status: "approved"},
	}
	assert.Equal(t, len(expectedClientFlats), len(flatsClient), "Client should see only approved flats")
}
