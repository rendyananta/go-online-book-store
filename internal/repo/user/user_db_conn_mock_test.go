// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package user is a generated GoMock package.
package user

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
)

// MockdbConnManager is a mock of dbConnManager interface.
type MockdbConnManager struct {
	ctrl     *gomock.Controller
	recorder *MockdbConnManagerMockRecorder
}

// MockdbConnManagerMockRecorder is the mock recorder for MockdbConnManager.
type MockdbConnManagerMockRecorder struct {
	mock *MockdbConnManager
}

// NewMockdbConnManager creates a new mock instance.
func NewMockdbConnManager(ctrl *gomock.Controller) *MockdbConnManager {
	mock := &MockdbConnManager{ctrl: ctrl}
	mock.recorder = &MockdbConnManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockdbConnManager) EXPECT() *MockdbConnManagerMockRecorder {
	return m.recorder
}

// Connection mocks base method.
func (m *MockdbConnManager) Connection(name string) (*sqlx.DB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connection", name)
	ret0, _ := ret[0].(*sqlx.DB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connection indicates an expected call of Connection.
func (mr *MockdbConnManagerMockRecorder) Connection(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connection", reflect.TypeOf((*MockdbConnManager)(nil).Connection), name)
}

// MockdbConnection is a mock of dbConnection interface.
type MockdbConnection struct {
	ctrl     *gomock.Controller
	recorder *MockdbConnectionMockRecorder
}

// MockdbConnectionMockRecorder is the mock recorder for MockdbConnection.
type MockdbConnectionMockRecorder struct {
	mock *MockdbConnection
}

// NewMockdbConnection creates a new mock instance.
func NewMockdbConnection(ctrl *gomock.Controller) *MockdbConnection {
	mock := &MockdbConnection{ctrl: ctrl}
	mock.recorder = &MockdbConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockdbConnection) EXPECT() *MockdbConnectionMockRecorder {
	return m.recorder
}

// ExecContext mocks base method.
func (m *MockdbConnection) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockdbConnectionMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockdbConnection)(nil).ExecContext), varargs...)
}

// Preparex mocks base method.
func (m *MockdbConnection) Preparex(query string) (*sqlx.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Preparex", query)
	ret0, _ := ret[0].(*sqlx.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Preparex indicates an expected call of Preparex.
func (mr *MockdbConnectionMockRecorder) Preparex(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Preparex", reflect.TypeOf((*MockdbConnection)(nil).Preparex), query)
}

// Rebind mocks base method.
func (m *MockdbConnection) Rebind(query string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rebind", query)
	ret0, _ := ret[0].(string)
	return ret0
}

// Rebind indicates an expected call of Rebind.
func (mr *MockdbConnectionMockRecorder) Rebind(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rebind", reflect.TypeOf((*MockdbConnection)(nil).Rebind), query)
}

// MockqueryGetter is a mock of queryGetter interface.
type MockqueryGetter struct {
	ctrl     *gomock.Controller
	recorder *MockqueryGetterMockRecorder
}

// MockqueryGetterMockRecorder is the mock recorder for MockqueryGetter.
type MockqueryGetterMockRecorder struct {
	mock *MockqueryGetter
}

// NewMockqueryGetter creates a new mock instance.
func NewMockqueryGetter(ctrl *gomock.Controller) *MockqueryGetter {
	mock := &MockqueryGetter{ctrl: ctrl}
	mock.recorder = &MockqueryGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockqueryGetter) EXPECT() *MockqueryGetterMockRecorder {
	return m.recorder
}

// GetContext mocks base method.
func (m *MockqueryGetter) GetContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetContext indicates an expected call of GetContext.
func (mr *MockqueryGetterMockRecorder) GetContext(ctx, dest interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContext", reflect.TypeOf((*MockqueryGetter)(nil).GetContext), varargs...)
}

// MockqueryExecer is a mock of queryExecer interface.
type MockqueryExecer struct {
	ctrl     *gomock.Controller
	recorder *MockqueryExecerMockRecorder
}

// MockqueryExecerMockRecorder is the mock recorder for MockqueryExecer.
type MockqueryExecerMockRecorder struct {
	mock *MockqueryExecer
}

// NewMockqueryExecer creates a new mock instance.
func NewMockqueryExecer(ctrl *gomock.Controller) *MockqueryExecer {
	mock := &MockqueryExecer{ctrl: ctrl}
	mock.recorder = &MockqueryExecerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockqueryExecer) EXPECT() *MockqueryExecerMockRecorder {
	return m.recorder
}

// ExecContext mocks base method.
func (m *MockqueryExecer) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockqueryExecerMockRecorder) ExecContext(ctx interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockqueryExecer)(nil).ExecContext), varargs...)
}
