package googlesheets

import (
	"context"
	"encoding/base64"
	"fmt"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const scopeSpreadsheets = "https://www.googleapis.com/auth/spreadsheets"

type Client struct {
	SheetsService *sheets.Service
}

func New(credsBase64 string) (*Client, error) {
	ctx := context.Background()

	credBytes, err := base64.StdEncoding.DecodeString(credsBase64)
	if err != nil {
		return nil, fmt.Errorf("decode base64 string: %w", err)
	}

	config, err := google.JWTConfigFromJSON(credBytes, scopeSpreadsheets)
	if err != nil {
		return nil, fmt.Errorf("jwt config from json: %w", err)
	}

	client := config.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("new service: %w", err)
	}

	return &Client{
		SheetsService: srv,
	}, nil
}

func (c *Client) GetValuesByRange(ctx context.Context, sheetID, range_ string) ([][]any, error) {
	resp, err := c.SheetsService.Spreadsheets.Values.Get(sheetID, range_).Do()
	if err != nil || resp.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("get values: %w", err)
	}

	return resp.Values, nil
}

func (c *Client) InsertValues(ctx context.Context, sheetID, range_ string, values [][]any) error {
	resp, err := c.SheetsService.Spreadsheets.Values.Append(
		sheetID, range_, &sheets.ValueRange{Values: values},
	).ValueInputOption("USER_ENTERED").Do()
	if err != nil || resp.HTTPStatusCode != 200 {
		return fmt.Errorf("insert values: %w", err)
	}

	return nil
}
