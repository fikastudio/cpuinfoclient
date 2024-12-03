package cpuinfoclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"
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

type cacheItem struct {
	expiresAt time.Time
	item      *SearchResponse
}

type Client struct {
	locker *sync.RWMutex
	cache  map[string]cacheItem
}

func New() *Client {
	locker := &sync.RWMutex{}

	cl := &Client{
		cache:  map[string]cacheItem{},
		locker: locker,
	}

	go func() {
		for {
			locker.Lock()
			for k, v := range cl.cache {
				if v.expiresAt.After(time.Now()) {
					delete(cl.cache, k)
				}
			}
			locker.Unlock()

			time.Sleep(time.Second * 60)
		}
	}()

	return cl
}

func (cl *Client) Search(ctx context.Context, name string) (*SearchResponse, error) {
	cl.locker.RLock()
	cachedVal, ok := cl.cache[name]
	cl.locker.RUnlock()
	if ok {
		return cachedVal.item, nil
	}

	params := url.Values{}
	params.Set("search", name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://cpuinfo.fly.dev/v1/cpu?"+params.Encode(), nil)
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

	cl.locker.Lock()
	cl.cache[name] = cacheItem{expiresAt: time.Now().Add(time.Minute * 60), item: &out}
	cl.locker.Unlock()

	return &out, nil
}
