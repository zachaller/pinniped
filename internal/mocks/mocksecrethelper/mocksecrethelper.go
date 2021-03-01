// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

// Code generated by MockGen. DO NOT EDIT.
// Source: go.pinniped.dev/internal/controller/supervisorconfig/generator (interfaces: SecretHelper)

// Package mocksecrethelper is a generated GoMock package.
package mocksecrethelper

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "go.pinniped.dev/generated/latest/apis/supervisor/config/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v10 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockSecretHelper is a mock of SecretHelper interface.
type MockSecretHelper struct {
	ctrl     *gomock.Controller
	recorder *MockSecretHelperMockRecorder
}

// MockSecretHelperMockRecorder is the mock recorder for MockSecretHelper.
type MockSecretHelperMockRecorder struct {
	mock *MockSecretHelper
}

// NewMockSecretHelper creates a new mock instance.
func NewMockSecretHelper(ctrl *gomock.Controller) *MockSecretHelper {
	mock := &MockSecretHelper{ctrl: ctrl}
	mock.recorder = &MockSecretHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretHelper) EXPECT() *MockSecretHelperMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockSecretHelper) Generate(arg0 *v1alpha1.FederationDomain) (*v1.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", arg0)
	ret0, _ := ret[0].(*v1.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockSecretHelperMockRecorder) Generate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockSecretHelper)(nil).Generate), arg0)
}

// Handles mocks base method.
func (m *MockSecretHelper) Handles(arg0 v10.Object) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handles", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Handles indicates an expected call of Handles.
func (mr *MockSecretHelperMockRecorder) Handles(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handles", reflect.TypeOf((*MockSecretHelper)(nil).Handles), arg0)
}

// IsValid mocks base method.
func (m *MockSecretHelper) IsValid(arg0 *v1alpha1.FederationDomain, arg1 *v1.Secret) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValid", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValid indicates an expected call of IsValid.
func (mr *MockSecretHelperMockRecorder) IsValid(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValid", reflect.TypeOf((*MockSecretHelper)(nil).IsValid), arg0, arg1)
}

// NamePrefix mocks base method.
func (m *MockSecretHelper) NamePrefix() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamePrefix")
	ret0, _ := ret[0].(string)
	return ret0
}

// NamePrefix indicates an expected call of NamePrefix.
func (mr *MockSecretHelperMockRecorder) NamePrefix() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamePrefix", reflect.TypeOf((*MockSecretHelper)(nil).NamePrefix))
}

// ObserveActiveSecretAndUpdateParentFederationDomain mocks base method.
func (m *MockSecretHelper) ObserveActiveSecretAndUpdateParentFederationDomain(arg0 *v1alpha1.FederationDomain, arg1 *v1.Secret) *v1alpha1.FederationDomain {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObserveActiveSecretAndUpdateParentFederationDomain", arg0, arg1)
	ret0, _ := ret[0].(*v1alpha1.FederationDomain)
	return ret0
}

// ObserveActiveSecretAndUpdateParentFederationDomain indicates an expected call of ObserveActiveSecretAndUpdateParentFederationDomain.
func (mr *MockSecretHelperMockRecorder) ObserveActiveSecretAndUpdateParentFederationDomain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObserveActiveSecretAndUpdateParentFederationDomain", reflect.TypeOf((*MockSecretHelper)(nil).ObserveActiveSecretAndUpdateParentFederationDomain), arg0, arg1)
}
