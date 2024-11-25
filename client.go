package cpuinfoclient

import (
	"context"
	"encoding/json"
	"net/http"
)

type SearchResponse struct {
	Schema           string `json:"$schema"`
	Cpus             []Cpus `json:"cpus"`
	ResultSetTDP     uint   `json:"resultSetTDP"`
	GlobalAverageTDP uint   `json:"globalAverageTDP"`
}
type Cpus struct {
	ProductName               string   `json:"productName"`
	LaunchDate                string   `json:"launchDate"`
	TotalCores                uint     `json:"totalCores"`
	MaxTurboFrequencyGHz      string   `json:"maxTurboFrequencyGHz"`
	ProcessorBaseFrequencyGHz string   `json:"processorBaseFrequencyGHz"`
	CacheMB                   uint     `json:"cacheMB"`
	TdpWatt                   uint     `json:"tdpWatt"`
	CPUInfoModelNames         []string `json:"cpuInfoModelNames"`
}

type Client struct{}

func New() *Client {
	return &Client{}
}

func (cl *Client) Search(ctx context.Context, name string) (*SearchResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://cpuinfo.fly.dev/v1/cpu?search="+name, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var out SearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
