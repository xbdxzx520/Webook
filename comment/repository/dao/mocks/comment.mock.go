// Code generated by MockGen. DO NOT EDIT.
// Source: ./comment.go
//
// Generated by this command:
//
//	mockgen -source=./comment.go -package=daomocks -destination=mocks/comment.mock.go CommentDAO
//
// Package daomocks is a generated GoMock package.
package daomocks

import (
	context "context"
	reflect "reflect"

	dao "gitee.com/geekbang/basic-go/webook/comment/repository/dao"
	gomock "go.uber.org/mock/gomock"
)

// MockCommentDAO is a mock of CommentDAO interface.
type MockCommentDAO struct {
	ctrl     *gomock.Controller
	recorder *MockCommentDAOMockRecorder
}

// MockCommentDAOMockRecorder is the mock recorder for MockCommentDAO.
type MockCommentDAOMockRecorder struct {
	mock *MockCommentDAO
}

// NewMockCommentDAO creates a new mock instance.
func NewMockCommentDAO(ctrl *gomock.Controller) *MockCommentDAO {
	mock := &MockCommentDAO{ctrl: ctrl}
	mock.recorder = &MockCommentDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentDAO) EXPECT() *MockCommentDAOMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockCommentDAO) Delete(ctx context.Context, u dao.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCommentDAOMockRecorder) Delete(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCommentDAO)(nil).Delete), ctx, u)
}

// FindByBiz mocks base method.
func (m *MockCommentDAO) FindByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]dao.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByBiz", ctx, biz, bizId, minID, limit)
	ret0, _ := ret[0].([]dao.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByBiz indicates an expected call of FindByBiz.
func (mr *MockCommentDAOMockRecorder) FindByBiz(ctx, biz, bizId, minID, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByBiz", reflect.TypeOf((*MockCommentDAO)(nil).FindByBiz), ctx, biz, bizId, minID, limit)
}

// FindCommentList mocks base method.
func (m *MockCommentDAO) FindCommentList(ctx context.Context, u dao.Comment) ([]dao.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCommentList", ctx, u)
	ret0, _ := ret[0].([]dao.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCommentList indicates an expected call of FindCommentList.
func (mr *MockCommentDAOMockRecorder) FindCommentList(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCommentList", reflect.TypeOf((*MockCommentDAO)(nil).FindCommentList), ctx, u)
}

// FindOneByIDs mocks base method.
func (m *MockCommentDAO) FindOneByIDs(ctx context.Context, id []int64) ([]dao.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByIDs", ctx, id)
	ret0, _ := ret[0].([]dao.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByIDs indicates an expected call of FindOneByIDs.
func (mr *MockCommentDAOMockRecorder) FindOneByIDs(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByIDs", reflect.TypeOf((*MockCommentDAO)(nil).FindOneByIDs), ctx, id)
}

// FindRepliesByPid mocks base method.
func (m *MockCommentDAO) FindRepliesByPid(ctx context.Context, pid int64, offset, limit int) ([]dao.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRepliesByPid", ctx, pid, offset, limit)
	ret0, _ := ret[0].([]dao.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRepliesByPid indicates an expected call of FindRepliesByPid.
func (mr *MockCommentDAOMockRecorder) FindRepliesByPid(ctx, pid, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRepliesByPid", reflect.TypeOf((*MockCommentDAO)(nil).FindRepliesByPid), ctx, pid, offset, limit)
}

// FindRepliesByRid mocks base method.
func (m *MockCommentDAO) FindRepliesByRid(ctx context.Context, rid, id, limit int64) ([]dao.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRepliesByRid", ctx, rid, id, limit)
	ret0, _ := ret[0].([]dao.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRepliesByRid indicates an expected call of FindRepliesByRid.
func (mr *MockCommentDAOMockRecorder) FindRepliesByRid(ctx, rid, id, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRepliesByRid", reflect.TypeOf((*MockCommentDAO)(nil).FindRepliesByRid), ctx, rid, id, limit)
}

// Insert mocks base method.
func (m *MockCommentDAO) Insert(ctx context.Context, u dao.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCommentDAOMockRecorder) Insert(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockCommentDAO)(nil).Insert), ctx, u)
}
