// Code generated by MockGen. DO NOT EDIT.
// Source: cache.go
//
// Generated by this command:
//
//	mockgen -destination cache_mock.go -source cache.go -package cache
//
// Package cache is a generated GoMock package.
package cache

import (
	io "io"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCache) Add(k string, x any, d time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", k, x, d)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockCacheMockRecorder) Add(k, x, d any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCache)(nil).Add), k, x, d)
}

// Delete mocks base method.
func (m *MockCache) Delete(k string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", k)
}

// Delete indicates an expected call of Delete.
func (mr *MockCacheMockRecorder) Delete(k any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCache)(nil).Delete), k)
}

// DeleteExpired mocks base method.
func (m *MockCache) DeleteExpired() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteExpired")
}

// DeleteExpired indicates an expected call of DeleteExpired.
func (mr *MockCacheMockRecorder) DeleteExpired() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpired", reflect.TypeOf((*MockCache)(nil).DeleteExpired))
}

// Flush mocks base method.
func (m *MockCache) Flush() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Flush")
}

// Flush indicates an expected call of Flush.
func (mr *MockCacheMockRecorder) Flush() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockCache)(nil).Flush))
}

// Get mocks base method.
func (m *MockCache) Get(k string) (any, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", k)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder) Get(k any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), k)
}

// GetWithExpiration mocks base method.
func (m *MockCache) GetWithExpiration(k string) (any, time.Time, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithExpiration", k)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// GetWithExpiration indicates an expected call of GetWithExpiration.
func (mr *MockCacheMockRecorder) GetWithExpiration(k any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithExpiration", reflect.TypeOf((*MockCache)(nil).GetWithExpiration), k)
}

// ItemCount mocks base method.
func (m *MockCache) ItemCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// ItemCount indicates an expected call of ItemCount.
func (mr *MockCacheMockRecorder) ItemCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemCount", reflect.TypeOf((*MockCache)(nil).ItemCount))
}

// Items mocks base method.
func (m *MockCache) Items() map[string]Item {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Items")
	ret0, _ := ret[0].(map[string]Item)
	return ret0
}

// Items indicates an expected call of Items.
func (mr *MockCacheMockRecorder) Items() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Items", reflect.TypeOf((*MockCache)(nil).Items))
}

// Keys mocks base method.
func (m *MockCache) Keys() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Keys")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Keys indicates an expected call of Keys.
func (mr *MockCacheMockRecorder) Keys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Keys", reflect.TypeOf((*MockCache)(nil).Keys))
}

// Load mocks base method.
func (m *MockCache) Load(r io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", r)
	ret0, _ := ret[0].(error)
	return ret0
}

// Load indicates an expected call of Load.
func (mr *MockCacheMockRecorder) Load(r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockCache)(nil).Load), r)
}

// LoadFile mocks base method.
func (m *MockCache) LoadFile(fname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadFile", fname)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadFile indicates an expected call of LoadFile.
func (mr *MockCacheMockRecorder) LoadFile(fname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadFile", reflect.TypeOf((*MockCache)(nil).LoadFile), fname)
}

// OnEvicted mocks base method.
func (m *MockCache) OnEvicted(f func(string, any)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnEvicted", f)
}

// OnEvicted indicates an expected call of OnEvicted.
func (mr *MockCacheMockRecorder) OnEvicted(f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEvicted", reflect.TypeOf((*MockCache)(nil).OnEvicted), f)
}

// Save mocks base method.
func (m *MockCache) Save(w io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", w)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockCacheMockRecorder) Save(w any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCache)(nil).Save), w)
}

// SaveFile mocks base method.
func (m *MockCache) SaveFile(fname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFile", fname)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFile indicates an expected call of SaveFile.
func (mr *MockCacheMockRecorder) SaveFile(fname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFile", reflect.TypeOf((*MockCache)(nil).SaveFile), fname)
}

// Set mocks base method.
func (m *MockCache) Set(k string, x any, d time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", k, x, d)
}

// Set indicates an expected call of Set.
func (mr *MockCacheMockRecorder) Set(k, x, d any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCache)(nil).Set), k, x, d)
}

// SetDefault mocks base method.
func (m *MockCache) SetDefault(k string, x any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDefault", k, x)
}

// SetDefault indicates an expected call of SetDefault.
func (mr *MockCacheMockRecorder) SetDefault(k, x any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefault", reflect.TypeOf((*MockCache)(nil).SetDefault), k, x)
}
