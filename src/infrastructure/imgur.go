package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ImgurResponse struct {
	Data struct {
		Link string `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

func ImgurUploadImage(imagePath, title, description string) (string, error) {
	clientID := os.Getenv("IMGUR_CLIENT_ID")

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}(file)

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copying file: %v", err)
	}

	err = writer.WriteField("type", "file")
	if err != nil {
		return "", err
	}
	err = writer.WriteField("title", title)
	if err != nil {
		return "", err
	}
	err = writer.WriteField("description", description)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", &requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Client-ID "+clientID)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("error closing body: %v", err.Error())
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("imgur response error: %s", string(body))
	}

	var imgurResp ImgurResponse
	err = json.Unmarshal(body, &imgurResp)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON response: %v", err)
	}

	if !imgurResp.Success {
		return "", fmt.Errorf("imgur API error: %d", imgurResp.Status)
	}

	return imgurResp.Data.Link, nil
}

func ImgurGetImage(imageID, savePath string) error {
	url := fmt.Sprintf("https://i.imgur.com/%s.jpg", imageID)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("error closing body: %v", err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching image: status %d", resp.StatusCode)
	}

	outFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}(outFile)

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving image: %v", err)
	}

	return nil
}
