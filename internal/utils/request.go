package utils

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Link  string `json:"link"`
	Depth int    `json:"id,omitempty"`
}

func (r *Request) UnmarshalJSON(d []byte) error {
	tmp := struct {
		Link  interface{} `json:"link"`
		Depth interface{} `json:"depth"`
	}{}

	if err := json.Unmarshal(d, &tmp); err != nil {
		return err
	}

	switch v := tmp.Depth.(type) {
	case float64:
		r.Depth = int(v)
	default:
		return fmt.Errorf("invalid value for depth: %v", v)
	}

	r.Link = tmp.Link.(string)
	return nil
}
