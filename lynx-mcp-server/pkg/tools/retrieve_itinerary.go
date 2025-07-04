package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	TOOL_RETRIEVE_ITINERARY                         string = "retrieve_itinerary"
	TOOL_RETRIEVE_ITINERARY_DESCRIPTION             string = "Retrieve itinerary of a file"
	TOOL_RETRIEVE_ITINERARY_ARG_FILE_ID             string = "fileIdentifier"
	TOOL_RETRIEVE_ITINERARY_ARG_FILE_ID_DESCRIPTION string = "File identifier"

	TOOL_RETRIEVE_ITINERARY_SCHEMA = `{
		"type": "object",
		"description": "Retrieve file itinerary from file identifier",
		"properties": {
			"fileIdentifier": {
				"type": "string",
				"description":  "File identifier"
			}
		},
		"required": ["fileIdentifier"],
		"outputSchema": {
			"type": "object",
			"properties": {
				"count": {
					"type": "integer",
					"description": "Number of results found"
				},
				"results": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"supplier": {
								"type": "string",
								"description": "Supplier name"
							},
							"productName": {
								"type": "string",
								"description": "Product name"
							},
							"date": {
								"type": "string",
								"description": "Date"
							},
							"location": {
								"type": "string",
								"description": "Location"
							},
							"status": {
								"type": "string",
								"description": "Status"
							}
						},
						"required": ["supplier", "productName", "date", "location", "status"]
					}
				}
			},
			"required": ["count", "results"]
		}
	}`

	LYNX_RETRIEVE_ITINERARY_URL string = "/lynx/service/file.rpc"
)

type RetrieveItinerary struct {
	TotalBuyPrice  string             `json:"totalBuyPrice"`
	TotalSellPrice string             `json:"totalSellPrice"`
	Count          int                `json:"count"`
	Products       []ItineraryProduct `json:"products"`
}

type ItineraryProduct struct {
	Supplier    string `json:"supplier"`
	ProductName string `json:"productName"`
	Date        string `json:"date"`
	Location    string `json:"location"`
	Status      string `json:"status"`
}

func HandleRetrieveItinerary(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	session, _, err := utils.GetOrCreateSession(ctx, lynxConfig)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	arguments := request.GetArguments()
	fileIdentifier, ok := arguments["fileIdentifier"].(string)

	if !ok {
		return nil, fmt.Errorf("invalid number arguments")
	}

	body := utils.BuildGWTRetrieveItineraryBody(&utils.GWTRetrieveItineraryArgs{
		FileIdentifier: fileIdentifier,
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, LYNX_RETRIEVE_ITINERARY_URL), strings.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("failed to create retrieve itinerary request: %w", err)
	}

	req.Header.Set("Content-Type", utils.GWT_CONTENT_TYPE)
	req.AddCookie(utils.CreateAuthCookie(lynxConfig, session))

	// Use retry utility with exponential backoff
	resp, bodyStr, err := utils.RetryHTTPRequest(ctx, client, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute retrieve itinerary request after retries: %w", err)
	}
	defer resp.Body.Close()

	// Parse the GWT response body
	responseBody, err := utils.ParseGWTResponseBody(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GWT response: %w", err)
	}

	// Convert parsed data to structured format
	fileSearchResponse, err := parseRetrieveItineraryResponse(responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file search response: %w", err)
	}

	return utils.NewToolResultJSON(fileSearchResponse), nil
}

func parseRetrieveItineraryResponse(responseBody any) (*RetrieveItinerary, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// GetRetrieveItinerarySchema returns the complete JSON schema for the retrieve itinerary tool
func GetRetrieveItinerarySchema() json.RawMessage {
	return json.RawMessage(TOOL_RETRIEVE_ITINERARY_SCHEMA)
}
