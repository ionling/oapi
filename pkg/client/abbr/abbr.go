// https://www.abbreviations.com/abbr_api.php

package abbr

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

const abbrAPI = "https://www.stands4.com/services/v2/abbr.php"

type SearchType string

const (
	STExact   SearchType = "e" // Exact match
	STReverse SearchType = "r" // Reverse lookup
)

type commonTerm struct {
	ID           string
	Term         string
	Definition   string
	Category     string
	CategoryName string
	Score        string
}

type Term struct {
	commonTerm
	ParentCategory     string
	ParentCategoryName string
}

type GetAbbrsReq struct {
	UID        string
	TokenID    string
	Term       string
	SearchType SearchType
}

type GetAbbrsRes struct {
	Result []*Term `json:"result"`
}

func GetAbbrs(ctx context.Context, req *GetAbbrsReq) (res *GetAbbrsRes, err error) {
	hreq, err := http.NewRequestWithContext(ctx, "GET", abbrAPI, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	q := hreq.URL.Query()
	q.Add("uid", req.UID)
	q.Add("tokenid", req.TokenID)
	q.Add("searchtype", string(req.SearchType))
	q.Add("format", "json")
	q.Add("term", req.Term)
	hreq.URL.RawQuery = q.Encode()

	hres, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}

	bs, err := ioutil.ReadAll(hres.Body)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	if res, err = jsonToGetAbbrsRes(bs); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return
}

type rawTerm struct {
	commonTerm `mapstructure:",squash"`
	// Sometimes, the following two fields will be `{}`, not a string.
	// Then we will convert they to empty strings.
	ParentCategory     any
	ParentCategoryName any
}

func termFromRaw(in *rawTerm) (res *Term) {
	if in == nil {
		return nil
	}

	res = &Term{
		commonTerm: in.commonTerm,
	}
	if m, ok := in.ParentCategory.(map[string]any); !ok || len(m) > 0 {
		res.ParentCategory = fmt.Sprintf("%v", in.ParentCategory)
	}
	if m, ok := in.ParentCategoryName.(map[string]any); !ok || len(m) > 0 {
		res.ParentCategoryName = fmt.Sprintf("%v", in.ParentCategoryName)
	}
	return
}

func jsonToGetAbbrsRes(bs []byte) (res *GetAbbrsRes, err error) {
	m := map[string]any{}
	if err := json.Unmarshal(bs, &m); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	v, ok := m["result"]
	if !ok {
		return nil, fmt.Errorf("no result field")
	}

	var rawTerms []*rawTerm
	t := reflect.TypeOf(v)
	kind := t.Kind()
	switch kind {
	case reflect.Slice:
		err = mapstructure.Decode(v, &rawTerms)
	case reflect.Map:
		var t rawTerm
		err = mapstructure.Decode(v, &t)
		rawTerms = append(rawTerms, &t)
	default:
		err = fmt.Errorf("unsupported kind: %v", kind)
	}

	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	res = &GetAbbrsRes{}
	for _, t := range rawTerms {
		res.Result = append(res.Result, termFromRaw(t))
	}
	return
}
