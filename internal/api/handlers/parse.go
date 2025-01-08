package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type (
	ParamOptions struct {
		Numered  []string
		Stringed []string
	}

	ParamFuncOption func(c *ParamOptions)

	Params struct {
		ints map[string]int
		strs map[string]string
	}
)

func NumOpts(keys ...string) ParamFuncOption {
	return func(p *ParamOptions) {
		p.Numered = keys
	}
}

func StrOpts(keys ...string) ParamFuncOption {
	return func(p *ParamOptions) {
		p.Stringed = keys
	}
}

func ParseURLParameters(r *http.Request, opts ...ParamFuncOption) (ps *Params, err error) {

	ps = &Params{}

	parseFor := &ParamOptions{}
	for _, opt := range opts {
		opt(parseFor)
	}

	ps.strs, err = QueryURLParameters(r, parseFor.Stringed...)
	if err != nil {
		return nil, err
	}

	ps.ints, err = QueryURLNumeredParameters(r, parseFor.Numered...)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func QueryURLParameters(r *http.Request, params ...string) (map[string]string, error) {

	data := make(map[string]string, len(params))

	for _, param := range params {

		str := chi.URLParam(r, param)

		if str == "" || str == " " {
			return nil, NewErrorResponse(
				http.StatusNotAcceptable,
				ErrEmptyParameter(param),
			)
		}

		data[param] = str
	}

	return data, nil
}

func QueryURLNumeredParameters(r *http.Request, params ...string) (map[string]int, error) {

	stringPs, err := QueryURLParameters(r, params...)
	if err != nil {
		return nil, err
	}

	numPs := make(map[string]int, len(params))

	for _, param := range params {

		num, err := strconv.Atoi(stringPs[param])

		if err != nil {
			return nil, NewErrorResponse(
				http.StatusNotAcceptable,
				ErrInterpretationToNumber(param),
			)
		}

		numPs[param] = num
	}

	return numPs, nil
}

func (m *Params) GetInt(key string) int {
	v, ok := m.ints[key]
	if !ok {
		return 0
	}

	return v
}

func (m *Params) GetIntList(keys ...string) []int {
	lst := make([]int, len(keys))

	for i, key := range keys {
		lst[i] = m.GetInt(key)
	}

	return lst
}

func (m *Params) GetStr(key string) string {
	v, ok := m.strs[key]
	if !ok {
		return ""
	}

	return v
}

func (m *Params) GetStrList(keys ...string) []string {
	lst := make([]string, len(keys))

	for i, key := range keys {
		lst[i] = m.GetStr(key)
	}

	return lst
}
