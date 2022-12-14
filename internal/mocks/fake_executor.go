// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"context"
	"sync"

	"github.com/ryantate13/klogs/exec"
)

type FakeExecutor struct {
	StreamStub        func(context.Context, chan<- error, ...string) (<-chan string, error)
	streamMutex       sync.RWMutex
	streamArgsForCall []struct {
		arg1 context.Context
		arg2 chan<- error
		arg3 []string
	}
	streamReturns struct {
		result1 <-chan string
		result2 error
	}
	streamReturnsOnCall map[int]struct {
		result1 <-chan string
		result2 error
	}
	SyncStub        func(context.Context, ...string) ([]string, error)
	syncMutex       sync.RWMutex
	syncArgsForCall []struct {
		arg1 context.Context
		arg2 []string
	}
	syncReturns struct {
		result1 []string
		result2 error
	}
	syncReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeExecutor) Stream(arg1 context.Context, arg2 chan<- error, arg3 ...string) (<-chan string, error) {
	fake.streamMutex.Lock()
	ret, specificReturn := fake.streamReturnsOnCall[len(fake.streamArgsForCall)]
	fake.streamArgsForCall = append(fake.streamArgsForCall, struct {
		arg1 context.Context
		arg2 chan<- error
		arg3 []string
	}{arg1, arg2, arg3})
	stub := fake.StreamStub
	fakeReturns := fake.streamReturns
	fake.recordInvocation("Stream", []interface{}{arg1, arg2, arg3})
	fake.streamMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeExecutor) StreamCallCount() int {
	fake.streamMutex.RLock()
	defer fake.streamMutex.RUnlock()
	return len(fake.streamArgsForCall)
}

func (fake *FakeExecutor) StreamCalls(stub func(context.Context, chan<- error, ...string) (<-chan string, error)) {
	fake.streamMutex.Lock()
	defer fake.streamMutex.Unlock()
	fake.StreamStub = stub
}

func (fake *FakeExecutor) StreamArgsForCall(i int) (context.Context, chan<- error, []string) {
	fake.streamMutex.RLock()
	defer fake.streamMutex.RUnlock()
	argsForCall := fake.streamArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeExecutor) StreamReturns(result1 <-chan string, result2 error) {
	fake.streamMutex.Lock()
	defer fake.streamMutex.Unlock()
	fake.StreamStub = nil
	fake.streamReturns = struct {
		result1 <-chan string
		result2 error
	}{result1, result2}
}

func (fake *FakeExecutor) StreamReturnsOnCall(i int, result1 <-chan string, result2 error) {
	fake.streamMutex.Lock()
	defer fake.streamMutex.Unlock()
	fake.StreamStub = nil
	if fake.streamReturnsOnCall == nil {
		fake.streamReturnsOnCall = make(map[int]struct {
			result1 <-chan string
			result2 error
		})
	}
	fake.streamReturnsOnCall[i] = struct {
		result1 <-chan string
		result2 error
	}{result1, result2}
}

func (fake *FakeExecutor) Sync(arg1 context.Context, arg2 ...string) ([]string, error) {
	fake.syncMutex.Lock()
	ret, specificReturn := fake.syncReturnsOnCall[len(fake.syncArgsForCall)]
	fake.syncArgsForCall = append(fake.syncArgsForCall, struct {
		arg1 context.Context
		arg2 []string
	}{arg1, arg2})
	stub := fake.SyncStub
	fakeReturns := fake.syncReturns
	fake.recordInvocation("Sync", []interface{}{arg1, arg2})
	fake.syncMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeExecutor) SyncCallCount() int {
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	return len(fake.syncArgsForCall)
}

func (fake *FakeExecutor) SyncCalls(stub func(context.Context, ...string) ([]string, error)) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
	fake.SyncStub = stub
}

func (fake *FakeExecutor) SyncArgsForCall(i int) (context.Context, []string) {
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	argsForCall := fake.syncArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeExecutor) SyncReturns(result1 []string, result2 error) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
	fake.SyncStub = nil
	fake.syncReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeExecutor) SyncReturnsOnCall(i int, result1 []string, result2 error) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
	fake.SyncStub = nil
	if fake.syncReturnsOnCall == nil {
		fake.syncReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.syncReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeExecutor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.streamMutex.RLock()
	defer fake.streamMutex.RUnlock()
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeExecutor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ exec.Executor = new(FakeExecutor)
