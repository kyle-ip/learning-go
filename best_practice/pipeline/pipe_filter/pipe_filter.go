// Package pipe_filter is to define the interfaces and the structures for pipe-filter style implementation
package pipe_filter

import (
	"errors"
	"strconv"
	"strings"
)

type Request interface{}

type Response interface{}

type Filter interface {
	Process(data Request) (Response, error)
}

var SplitFilterWrongFormatError = errors.New("input data should be string")

// ========== SplitFilter ==========

type SplitFilter struct {
	delimiter string
}

func NewSplitFilter(delimiter string) *SplitFilter {
	return &SplitFilter{delimiter}
}

func (sf *SplitFilter) Process(data Request) (Response, error) {
	str, ok := data.(string) //检查数据格式/类型，是否可以处理
	if !ok {
		return nil, SplitFilterWrongFormatError
	}
	parts := strings.Split(str, sf.delimiter)
	return parts, nil
}

// ========== SumFilter ==========

var SumFilterWrongFormatError = errors.New("input data should be []int")

type SumFilter struct {
}

func NewSumFilter() *SumFilter {
	return &SumFilter{}
}

func (sf *SumFilter) Process(data Request) (Response, error) {
	elems, ok := data.([]int)
	if !ok {
		return nil, SumFilterWrongFormatError
	}
	ret := 0
	for _, elem := range elems {
		ret += elem
	}
	return ret, nil
}

// ========== ToIntFilter ==========

var ToIntFilterWrongFormatError = errors.New("input data should be []string")

type ToIntFilter struct {
}

func NewToIntFilter() *ToIntFilter {
	return &ToIntFilter{}
}

func (tif *ToIntFilter) Process(data Request) (Response, error) {
	parts, ok := data.([]string)
	if !ok {
		return nil, ToIntFilterWrongFormatError
	}
	var ret []int
	for _, part := range parts {
		s, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

// ========== StraightPipeline ==========

func NewStraightPipeline(name string, filters ...Filter) *StraightPipeline {
	return &StraightPipeline{
		Name:    name,
		Filters: &filters,
	}
}

type StraightPipeline struct {
	Name    string
	Filters *[]Filter
}

func (f *StraightPipeline) Process(data Request) (Response, error) {
	var ret interface{}
	var err error
	for _, filter := range *f.Filters {
		ret, err = filter.Process(data)
		if err != nil {
			return ret, err
		}
		data = ret
	}
	return ret, err
}
