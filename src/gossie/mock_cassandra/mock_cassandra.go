// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/apesternikov/gossie/src/cassandra (interfaces: Cassandra)

package mock_cassandra

import (
	cassandra "github.com/apesternikov/gossie/src/cassandra"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Cassandra interface
type MockCassandra struct {
	ctrl     *gomock.Controller
	recorder *_MockCassandraRecorder
}

// Recorder for MockCassandra (not exported)
type _MockCassandraRecorder struct {
	mock *MockCassandra
}

func NewMockCassandra(ctrl *gomock.Controller) *MockCassandra {
	mock := &MockCassandra{ctrl: ctrl}
	mock.recorder = &_MockCassandraRecorder{mock}
	return mock
}

func (_m *MockCassandra) EXPECT() *_MockCassandraRecorder {
	return _m.recorder
}

func (_m *MockCassandra) Add(_param0 []byte, _param1 *cassandra.ColumnParent, _param2 *cassandra.CounterColumn, _param3 cassandra.ConsistencyLevel) error {
	ret := _m.ctrl.Call(_m, "Add", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) Add(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Add", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) AtomicBatchMutate(_param0 map[string]map[string][]*cassandra.Mutation, _param1 cassandra.ConsistencyLevel) error {
	ret := _m.ctrl.Call(_m, "AtomicBatchMutate", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) AtomicBatchMutate(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AtomicBatchMutate", arg0, arg1)
}

func (_m *MockCassandra) BatchMutate(_param0 map[string]map[string][]*cassandra.Mutation, _param1 cassandra.ConsistencyLevel) error {
	ret := _m.ctrl.Call(_m, "BatchMutate", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) BatchMutate(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "BatchMutate", arg0, arg1)
}

func (_m *MockCassandra) Cas(_param0 []byte, _param1 string, _param2 []*cassandra.Column, _param3 []*cassandra.Column, _param4 cassandra.ConsistencyLevel, _param5 cassandra.ConsistencyLevel) (*cassandra.CASResult_, error) {
	ret := _m.ctrl.Call(_m, "Cas", _param0, _param1, _param2, _param3, _param4, _param5)
	ret0, _ := ret[0].(*cassandra.CASResult_)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) Cas(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Cas", arg0, arg1, arg2, arg3, arg4, arg5)
}

func (_m *MockCassandra) DescribeClusterName() (string, error) {
	ret := _m.ctrl.Call(_m, "DescribeClusterName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeClusterName() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeClusterName")
}

func (_m *MockCassandra) DescribeKeyspace(_param0 string) (*cassandra.KsDef, error) {
	ret := _m.ctrl.Call(_m, "DescribeKeyspace", _param0)
	ret0, _ := ret[0].(*cassandra.KsDef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeKeyspace(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeKeyspace", arg0)
}

func (_m *MockCassandra) DescribeKeyspaces() ([]*cassandra.KsDef, error) {
	ret := _m.ctrl.Call(_m, "DescribeKeyspaces")
	ret0, _ := ret[0].([]*cassandra.KsDef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeKeyspaces() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeKeyspaces")
}

func (_m *MockCassandra) DescribeLocalRing(_param0 string) ([]*cassandra.TokenRange, error) {
	ret := _m.ctrl.Call(_m, "DescribeLocalRing", _param0)
	ret0, _ := ret[0].([]*cassandra.TokenRange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeLocalRing(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeLocalRing", arg0)
}

func (_m *MockCassandra) DescribePartitioner() (string, error) {
	ret := _m.ctrl.Call(_m, "DescribePartitioner")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribePartitioner() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribePartitioner")
}

func (_m *MockCassandra) DescribeRing(_param0 string) ([]*cassandra.TokenRange, error) {
	ret := _m.ctrl.Call(_m, "DescribeRing", _param0)
	ret0, _ := ret[0].([]*cassandra.TokenRange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeRing(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeRing", arg0)
}

func (_m *MockCassandra) DescribeSchemaVersions() (map[string][]string, error) {
	ret := _m.ctrl.Call(_m, "DescribeSchemaVersions")
	ret0, _ := ret[0].(map[string][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeSchemaVersions() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeSchemaVersions")
}

func (_m *MockCassandra) DescribeSnitch() (string, error) {
	ret := _m.ctrl.Call(_m, "DescribeSnitch")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeSnitch() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeSnitch")
}

func (_m *MockCassandra) DescribeSplits(_param0 string, _param1 string, _param2 string, _param3 int32) ([]string, error) {
	ret := _m.ctrl.Call(_m, "DescribeSplits", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeSplits(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeSplits", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) DescribeSplitsEx(_param0 string, _param1 string, _param2 string, _param3 int32) ([]*cassandra.CfSplit, error) {
	ret := _m.ctrl.Call(_m, "DescribeSplitsEx", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]*cassandra.CfSplit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeSplitsEx(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeSplitsEx", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) DescribeTokenMap() (map[string]string, error) {
	ret := _m.ctrl.Call(_m, "DescribeTokenMap")
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeTokenMap() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeTokenMap")
}

func (_m *MockCassandra) DescribeVersion() (string, error) {
	ret := _m.ctrl.Call(_m, "DescribeVersion")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) DescribeVersion() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeVersion")
}

func (_m *MockCassandra) ExecuteCql3Query(_param0 []byte, _param1 cassandra.Compression, _param2 cassandra.ConsistencyLevel) (*cassandra.CqlResult_, error) {
	ret := _m.ctrl.Call(_m, "ExecuteCql3Query", _param0, _param1, _param2)
	ret0, _ := ret[0].(*cassandra.CqlResult_)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) ExecuteCql3Query(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ExecuteCql3Query", arg0, arg1, arg2)
}

func (_m *MockCassandra) ExecuteCqlQuery(_param0 []byte, _param1 cassandra.Compression) (*cassandra.CqlResult_, error) {
	ret := _m.ctrl.Call(_m, "ExecuteCqlQuery", _param0, _param1)
	ret0, _ := ret[0].(*cassandra.CqlResult_)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) ExecuteCqlQuery(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ExecuteCqlQuery", arg0, arg1)
}

func (_m *MockCassandra) ExecutePreparedCql3Query(_param0 int32, _param1 [][]byte, _param2 cassandra.ConsistencyLevel) (*cassandra.CqlResult_, error) {
	ret := _m.ctrl.Call(_m, "ExecutePreparedCql3Query", _param0, _param1, _param2)
	ret0, _ := ret[0].(*cassandra.CqlResult_)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) ExecutePreparedCql3Query(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ExecutePreparedCql3Query", arg0, arg1, arg2)
}

func (_m *MockCassandra) ExecutePreparedCqlQuery(_param0 int32, _param1 [][]byte) (*cassandra.CqlResult_, error) {
	ret := _m.ctrl.Call(_m, "ExecutePreparedCqlQuery", _param0, _param1)
	ret0, _ := ret[0].(*cassandra.CqlResult_)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) ExecutePreparedCqlQuery(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ExecutePreparedCqlQuery", arg0, arg1)
}

func (_m *MockCassandra) Get(_param0 []byte, _param1 *cassandra.ColumnPath, _param2 cassandra.ConsistencyLevel) (*cassandra.ColumnOrSuperColumn, error) {
	ret := _m.ctrl.Call(_m, "Get", _param0, _param1, _param2)
	ret0, _ := ret[0].(*cassandra.ColumnOrSuperColumn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0, arg1, arg2)
}

func (_m *MockCassandra) GetCount(_param0 []byte, _param1 *cassandra.ColumnParent, _param2 *cassandra.SlicePredicate, _param3 cassandra.ConsistencyLevel) (int32, error) {
	ret := _m.ctrl.Call(_m, "GetCount", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) GetCount(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetCount", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) GetIndexedSlices(_param0 *cassandra.ColumnParent, _param1 *cassandra.IndexClause, _param2 *cassandra.SlicePredicate, _param3 cassandra.ConsistencyLevel) ([]*cassandra.KeySlice, error) {
	ret := _m.ctrl.Call(_m, "GetIndexedSlices", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]*cassandra.KeySlice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) GetIndexedSlices(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetIndexedSlices", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) GetPagedSlice(_param0 string, _param1 *cassandra.KeyRange, _param2 []byte, _param3 cassandra.ConsistencyLevel) ([]*cassandra.KeySlice, error) {
	ret := _m.ctrl.Call(_m, "GetPagedSlice", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]*cassandra.KeySlice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) GetPagedSlice(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetPagedSlice", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) GetRangeSlices(_param0 *cassandra.ColumnParent, _param1 *cassandra.SlicePredicate, _param2 *cassandra.KeyRange, _param3 cassandra.ConsistencyLevel) ([]*cassandra.KeySlice, error) {
	ret := _m.ctrl.Call(_m, "GetRangeSlices", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]*cassandra.KeySlice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) GetRangeSlices(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRangeSlices", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) GetSlice(_param0 []byte, _param1 *cassandra.ColumnParent, _param2 *cassandra.SlicePredicate, _param3 cassandra.ConsistencyLevel) ([]*cassandra.ColumnOrSuperColumn, error) {
	ret := _m.ctrl.Call(_m, "GetSlice", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]*cassandra.ColumnOrSuperColumn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) GetSlice(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSlice", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) Insert(_param0 []byte, _param1 *cassandra.ColumnParent, _param2 *cassandra.Column, _param3 cassandra.ConsistencyLevel) error {
	ret := _m.ctrl.Call(_m, "Insert", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) Insert(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Insert", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) Login(_param0 *cassandra.AuthenticationRequest) error {
	ret := _m.ctrl.Call(_m, "Login", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) Login(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Login", arg0)
}

func (_m *MockCassandra) MultigetCount(_param0 [][]byte, _param1 *cassandra.ColumnParent, _param2 *cassandra.SlicePredicate, _param3 cassandra.ConsistencyLevel) (map[string]int32, error) {
	ret := _m.ctrl.Call(_m, "MultigetCount", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(map[string]int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) MultigetCount(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MultigetCount", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) MultigetSlice(_param0 [][]byte, _param1 *cassandra.ColumnParent, _param2 *cassandra.SlicePredicate, _param3 cassandra.ConsistencyLevel) (map[string][]*cassandra.ColumnOrSuperColumn, error) {
	ret := _m.ctrl.Call(_m, "MultigetSlice", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(map[string][]*cassandra.ColumnOrSuperColumn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) MultigetSlice(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MultigetSlice", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) PrepareCql3Query(_param0 []byte, _param1 cassandra.Compression) (*cassandra.CqlPreparedResult_, error) {
	ret := _m.ctrl.Call(_m, "PrepareCql3Query", _param0, _param1)
	ret0, _ := ret[0].(*cassandra.CqlPreparedResult_)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) PrepareCql3Query(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PrepareCql3Query", arg0, arg1)
}

func (_m *MockCassandra) PrepareCqlQuery(_param0 []byte, _param1 cassandra.Compression) (*cassandra.CqlPreparedResult_, error) {
	ret := _m.ctrl.Call(_m, "PrepareCqlQuery", _param0, _param1)
	ret0, _ := ret[0].(*cassandra.CqlPreparedResult_)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) PrepareCqlQuery(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PrepareCqlQuery", arg0, arg1)
}

func (_m *MockCassandra) Remove(_param0 []byte, _param1 *cassandra.ColumnPath, _param2 int64, _param3 cassandra.ConsistencyLevel) error {
	ret := _m.ctrl.Call(_m, "Remove", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) Remove(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Remove", arg0, arg1, arg2, arg3)
}

func (_m *MockCassandra) RemoveCounter(_param0 []byte, _param1 *cassandra.ColumnPath, _param2 cassandra.ConsistencyLevel) error {
	ret := _m.ctrl.Call(_m, "RemoveCounter", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) RemoveCounter(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoveCounter", arg0, arg1, arg2)
}

func (_m *MockCassandra) SetCqlVersion(_param0 string) error {
	ret := _m.ctrl.Call(_m, "SetCqlVersion", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) SetCqlVersion(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetCqlVersion", arg0)
}

func (_m *MockCassandra) SetKeyspace(_param0 string) error {
	ret := _m.ctrl.Call(_m, "SetKeyspace", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) SetKeyspace(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetKeyspace", arg0)
}

func (_m *MockCassandra) SystemAddColumnFamily(_param0 *cassandra.CfDef) (string, error) {
	ret := _m.ctrl.Call(_m, "SystemAddColumnFamily", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) SystemAddColumnFamily(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SystemAddColumnFamily", arg0)
}

func (_m *MockCassandra) SystemAddKeyspace(_param0 *cassandra.KsDef) (string, error) {
	ret := _m.ctrl.Call(_m, "SystemAddKeyspace", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) SystemAddKeyspace(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SystemAddKeyspace", arg0)
}

func (_m *MockCassandra) SystemDropColumnFamily(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "SystemDropColumnFamily", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) SystemDropColumnFamily(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SystemDropColumnFamily", arg0)
}

func (_m *MockCassandra) SystemDropKeyspace(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "SystemDropKeyspace", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) SystemDropKeyspace(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SystemDropKeyspace", arg0)
}

func (_m *MockCassandra) SystemUpdateColumnFamily(_param0 *cassandra.CfDef) (string, error) {
	ret := _m.ctrl.Call(_m, "SystemUpdateColumnFamily", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) SystemUpdateColumnFamily(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SystemUpdateColumnFamily", arg0)
}

func (_m *MockCassandra) SystemUpdateKeyspace(_param0 *cassandra.KsDef) (string, error) {
	ret := _m.ctrl.Call(_m, "SystemUpdateKeyspace", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) SystemUpdateKeyspace(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SystemUpdateKeyspace", arg0)
}

func (_m *MockCassandra) TraceNextQuery() ([]byte, error) {
	ret := _m.ctrl.Call(_m, "TraceNextQuery")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCassandraRecorder) TraceNextQuery() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "TraceNextQuery")
}

func (_m *MockCassandra) Truncate(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Truncate", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCassandraRecorder) Truncate(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Truncate", arg0)
}
