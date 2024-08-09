package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const baseURL = "http://localhost:8080"

var jwtToken string

func setup() {
	cmd := exec.Command("docker-compose", "-f", "test/integration/docker-compose.yml", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "./../.."
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to start docker-compose: %v", err)
	}

	time.Sleep(10 * time.Second)
	log.Println("Starting...")

	token, err := getJWTToken("moderator")
	if err != nil {
		log.Fatalf("failed to get JWT token: %v", err)
	}
	jwtToken = token
}

func teardown() {
	cmd := exec.Command("docker-compose", "down")
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to stop docker-compose: %v", err)
	}
}

func getJWTToken(userType string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/dummyLogin?user_type=%s", baseURL, userType))
	if err != nil {
		return "", fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	token, ok := response["token"]
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}

func authRequest(method, url string, body []byte, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
