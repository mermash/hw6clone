// Code generated by MockGen. DO NOT EDIT.
// Source: posts_handlers.go

// Package main is a generated GoMock package.
package main

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPostRepoI is a mock of PostRepoI interface.
type MockPostRepoI struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepoIMockRecorder
}

// MockPostRepoIMockRecorder is the mock recorder for MockPostRepoI.
type MockPostRepoIMockRecorder struct {
	mock *MockPostRepoI
}

// NewMockPostRepoI creates a new mock instance.
func NewMockPostRepoI(ctrl *gomock.Controller) *MockPostRepoI {
	mock := &MockPostRepoI{ctrl: ctrl}
	mock.recorder = &MockPostRepoIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostRepoI) EXPECT() *MockPostRepoIMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockPostRepoI) Add(post *Post) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", post)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockPostRepoIMockRecorder) Add(post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockPostRepoI)(nil).Add), post)
}

// Delete mocks base method.
func (m *MockPostRepoI) Delete(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockPostRepoIMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPostRepoI)(nil).Delete), id)
}

// DownVote mocks base method.
func (m *MockPostRepoI) DownVote(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownVote", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownVote indicates an expected call of DownVote.
func (mr *MockPostRepoIMockRecorder) DownVote(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownVote", reflect.TypeOf((*MockPostRepoI)(nil).DownVote), id)
}

// GetAll mocks base method.
func (m *MockPostRepoI) GetAll() ([]*PostComplexData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*PostComplexData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPostRepoIMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPostRepoI)(nil).GetAll))
}

// GetByCategoryName mocks base method.
func (m *MockPostRepoI) GetByCategoryName(categoryName string) ([]*PostComplexData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCategoryName", categoryName)
	ret0, _ := ret[0].([]*PostComplexData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCategoryName indicates an expected call of GetByCategoryName.
func (mr *MockPostRepoIMockRecorder) GetByCategoryName(categoryName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCategoryName", reflect.TypeOf((*MockPostRepoI)(nil).GetByCategoryName), categoryName)
}

// GetById mocks base method.
func (m *MockPostRepoI) GetById(id string) (*PostComplexData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(*PostComplexData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockPostRepoIMockRecorder) GetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockPostRepoI)(nil).GetById), id)
}

// GetByUserLogin mocks base method.
func (m *MockPostRepoI) GetByUserLogin(userLogin string) ([]*PostComplexData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserLogin", userLogin)
	ret0, _ := ret[0].([]*PostComplexData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserLogin indicates an expected call of GetByUserLogin.
func (mr *MockPostRepoIMockRecorder) GetByUserLogin(userLogin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserLogin", reflect.TypeOf((*MockPostRepoI)(nil).GetByUserLogin), userLogin)
}

// UpVote mocks base method.
func (m *MockPostRepoI) UpVote(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpVote", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpVote indicates an expected call of UpVote.
func (mr *MockPostRepoIMockRecorder) UpVote(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpVote", reflect.TypeOf((*MockPostRepoI)(nil).UpVote), id)
}

// MockCommentRepoI is a mock of CommentRepoI interface.
type MockCommentRepoI struct {
	ctrl     *gomock.Controller
	recorder *MockCommentRepoIMockRecorder
}

// MockCommentRepoIMockRecorder is the mock recorder for MockCommentRepoI.
type MockCommentRepoIMockRecorder struct {
	mock *MockCommentRepoI
}

// NewMockCommentRepoI creates a new mock instance.
func NewMockCommentRepoI(ctrl *gomock.Controller) *MockCommentRepoI {
	mock := &MockCommentRepoI{ctrl: ctrl}
	mock.recorder = &MockCommentRepoIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentRepoI) EXPECT() *MockCommentRepoIMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCommentRepoI) Add(comment *Comment) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", comment)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockCommentRepoIMockRecorder) Add(comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCommentRepoI)(nil).Add), comment)
}

// Delete mocks base method.
func (m *MockCommentRepoI) Delete(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockCommentRepoIMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCommentRepoI)(nil).Delete), id)
}

// GetCommentsByPostIds mocks base method.
func (m *MockCommentRepoI) GetCommentsByPostIds(postIds []string) (map[string][]*CommentComplexData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentsByPostIds", postIds)
	ret0, _ := ret[0].(map[string][]*CommentComplexData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentsByPostIds indicates an expected call of GetCommentsByPostIds.
func (mr *MockCommentRepoIMockRecorder) GetCommentsByPostIds(postIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentsByPostIds", reflect.TypeOf((*MockCommentRepoI)(nil).GetCommentsByPostIds), postIds)
}

// MockVoteRepoI is a mock of VoteRepoI interface.
type MockVoteRepoI struct {
	ctrl     *gomock.Controller
	recorder *MockVoteRepoIMockRecorder
}

// MockVoteRepoIMockRecorder is the mock recorder for MockVoteRepoI.
type MockVoteRepoIMockRecorder struct {
	mock *MockVoteRepoI
}

// NewMockVoteRepoI creates a new mock instance.
func NewMockVoteRepoI(ctrl *gomock.Controller) *MockVoteRepoI {
	mock := &MockVoteRepoI{ctrl: ctrl}
	mock.recorder = &MockVoteRepoIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVoteRepoI) EXPECT() *MockVoteRepoIMockRecorder {
	return m.recorder
}

// GetVotesByPostIds mocks base method.
func (m *MockVoteRepoI) GetVotesByPostIds(postIds []string) (map[string][]*Vote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVotesByPostIds", postIds)
	ret0, _ := ret[0].(map[string][]*Vote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVotesByPostIds indicates an expected call of GetVotesByPostIds.
func (mr *MockVoteRepoIMockRecorder) GetVotesByPostIds(postIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVotesByPostIds", reflect.TypeOf((*MockVoteRepoI)(nil).GetVotesByPostIds), postIds)
}

// MockDictionaryRepoI is a mock of DictionaryRepoI interface.
type MockDictionaryRepoI struct {
	ctrl     *gomock.Controller
	recorder *MockDictionaryRepoIMockRecorder
}

// MockDictionaryRepoIMockRecorder is the mock recorder for MockDictionaryRepoI.
type MockDictionaryRepoIMockRecorder struct {
	mock *MockDictionaryRepoI
}

// NewMockDictionaryRepoI creates a new mock instance.
func NewMockDictionaryRepoI(ctrl *gomock.Controller) *MockDictionaryRepoI {
	mock := &MockDictionaryRepoI{ctrl: ctrl}
	mock.recorder = &MockDictionaryRepoIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDictionaryRepoI) EXPECT() *MockDictionaryRepoIMockRecorder {
	return m.recorder
}

// GetCategoryByName mocks base method.
func (m *MockDictionaryRepoI) GetCategoryByName(name string) (*Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryByName", name)
	ret0, _ := ret[0].(*Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryByName indicates an expected call of GetCategoryByName.
func (mr *MockDictionaryRepoIMockRecorder) GetCategoryByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryByName", reflect.TypeOf((*MockDictionaryRepoI)(nil).GetCategoryByName), name)
}

// MockDTOConverterI is a mock of DTOConverterI interface.
type MockDTOConverterI struct {
	ctrl     *gomock.Controller
	recorder *MockDTOConverterIMockRecorder
}

// MockDTOConverterIMockRecorder is the mock recorder for MockDTOConverterI.
type MockDTOConverterIMockRecorder struct {
	mock *MockDTOConverterI
}

// NewMockDTOConverterI creates a new mock instance.
func NewMockDTOConverterI(ctrl *gomock.Controller) *MockDTOConverterI {
	mock := &MockDTOConverterI{ctrl: ctrl}
	mock.recorder = &MockDTOConverterIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDTOConverterI) EXPECT() *MockDTOConverterIMockRecorder {
	return m.recorder
}

// CommentsConvertToDTO mocks base method.
func (m *MockDTOConverterI) CommentsConvertToDTO(data []*CommentComplexData) []*CommentDTO {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommentsConvertToDTO", data)
	ret0, _ := ret[0].([]*CommentDTO)
	return ret0
}

// CommentsConvertToDTO indicates an expected call of CommentsConvertToDTO.
func (mr *MockDTOConverterIMockRecorder) CommentsConvertToDTO(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommentsConvertToDTO", reflect.TypeOf((*MockDTOConverterI)(nil).CommentsConvertToDTO), data)
}

// PostConvertToDTO mocks base method.
func (m *MockDTOConverterI) PostConvertToDTO(data *PostComplexData) (*PostDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostConvertToDTO", data)
	ret0, _ := ret[0].(*PostDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostConvertToDTO indicates an expected call of PostConvertToDTO.
func (mr *MockDTOConverterIMockRecorder) PostConvertToDTO(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostConvertToDTO", reflect.TypeOf((*MockDTOConverterI)(nil).PostConvertToDTO), data)
}

// PostsConvertToDTO mocks base method.
func (m *MockDTOConverterI) PostsConvertToDTO(data []*PostComplexData) ([]*PostDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostsConvertToDTO", data)
	ret0, _ := ret[0].([]*PostDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostsConvertToDTO indicates an expected call of PostsConvertToDTO.
func (mr *MockDTOConverterIMockRecorder) PostsConvertToDTO(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostsConvertToDTO", reflect.TypeOf((*MockDTOConverterI)(nil).PostsConvertToDTO), data)
}

// VotesConvertToDTO mocks base method.
func (m *MockDTOConverterI) VotesConvertToDTO(data []*Vote) []*VoteDTO {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VotesConvertToDTO", data)
	ret0, _ := ret[0].([]*VoteDTO)
	return ret0
}

// VotesConvertToDTO indicates an expected call of VotesConvertToDTO.
func (mr *MockDTOConverterIMockRecorder) VotesConvertToDTO(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VotesConvertToDTO", reflect.TypeOf((*MockDTOConverterI)(nil).VotesConvertToDTO), data)
}

// MockTimeGetterI is a mock of TimeGetterI interface.
type MockTimeGetterI struct {
	ctrl     *gomock.Controller
	recorder *MockTimeGetterIMockRecorder
}

// MockTimeGetterIMockRecorder is the mock recorder for MockTimeGetterI.
type MockTimeGetterIMockRecorder struct {
	mock *MockTimeGetterI
}

// NewMockTimeGetterI creates a new mock instance.
func NewMockTimeGetterI(ctrl *gomock.Controller) *MockTimeGetterI {
	mock := &MockTimeGetterI{ctrl: ctrl}
	mock.recorder = &MockTimeGetterIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimeGetterI) EXPECT() *MockTimeGetterIMockRecorder {
	return m.recorder
}

// GetCreated mocks base method.
func (m *MockTimeGetterI) GetCreated() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCreated")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCreated indicates an expected call of GetCreated.
func (mr *MockTimeGetterIMockRecorder) GetCreated() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCreated", reflect.TypeOf((*MockTimeGetterI)(nil).GetCreated))
}

// MockUUIDGetterI is a mock of UUIDGetterI interface.
type MockUUIDGetterI struct {
	ctrl     *gomock.Controller
	recorder *MockUUIDGetterIMockRecorder
}

// MockUUIDGetterIMockRecorder is the mock recorder for MockUUIDGetterI.
type MockUUIDGetterIMockRecorder struct {
	mock *MockUUIDGetterI
}

// NewMockUUIDGetterI creates a new mock instance.
func NewMockUUIDGetterI(ctrl *gomock.Controller) *MockUUIDGetterI {
	mock := &MockUUIDGetterI{ctrl: ctrl}
	mock.recorder = &MockUUIDGetterIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUUIDGetterI) EXPECT() *MockUUIDGetterIMockRecorder {
	return m.recorder
}

// GetUUID mocks base method.
func (m *MockUUIDGetterI) GetUUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetUUID indicates an expected call of GetUUID.
func (mr *MockUUIDGetterIMockRecorder) GetUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUUID", reflect.TypeOf((*MockUUIDGetterI)(nil).GetUUID))
}
