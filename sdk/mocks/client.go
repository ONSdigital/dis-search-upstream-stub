// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/dis-search-upstream-stub/sdk"
	apiError "github.com/ONSdigital/dis-search-upstream-stub/sdk/errors"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"sync"
)

// Ensure, that ClienterMock does implement sdk.Clienter.
// If this is not the case, regenerate this file with moq.
var _ sdk.Clienter = &ClienterMock{}

// ClienterMock is a mock implementation of sdk.Clienter.
//
//	func TestSomethingThatUsesClienter(t *testing.T) {
//
//		// make and configure a mocked sdk.Clienter
//		mockedClienter := &ClienterMock{
//			CheckerFunc: func(ctx context.Context, check *health.CheckState) error {
//				panic("mock out the Checker method")
//			},
//			GetResourcesFunc: func(ctx context.Context, options sdk.Options) (*models.Resources, apiError.Error) {
//				panic("mock out the GetResources method")
//			},
//		}
//
//		// use mockedClienter in code that requires sdk.Clienter
//		// and then make assertions.
//
//	}
type ClienterMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(ctx context.Context, check *health.CheckState) error

	// GetResourcesFunc mocks the GetResources method.
	GetResourcesFunc func(ctx context.Context, options sdk.Options) (*models.Resources, apiError.Error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Check is the check argument value.
			Check *health.CheckState
		}
		// GetResources holds details about calls to the GetResources method.
		GetResources []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Options is the options argument value.
			Options sdk.Options
		}
	}
	lockChecker      sync.RWMutex
	lockGetResources sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *ClienterMock) Checker(ctx context.Context, check *health.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("ClienterMock.CheckerFunc: method is nil but Clienter.Checker was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Check *health.CheckState
	}{
		Ctx:   ctx,
		Check: check,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(ctx, check)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//
//	len(mockedClienter.CheckerCalls())
func (mock *ClienterMock) CheckerCalls() []struct {
	Ctx   context.Context
	Check *health.CheckState
} {
	var calls []struct {
		Ctx   context.Context
		Check *health.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// GetResources calls GetResourcesFunc.
func (mock *ClienterMock) GetResources(ctx context.Context, options sdk.Options) (*models.Resources, apiError.Error) {
	if mock.GetResourcesFunc == nil {
		panic("ClienterMock.GetResourcesFunc: method is nil but Clienter.GetResources was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Options sdk.Options
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
//	len(mockedClienter.GetResourcesCalls())
func (mock *ClienterMock) GetResourcesCalls() []struct {
	Ctx     context.Context
	Options sdk.Options
} {
	var calls []struct {
		Ctx     context.Context
		Options sdk.Options
	}
	mock.lockGetResources.RLock()
	calls = mock.calls.GetResources
	mock.lockGetResources.RUnlock()
	return calls
}