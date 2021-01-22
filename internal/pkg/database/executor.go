package database

import (
	"errors"
	"reflect"
	"time"
)

var TimeoutError = errors.New("statement timed out")

func PrepareAsync(timeout time.Duration, query string, values ...interface{}) error {
	ch := make(chan error, 1)

	go func() {
		ch <- PrepareStatement(query, values...)
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(timeout):
		return TimeoutError
	}
}

func NamedPrepareAsync(timeout time.Duration, query string, values interface{}) error {
	ch := make(chan error, 1)

	go func() {
		ch <- NamedPrepareStatement(query, values)
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(timeout):
		return TimeoutError
	}
}

func QueryAsync(timeout time.Duration, structType reflect.Type, query string, values ...interface{}) (interface{}, error) {
	ch := make(chan *QueryResult, 1)

	go func() {
		ch <- Query(structType, query, values...)
	}()

	select {
	case result := <-ch:
		return result.Results, result.Error
	case <-time.After(timeout):
		return nil, TimeoutError
	}
}

func NamedQueryAsync(timeout time.Duration, object interface{}, query string, values interface{}) error {
	ch := make(chan error, 1)

	go func() {
		ch <- NamedQuery(object, query, values)
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(timeout):
		return TimeoutError
	}
}
