// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"sync"

	"github.com/ONSdigital/dis-search-upstream-stub/api"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
)

// Ensure, that DataStorerMock does implement api.DataStorer.
// If this is not the case, regenerate this file with moq.
var _ api.DataStorer = &DataStorerMock{}

// DataStorerMock is a mock implementation of api.DataStorer.
//
//	func TestSomethingThatUsesDataStorer(t *testing.T) {
//
//		// make and configure a mocked api.DataStorer
//		mockedDataStorer := &DataStorerMock{
//			GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
//				panic("mock out the GetResources method")
//			},
//		}
//
//		// use mockedDataStorer in code that requires api.DataStorer
//		// and then make assertions.
//
//	}
type DataStorerMock struct {
	// GetResourcesFunc mocks the GetResources method.
	GetResourcesFunc func(ctx context.Context, options data.Options) (*models.Resources, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetResources holds details about calls to the GetResources method.
		GetResources []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Options is the options argument value.
			Options data.Options
		}
	}
	lockGetResources sync.RWMutex
}

// GetResources calls GetResourcesFunc.
func (mock *DataStorerMock) GetResources(ctx context.Context, options data.Options) (*models.Resources, error) {
	if mock.GetResourcesFunc == nil {
		panic("DataStorerMock.GetResourcesFunc: method is nil but DataStorer.GetResources was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Options data.Options
	}{
		Ctx:     ctx,
		Options: options,
	}
	mock.lockGetResources.Lock()
	mock.calls.GetResources = append(mock.calls.GetResources, callInfo)
	mock.lockGetResources.Unlock()
	return mock.GetResourcesFunc(ctx, options)
}

// GetResourcesCalls gets all the calls that were made to GetResources.
// Check the length with:
//
//	len(mockedDataStorer.GetResourcesCalls())
func (mock *DataStorerMock) GetResourcesCalls() []struct {
	Ctx     context.Context
	Options data.Options
} {
	var calls []struct {
		Ctx     context.Context
		Options data.Options
	}
	mock.lockGetResources.RLock()
	calls = mock.calls.GetResources
	mock.lockGetResources.RUnlock()
	return calls
}
