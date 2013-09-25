package cassandra

import (
	"fmt"
	"github.com/pomack/thrift4go/lib/go/src/thrift"
	"math"
)

// This is a temporary safety measure to ensure that the `math'
// import does not trip up any generated output that may not
// happen to use the math import due to not having emited enums.
//
// Future clean-ups will deprecate the need for this.
func init() {
	var temporaryAndUnused int32 = math.MinInt32
	temporaryAndUnused++
}

/**
 *The ConsistencyLevel is an enum that controls both read and write
 *behavior based on the ReplicationFactor of the keyspace.  The
 *different consistency levels have different meanings, depending on
 *if you're doing a write or read operation.
 *
 *If W + R > ReplicationFactor, where W is the number of nodes to
 *block for on write, and R the number to block for on reads, you
 *will have strongly consistent behavior; that is, readers will
 *always see the most recent write. Of these, the most interesting is
 *to do QUORUM reads and writes, which gives you consistency while
 *still allowing availability in the face of node failures up to half
 *of <ReplicationFactor>. Of course if latency is more important than
 *consistency then you can use lower values for either or both.
 *
 *Some ConsistencyLevels (ONE, TWO, THREE) refer to a specific number
 *of replicas rather than a logical concept that adjusts
 *automatically with the replication factor.  Of these, only ONE is
 *commonly used; TWO and (even more rarely) THREE are only useful
 *when you care more about guaranteeing a certain level of
 *durability, than consistency.
 *
 *Write consistency levels make the following guarantees before reporting success to the client:
 *  ANY          Ensure that the write has been written once somewhere, including possibly being hinted in a non-target node.
 *  ONE          Ensure that the write has been written to at least 1 node's commit log and memory table
 *  TWO          Ensure that the write has been written to at least 2 node's commit log and memory table
 *  THREE        Ensure that the write has been written to at least 3 node's commit log and memory table
 *  QUORUM       Ensure that the write has been written to <ReplicationFactor> / 2 + 1 nodes
 *  LOCAL_QUORUM Ensure that the write has been written to <ReplicationFactor> / 2 + 1 nodes, within the local datacenter (requires NetworkTopologyStrategy)
 *  EACH_QUORUM  Ensure that the write has been written to <ReplicationFactor> / 2 + 1 nodes in each datacenter (requires NetworkTopologyStrategy)
 *  ALL          Ensure that the write is written to <code>&lt;ReplicationFactor&gt;</code> nodes before responding to the client.
 *
 *Read consistency levels make the following guarantees before returning successful results to the client:
 *  ANY          Not supported. You probably want ONE instead.
 *  ONE          Returns the record obtained from a single replica.
 *  TWO          Returns the record with the most recent timestamp once two replicas have replied.
 *  THREE        Returns the record with the most recent timestamp once three replicas have replied.
 *  QUORUM       Returns the record with the most recent timestamp once a majority of replicas have replied.
 *  LOCAL_QUORUM Returns the record with the most recent timestamp once a majority of replicas within the local datacenter have replied.
 *  EACH_QUORUM  Returns the record with the most recent timestamp once a majority of replicas within each datacenter have replied.
 *  ALL          Returns the record with the most recent timestamp once all replicas have replied (implies no replica may be down)..
 */
type ConsistencyLevel int64

const (
	ONE          ConsistencyLevel = 1
	QUORUM       ConsistencyLevel = 2
	LOCAL_QUORUM ConsistencyLevel = 3
	EACH_QUORUM  ConsistencyLevel = 4
	ALL          ConsistencyLevel = 5
	ANY          ConsistencyLevel = 6
	TWO          ConsistencyLevel = 7
	THREE        ConsistencyLevel = 8
)

func (p ConsistencyLevel) String() string {
	switch p {
	case ONE:
		return "ONE"
	case QUORUM:
		return "QUORUM"
	case LOCAL_QUORUM:
		return "LOCAL_QUORUM"
	case EACH_QUORUM:
		return "EACH_QUORUM"
	case ALL:
		return "ALL"
	case ANY:
		return "ANY"
	case TWO:
		return "TWO"
	case THREE:
		return "THREE"
	}
	return "<UNSET>"
}

func FromConsistencyLevelString(s string) ConsistencyLevel {
	switch s {
	case "ONE":
		return ONE
	case "QUORUM":
		return QUORUM
	case "LOCAL_QUORUM":
		return LOCAL_QUORUM
	case "EACH_QUORUM":
		return EACH_QUORUM
	case "ALL":
		return ALL
	case "ANY":
		return ANY
	case "TWO":
		return TWO
	case "THREE":
		return THREE
	}
	return ConsistencyLevel(-10000)
}

func (p ConsistencyLevel) Value() int {
	return int(p)
}

func (p ConsistencyLevel) IsEnum() bool {
	return true
}

type IndexOperator int64

const (
	EQ  IndexOperator = 0
	GTE IndexOperator = 1
	GT  IndexOperator = 2
	LTE IndexOperator = 3
	LT  IndexOperator = 4
)

func (p IndexOperator) String() string {
	switch p {
	case EQ:
		return "EQ"
	case GTE:
		return "GTE"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case LT:
		return "LT"
	}
	return "<UNSET>"
}

func FromIndexOperatorString(s string) IndexOperator {
	switch s {
	case "EQ":
		return EQ
	case "GTE":
		return GTE
	case "GT":
		return GT
	case "LTE":
		return LTE
	case "LT":
		return LT
	}
	return IndexOperator(-10000)
}

func (p IndexOperator) Value() int {
	return int(p)
}

func (p IndexOperator) IsEnum() bool {
	return true
}

type IndexType int64

const (
	KEYS   IndexType = 0
	CUSTOM IndexType = 1
)

func (p IndexType) String() string {
	switch p {
	case KEYS:
		return "KEYS"
	case CUSTOM:
		return "CUSTOM"
	}
	return "<UNSET>"
}

func FromIndexTypeString(s string) IndexType {
	switch s {
	case "KEYS":
		return KEYS
	case "CUSTOM":
		return CUSTOM
	}
	return IndexType(-10000)
}

func (p IndexType) Value() int {
	return int(p)
}

func (p IndexType) IsEnum() bool {
	return true
}

/**
 *CQL query compression
 */
type Compression int64

const (
	GZIP Compression = 1
	NONE Compression = 2
)

func (p Compression) String() string {
	switch p {
	case GZIP:
		return "GZIP"
	case NONE:
		return "NONE"
	}
	return "<UNSET>"
}

func FromCompressionString(s string) Compression {
	switch s {
	case "GZIP":
		return GZIP
	case "NONE":
		return NONE
	}
	return Compression(-10000)
}

func (p Compression) Value() int {
	return int(p)
}

func (p Compression) IsEnum() bool {
	return true
}

type CqlResultType int64

const (
	ROWS CqlResultType = 1
	VOID CqlResultType = 2
	INT  CqlResultType = 3
)

func (p CqlResultType) String() string {
	switch p {
	case ROWS:
		return "ROWS"
	case VOID:
		return "VOID"
	case INT:
		return "INT"
	}
	return "<UNSET>"
}

func FromCqlResultTypeString(s string) CqlResultType {
	switch s {
	case "ROWS":
		return ROWS
	case "VOID":
		return VOID
	case "INT":
		return INT
	}
	return CqlResultType(-10000)
}

func (p CqlResultType) Value() int {
	return int(p)
}

func (p CqlResultType) IsEnum() bool {
	return true
}

/**
 * Basic unit of data within a ColumnFamily.
 * @param name, the name by which this column is set and retrieved.  Maximum 64KB long.
 * @param value. The data associated with the name.  Maximum 2GB long, but in practice you should limit it to small numbers of MB (since Thrift must read the full value into memory to operate on it).
 * @param timestamp. The timestamp is used for conflict detection/resolution when two columns with same name need to be compared.
 * @param ttl. An optional, positive delay (in seconds) after which the column will be automatically deleted.
 * 
 * Attributes:
 *  - Name
 *  - Value
 *  - Timestamp
 *  - Ttl
 */
type Column struct {
	thrift.TStruct
	Name      []byte "name"      // 1
	Value     []byte "value"     // 2
	Timestamp int64  "timestamp" // 3
	Ttl       int32  "ttl"       // 4
}

func NewColumn() *Column {
	output := &Column{
		TStruct: thrift.NewTStruct("Column", []thrift.TField{
			thrift.NewTField("name", thrift.BINARY, 1),
			thrift.NewTField("value", thrift.BINARY, 2),
			thrift.NewTField("timestamp", thrift.I64, 3),
			thrift.NewTField("ttl", thrift.I32, 4),
		}),
	}
	{
	}
	return output
}

func (p *Column) IsSetValue() bool {
	return p.Value != nil
}

func (p *Column) IsSetTimestamp() bool {
	return p.Timestamp != 0
}

func (p *Column) IsSetTtl() bool {
	return p.Ttl != 0
}

func (p *Column) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "value" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "timestamp" {
			if fieldTypeId == thrift.I64 {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "ttl" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *Column) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v0, err1 := iprot.ReadBinary()
	if err1 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "name", p.ThriftName(), err1)
	}
	p.Name = v0
	return err
}

func (p *Column) ReadFieldName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *Column) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v2, err3 := iprot.ReadBinary()
	if err3 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "value", p.ThriftName(), err3)
	}
	p.Value = v2
	return err
}

func (p *Column) ReadFieldValue(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *Column) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v4, err5 := iprot.ReadI64()
	if err5 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "timestamp", p.ThriftName(), err5)
	}
	p.Timestamp = v4
	return err
}

func (p *Column) ReadFieldTimestamp(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *Column) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v6, err7 := iprot.ReadI32()
	if err7 != nil {
		return thrift.NewTProtocolExceptionReadField(4, "ttl", p.ThriftName(), err7)
	}
	p.Ttl = v6
	return err
}

func (p *Column) ReadFieldTtl(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *Column) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("Column")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *Column) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Name != nil {
		err = oprot.WriteFieldBegin("name", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Name)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
	}
	return err
}

func (p *Column) WriteFieldName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *Column) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Value != nil {
		if p.IsSetValue() {
			err = oprot.WriteFieldBegin("value", thrift.BINARY, 2)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
			}
			err = oprot.WriteBinary(p.Value)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *Column) WriteFieldValue(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *Column) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetTimestamp() {
		err = oprot.WriteFieldBegin("timestamp", thrift.I64, 3)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "timestamp", p.ThriftName(), err)
		}
		err = oprot.WriteI64(int64(p.Timestamp))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "timestamp", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "timestamp", p.ThriftName(), err)
		}
	}
	return err
}

func (p *Column) WriteFieldTimestamp(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *Column) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetTtl() {
		err = oprot.WriteFieldBegin("ttl", thrift.I32, 4)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "ttl", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.Ttl))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "ttl", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "ttl", p.ThriftName(), err)
		}
	}
	return err
}

func (p *Column) WriteFieldTtl(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *Column) TStructName() string {
	return "Column"
}

func (p *Column) ThriftName() string {
	return "Column"
}

func (p *Column) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Column(%+v)", *p)
}

func (p *Column) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*Column)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *Column) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Name
	case 2:
		return p.Value
	case 3:
		return p.Timestamp
	case 4:
		return p.Ttl
	}
	return nil
}

func (p *Column) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("name", thrift.BINARY, 1),
		thrift.NewTField("value", thrift.BINARY, 2),
		thrift.NewTField("timestamp", thrift.I64, 3),
		thrift.NewTField("ttl", thrift.I32, 4),
	})
}

/**
 * A named list of columns.
 * @param name. see Column.name.
 * @param columns. A collection of standard Columns.  The columns within a super column are defined in an adhoc manner.
 *                 Columns within a super column do not have to have matching structures (similarly named child columns).
 * 
 * Attributes:
 *  - Name
 *  - Columns
 */
type SuperColumn struct {
	thrift.TStruct
	Name    []byte       "name"    // 1
	Columns thrift.TList "columns" // 2
}

func NewSuperColumn() *SuperColumn {
	output := &SuperColumn{
		TStruct: thrift.NewTStruct("SuperColumn", []thrift.TField{
			thrift.NewTField("name", thrift.BINARY, 1),
			thrift.NewTField("columns", thrift.LIST, 2),
		}),
	}
	{
	}
	return output
}

func (p *SuperColumn) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "columns" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *SuperColumn) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v8, err9 := iprot.ReadBinary()
	if err9 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "name", p.ThriftName(), err9)
	}
	p.Name = v8
	return err
}

func (p *SuperColumn) ReadFieldName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *SuperColumn) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype15, _size12, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Columns", "", err)
	}
	p.Columns = thrift.NewTList(_etype15, _size12)
	for _i16 := 0; _i16 < _size12; _i16++ {
		_elem17 := NewColumn()
		err20 := _elem17.Read(iprot)
		if err20 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem17Column", err20)
		}
		p.Columns.Push(_elem17)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *SuperColumn) ReadFieldColumns(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *SuperColumn) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("SuperColumn")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *SuperColumn) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Name != nil {
		err = oprot.WriteFieldBegin("name", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Name)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
	}
	return err
}

func (p *SuperColumn) WriteFieldName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *SuperColumn) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Columns != nil {
		err = oprot.WriteFieldBegin("columns", thrift.LIST, 2)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "columns", p.ThriftName(), err)
		}
		err = oprot.WriteListBegin(thrift.STRUCT, p.Columns.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		for Iter21 := range p.Columns.Iter() {
			Iter22 := Iter21.(*Column)
			err = Iter22.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("Column", err)
			}
		}
		err = oprot.WriteListEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "columns", p.ThriftName(), err)
		}
	}
	return err
}

func (p *SuperColumn) WriteFieldColumns(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *SuperColumn) TStructName() string {
	return "SuperColumn"
}

func (p *SuperColumn) ThriftName() string {
	return "SuperColumn"
}

func (p *SuperColumn) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SuperColumn(%+v)", *p)
}

func (p *SuperColumn) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*SuperColumn)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *SuperColumn) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Name
	case 2:
		return p.Columns
	}
	return nil
}

func (p *SuperColumn) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("name", thrift.BINARY, 1),
		thrift.NewTField("columns", thrift.LIST, 2),
	})
}

/**
 * Attributes:
 *  - Name
 *  - Value
 */
type CounterColumn struct {
	thrift.TStruct
	Name  []byte "name"  // 1
	Value int64  "value" // 2
}

func NewCounterColumn() *CounterColumn {
	output := &CounterColumn{
		TStruct: thrift.NewTStruct("CounterColumn", []thrift.TField{
			thrift.NewTField("name", thrift.BINARY, 1),
			thrift.NewTField("value", thrift.I64, 2),
		}),
	}
	{
	}
	return output
}

func (p *CounterColumn) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "value" {
			if fieldTypeId == thrift.I64 {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CounterColumn) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v23, err24 := iprot.ReadBinary()
	if err24 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "name", p.ThriftName(), err24)
	}
	p.Name = v23
	return err
}

func (p *CounterColumn) ReadFieldName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *CounterColumn) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v25, err26 := iprot.ReadI64()
	if err26 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "value", p.ThriftName(), err26)
	}
	p.Value = v25
	return err
}

func (p *CounterColumn) ReadFieldValue(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *CounterColumn) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("CounterColumn")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CounterColumn) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Name != nil {
		err = oprot.WriteFieldBegin("name", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Name)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CounterColumn) WriteFieldName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *CounterColumn) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("value", thrift.I64, 2)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
	}
	err = oprot.WriteI64(int64(p.Value))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
	}
	return err
}

func (p *CounterColumn) WriteFieldValue(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *CounterColumn) TStructName() string {
	return "CounterColumn"
}

func (p *CounterColumn) ThriftName() string {
	return "CounterColumn"
}

func (p *CounterColumn) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CounterColumn(%+v)", *p)
}

func (p *CounterColumn) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*CounterColumn)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *CounterColumn) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Name
	case 2:
		return p.Value
	}
	return nil
}

func (p *CounterColumn) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("name", thrift.BINARY, 1),
		thrift.NewTField("value", thrift.I64, 2),
	})
}

/**
 * Attributes:
 *  - Name
 *  - Columns
 */
type CounterSuperColumn struct {
	thrift.TStruct
	Name    []byte       "name"    // 1
	Columns thrift.TList "columns" // 2
}

func NewCounterSuperColumn() *CounterSuperColumn {
	output := &CounterSuperColumn{
		TStruct: thrift.NewTStruct("CounterSuperColumn", []thrift.TField{
			thrift.NewTField("name", thrift.BINARY, 1),
			thrift.NewTField("columns", thrift.LIST, 2),
		}),
	}
	{
	}
	return output
}

func (p *CounterSuperColumn) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "columns" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CounterSuperColumn) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v27, err28 := iprot.ReadBinary()
	if err28 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "name", p.ThriftName(), err28)
	}
	p.Name = v27
	return err
}

func (p *CounterSuperColumn) ReadFieldName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *CounterSuperColumn) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype34, _size31, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Columns", "", err)
	}
	p.Columns = thrift.NewTList(_etype34, _size31)
	for _i35 := 0; _i35 < _size31; _i35++ {
		_elem36 := NewCounterColumn()
		err39 := _elem36.Read(iprot)
		if err39 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem36CounterColumn", err39)
		}
		p.Columns.Push(_elem36)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *CounterSuperColumn) ReadFieldColumns(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *CounterSuperColumn) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("CounterSuperColumn")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CounterSuperColumn) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Name != nil {
		err = oprot.WriteFieldBegin("name", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Name)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CounterSuperColumn) WriteFieldName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *CounterSuperColumn) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Columns != nil {
		err = oprot.WriteFieldBegin("columns", thrift.LIST, 2)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "columns", p.ThriftName(), err)
		}
		err = oprot.WriteListBegin(thrift.STRUCT, p.Columns.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		for Iter40 := range p.Columns.Iter() {
			Iter41 := Iter40.(*CounterColumn)
			err = Iter41.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("CounterColumn", err)
			}
		}
		err = oprot.WriteListEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "columns", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CounterSuperColumn) WriteFieldColumns(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *CounterSuperColumn) TStructName() string {
	return "CounterSuperColumn"
}

func (p *CounterSuperColumn) ThriftName() string {
	return "CounterSuperColumn"
}

func (p *CounterSuperColumn) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CounterSuperColumn(%+v)", *p)
}

func (p *CounterSuperColumn) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*CounterSuperColumn)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *CounterSuperColumn) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Name
	case 2:
		return p.Columns
	}
	return nil
}

func (p *CounterSuperColumn) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("name", thrift.BINARY, 1),
		thrift.NewTField("columns", thrift.LIST, 2),
	})
}

/**
 * Methods for fetching rows/records from Cassandra will return either a single instance of ColumnOrSuperColumn or a list
 * of ColumnOrSuperColumns (get_slice()). If you're looking up a SuperColumn (or list of SuperColumns) then the resulting
 * instances of ColumnOrSuperColumn will have the requested SuperColumn in the attribute super_column. For queries resulting
 * in Columns, those values will be in the attribute column. This change was made between 0.3 and 0.4 to standardize on
 * single query methods that may return either a SuperColumn or Column.
 * 
 * If the query was on a counter column family, you will either get a counter_column (instead of a column) or a
 * counter_super_column (instead of a super_column)
 * 
 * @param column. The Column returned by get() or get_slice().
 * @param super_column. The SuperColumn returned by get() or get_slice().
 * @param counter_column. The Counterolumn returned by get() or get_slice().
 * @param counter_super_column. The CounterSuperColumn returned by get() or get_slice().
 * 
 * Attributes:
 *  - Column
 *  - SuperColumn
 *  - CounterColumn
 *  - CounterSuperColumn
 */
type ColumnOrSuperColumn struct {
	thrift.TStruct
	Column             *Column             "column"               // 1
	SuperColumn        *SuperColumn        "super_column"         // 2
	CounterColumn      *CounterColumn      "counter_column"       // 3
	CounterSuperColumn *CounterSuperColumn "counter_super_column" // 4
}

func NewColumnOrSuperColumn() *ColumnOrSuperColumn {
	output := &ColumnOrSuperColumn{
		TStruct: thrift.NewTStruct("ColumnOrSuperColumn", []thrift.TField{
			thrift.NewTField("column", thrift.STRUCT, 1),
			thrift.NewTField("super_column", thrift.STRUCT, 2),
			thrift.NewTField("counter_column", thrift.STRUCT, 3),
			thrift.NewTField("counter_super_column", thrift.STRUCT, 4),
		}),
	}
	{
	}
	return output
}

func (p *ColumnOrSuperColumn) IsSetColumn() bool {
	return p.Column != nil
}

func (p *ColumnOrSuperColumn) IsSetSuperColumn() bool {
	return p.SuperColumn != nil
}

func (p *ColumnOrSuperColumn) IsSetCounterColumn() bool {
	return p.CounterColumn != nil
}

func (p *ColumnOrSuperColumn) IsSetCounterSuperColumn() bool {
	return p.CounterSuperColumn != nil
}

func (p *ColumnOrSuperColumn) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "column" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "super_column" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "counter_column" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "counter_super_column" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *ColumnOrSuperColumn) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.Column = NewColumn()
	err44 := p.Column.Read(iprot)
	if err44 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.ColumnColumn", err44)
	}
	return err
}

func (p *ColumnOrSuperColumn) ReadFieldColumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *ColumnOrSuperColumn) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.SuperColumn = NewSuperColumn()
	err47 := p.SuperColumn.Read(iprot)
	if err47 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.SuperColumnSuperColumn", err47)
	}
	return err
}

func (p *ColumnOrSuperColumn) ReadFieldSuperColumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *ColumnOrSuperColumn) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.CounterColumn = NewCounterColumn()
	err50 := p.CounterColumn.Read(iprot)
	if err50 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.CounterColumnCounterColumn", err50)
	}
	return err
}

func (p *ColumnOrSuperColumn) ReadFieldCounterColumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *ColumnOrSuperColumn) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.CounterSuperColumn = NewCounterSuperColumn()
	err53 := p.CounterSuperColumn.Read(iprot)
	if err53 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.CounterSuperColumnCounterSuperColumn", err53)
	}
	return err
}

func (p *ColumnOrSuperColumn) ReadFieldCounterSuperColumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *ColumnOrSuperColumn) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("ColumnOrSuperColumn")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *ColumnOrSuperColumn) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Column != nil {
		if p.IsSetColumn() {
			err = oprot.WriteFieldBegin("column", thrift.STRUCT, 1)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "column", p.ThriftName(), err)
			}
			err = p.Column.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("Column", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "column", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *ColumnOrSuperColumn) WriteFieldColumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *ColumnOrSuperColumn) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.SuperColumn != nil {
		if p.IsSetSuperColumn() {
			err = oprot.WriteFieldBegin("super_column", thrift.STRUCT, 2)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "super_column", p.ThriftName(), err)
			}
			err = p.SuperColumn.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("SuperColumn", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "super_column", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *ColumnOrSuperColumn) WriteFieldSuperColumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *ColumnOrSuperColumn) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.CounterColumn != nil {
		if p.IsSetCounterColumn() {
			err = oprot.WriteFieldBegin("counter_column", thrift.STRUCT, 3)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(3, "counter_column", p.ThriftName(), err)
			}
			err = p.CounterColumn.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("CounterColumn", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(3, "counter_column", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *ColumnOrSuperColumn) WriteFieldCounterColumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *ColumnOrSuperColumn) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.CounterSuperColumn != nil {
		if p.IsSetCounterSuperColumn() {
			err = oprot.WriteFieldBegin("counter_super_column", thrift.STRUCT, 4)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "counter_super_column", p.ThriftName(), err)
			}
			err = p.CounterSuperColumn.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("CounterSuperColumn", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "counter_super_column", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *ColumnOrSuperColumn) WriteFieldCounterSuperColumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *ColumnOrSuperColumn) TStructName() string {
	return "ColumnOrSuperColumn"
}

func (p *ColumnOrSuperColumn) ThriftName() string {
	return "ColumnOrSuperColumn"
}

func (p *ColumnOrSuperColumn) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ColumnOrSuperColumn(%+v)", *p)
}

func (p *ColumnOrSuperColumn) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*ColumnOrSuperColumn)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *ColumnOrSuperColumn) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Column
	case 2:
		return p.SuperColumn
	case 3:
		return p.CounterColumn
	case 4:
		return p.CounterSuperColumn
	}
	return nil
}

func (p *ColumnOrSuperColumn) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("column", thrift.STRUCT, 1),
		thrift.NewTField("super_column", thrift.STRUCT, 2),
		thrift.NewTField("counter_column", thrift.STRUCT, 3),
		thrift.NewTField("counter_super_column", thrift.STRUCT, 4),
	})
}

/**
 * A specific column was requested that does not exist.
 */
type NotFoundException struct {
	thrift.TStruct
}

func NewNotFoundException() *NotFoundException {
	output := &NotFoundException{
		TStruct: thrift.NewTStruct("NotFoundException", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *NotFoundException) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *NotFoundException) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("NotFoundException")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *NotFoundException) TStructName() string {
	return "NotFoundException"
}

func (p *NotFoundException) ThriftName() string {
	return "NotFoundException"
}

func (p *NotFoundException) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("NotFoundException(%+v)", *p)
}

func (p *NotFoundException) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*NotFoundException)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *NotFoundException) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *NotFoundException) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Invalid request could mean keyspace or column family does not exist, required parameters are missing, or a parameter is malformed.
 * why contains an associated error message.
 * 
 * Attributes:
 *  - Why
 */
type InvalidRequestException struct {
	thrift.TStruct
	Why string "why" // 1
}

func NewInvalidRequestException() *InvalidRequestException {
	output := &InvalidRequestException{
		TStruct: thrift.NewTStruct("InvalidRequestException", []thrift.TField{
			thrift.NewTField("why", thrift.STRING, 1),
		}),
	}
	{
	}
	return output
}

func (p *InvalidRequestException) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "why" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *InvalidRequestException) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v54, err55 := iprot.ReadString()
	if err55 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "why", p.ThriftName(), err55)
	}
	p.Why = v54
	return err
}

func (p *InvalidRequestException) ReadFieldWhy(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *InvalidRequestException) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("InvalidRequestException")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *InvalidRequestException) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("why", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Why))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	return err
}

func (p *InvalidRequestException) WriteFieldWhy(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *InvalidRequestException) TStructName() string {
	return "InvalidRequestException"
}

func (p *InvalidRequestException) ThriftName() string {
	return "InvalidRequestException"
}

func (p *InvalidRequestException) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("InvalidRequestException(%+v)", *p)
}

func (p *InvalidRequestException) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*InvalidRequestException)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *InvalidRequestException) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Why
	}
	return nil
}

func (p *InvalidRequestException) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("why", thrift.STRING, 1),
	})
}

/**
 * Not all the replicas required could be created and/or read.
 */
type UnavailableException struct {
	thrift.TStruct
}

func NewUnavailableException() *UnavailableException {
	output := &UnavailableException{
		TStruct: thrift.NewTStruct("UnavailableException", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *UnavailableException) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *UnavailableException) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("UnavailableException")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *UnavailableException) TStructName() string {
	return "UnavailableException"
}

func (p *UnavailableException) ThriftName() string {
	return "UnavailableException"
}

func (p *UnavailableException) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UnavailableException(%+v)", *p)
}

func (p *UnavailableException) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*UnavailableException)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *UnavailableException) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *UnavailableException) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * RPC timeout was exceeded.  either a node failed mid-operation, or load was too high, or the requested op was too large.
 */
type TimedOutException struct {
	thrift.TStruct
}

func NewTimedOutException() *TimedOutException {
	output := &TimedOutException{
		TStruct: thrift.NewTStruct("TimedOutException", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *TimedOutException) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *TimedOutException) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("TimedOutException")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *TimedOutException) TStructName() string {
	return "TimedOutException"
}

func (p *TimedOutException) ThriftName() string {
	return "TimedOutException"
}

func (p *TimedOutException) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TimedOutException(%+v)", *p)
}

func (p *TimedOutException) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*TimedOutException)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *TimedOutException) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *TimedOutException) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * invalid authentication request (invalid keyspace, user does not exist, or credentials invalid)
 * 
 * Attributes:
 *  - Why
 */
type AuthenticationException struct {
	thrift.TStruct
	Why string "why" // 1
}

func NewAuthenticationException() *AuthenticationException {
	output := &AuthenticationException{
		TStruct: thrift.NewTStruct("AuthenticationException", []thrift.TField{
			thrift.NewTField("why", thrift.STRING, 1),
		}),
	}
	{
	}
	return output
}

func (p *AuthenticationException) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "why" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *AuthenticationException) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v56, err57 := iprot.ReadString()
	if err57 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "why", p.ThriftName(), err57)
	}
	p.Why = v56
	return err
}

func (p *AuthenticationException) ReadFieldWhy(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *AuthenticationException) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("AuthenticationException")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *AuthenticationException) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("why", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Why))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	return err
}

func (p *AuthenticationException) WriteFieldWhy(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *AuthenticationException) TStructName() string {
	return "AuthenticationException"
}

func (p *AuthenticationException) ThriftName() string {
	return "AuthenticationException"
}

func (p *AuthenticationException) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AuthenticationException(%+v)", *p)
}

func (p *AuthenticationException) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*AuthenticationException)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *AuthenticationException) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Why
	}
	return nil
}

func (p *AuthenticationException) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("why", thrift.STRING, 1),
	})
}

/**
 * invalid authorization request (user does not have access to keyspace)
 * 
 * Attributes:
 *  - Why
 */
type AuthorizationException struct {
	thrift.TStruct
	Why string "why" // 1
}

func NewAuthorizationException() *AuthorizationException {
	output := &AuthorizationException{
		TStruct: thrift.NewTStruct("AuthorizationException", []thrift.TField{
			thrift.NewTField("why", thrift.STRING, 1),
		}),
	}
	{
	}
	return output
}

func (p *AuthorizationException) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "why" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *AuthorizationException) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v58, err59 := iprot.ReadString()
	if err59 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "why", p.ThriftName(), err59)
	}
	p.Why = v58
	return err
}

func (p *AuthorizationException) ReadFieldWhy(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *AuthorizationException) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("AuthorizationException")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *AuthorizationException) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("why", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Why))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "why", p.ThriftName(), err)
	}
	return err
}

func (p *AuthorizationException) WriteFieldWhy(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *AuthorizationException) TStructName() string {
	return "AuthorizationException"
}

func (p *AuthorizationException) ThriftName() string {
	return "AuthorizationException"
}

func (p *AuthorizationException) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AuthorizationException(%+v)", *p)
}

func (p *AuthorizationException) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*AuthorizationException)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *AuthorizationException) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Why
	}
	return nil
}

func (p *AuthorizationException) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("why", thrift.STRING, 1),
	})
}

/**
 * schemas are not in agreement across all nodes
 */
type SchemaDisagreementException struct {
	thrift.TStruct
}

func NewSchemaDisagreementException() *SchemaDisagreementException {
	output := &SchemaDisagreementException{
		TStruct: thrift.NewTStruct("SchemaDisagreementException", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *SchemaDisagreementException) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *SchemaDisagreementException) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("SchemaDisagreementException")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *SchemaDisagreementException) TStructName() string {
	return "SchemaDisagreementException"
}

func (p *SchemaDisagreementException) ThriftName() string {
	return "SchemaDisagreementException"
}

func (p *SchemaDisagreementException) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SchemaDisagreementException(%+v)", *p)
}

func (p *SchemaDisagreementException) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*SchemaDisagreementException)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *SchemaDisagreementException) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *SchemaDisagreementException) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * ColumnParent is used when selecting groups of columns from the same ColumnFamily. In directory structure terms, imagine
 * ColumnParent as ColumnPath + '/../'.
 * 
 * See also <a href="cassandra.html#Struct_ColumnPath">ColumnPath</a>
 * 
 * Attributes:
 *  - ColumnFamily
 *  - SuperColumn
 */
type ColumnParent struct {
	thrift.TStruct
	_            interface{} // nil # 1
	_            interface{} // nil # 2
	ColumnFamily string      "column_family" // 3
	SuperColumn  []byte      "super_column"  // 4
}

func NewColumnParent() *ColumnParent {
	output := &ColumnParent{
		TStruct: thrift.NewTStruct("ColumnParent", []thrift.TField{
			thrift.NewTField("column_family", thrift.STRING, 3),
			thrift.NewTField("super_column", thrift.BINARY, 4),
		}),
	}
	{
	}
	return output
}

func (p *ColumnParent) IsSetSuperColumn() bool {
	return p.SuperColumn != nil
}

func (p *ColumnParent) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 3 || fieldName == "column_family" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "super_column" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *ColumnParent) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v60, err61 := iprot.ReadString()
	if err61 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "column_family", p.ThriftName(), err61)
	}
	p.ColumnFamily = v60
	return err
}

func (p *ColumnParent) ReadFieldColumnFamily(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *ColumnParent) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v62, err63 := iprot.ReadBinary()
	if err63 != nil {
		return thrift.NewTProtocolExceptionReadField(4, "super_column", p.ThriftName(), err63)
	}
	p.SuperColumn = v62
	return err
}

func (p *ColumnParent) ReadFieldSuperColumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *ColumnParent) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("ColumnParent")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *ColumnParent) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("column_family", thrift.STRING, 3)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "column_family", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.ColumnFamily))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "column_family", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "column_family", p.ThriftName(), err)
	}
	return err
}

func (p *ColumnParent) WriteFieldColumnFamily(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *ColumnParent) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.SuperColumn != nil {
		if p.IsSetSuperColumn() {
			err = oprot.WriteFieldBegin("super_column", thrift.BINARY, 4)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "super_column", p.ThriftName(), err)
			}
			err = oprot.WriteBinary(p.SuperColumn)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "super_column", p.ThriftName(), err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "super_column", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *ColumnParent) WriteFieldSuperColumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *ColumnParent) TStructName() string {
	return "ColumnParent"
}

func (p *ColumnParent) ThriftName() string {
	return "ColumnParent"
}

func (p *ColumnParent) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ColumnParent(%+v)", *p)
}

func (p *ColumnParent) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*ColumnParent)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *ColumnParent) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 3:
		return p.ColumnFamily
	case 4:
		return p.SuperColumn
	}
	return nil
}

func (p *ColumnParent) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("column_family", thrift.STRING, 3),
		thrift.NewTField("super_column", thrift.BINARY, 4),
	})
}

/**
 * The ColumnPath is the path to a single column in Cassandra. It might make sense to think of ColumnPath and
 * ColumnParent in terms of a directory structure.
 * 
 * ColumnPath is used to looking up a single column.
 * 
 * @param column_family. The name of the CF of the column being looked up.
 * @param super_column. The super column name.
 * @param column. The column name.
 * 
 * Attributes:
 *  - ColumnFamily
 *  - SuperColumn
 *  - Column
 */
type ColumnPath struct {
	thrift.TStruct
	_            interface{} // nil # 1
	_            interface{} // nil # 2
	ColumnFamily string      "column_family" // 3
	SuperColumn  []byte      "super_column"  // 4
	Column       []byte      "column"        // 5
}

func NewColumnPath() *ColumnPath {
	output := &ColumnPath{
		TStruct: thrift.NewTStruct("ColumnPath", []thrift.TField{
			thrift.NewTField("column_family", thrift.STRING, 3),
			thrift.NewTField("super_column", thrift.BINARY, 4),
			thrift.NewTField("column", thrift.BINARY, 5),
		}),
	}
	{
	}
	return output
}

func (p *ColumnPath) IsSetSuperColumn() bool {
	return p.SuperColumn != nil
}

func (p *ColumnPath) IsSetColumn() bool {
	return p.Column != nil
}

func (p *ColumnPath) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 3 || fieldName == "column_family" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "super_column" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 5 || fieldName == "column" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *ColumnPath) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v64, err65 := iprot.ReadString()
	if err65 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "column_family", p.ThriftName(), err65)
	}
	p.ColumnFamily = v64
	return err
}

func (p *ColumnPath) ReadFieldColumnFamily(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *ColumnPath) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v66, err67 := iprot.ReadBinary()
	if err67 != nil {
		return thrift.NewTProtocolExceptionReadField(4, "super_column", p.ThriftName(), err67)
	}
	p.SuperColumn = v66
	return err
}

func (p *ColumnPath) ReadFieldSuperColumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *ColumnPath) ReadField5(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v68, err69 := iprot.ReadBinary()
	if err69 != nil {
		return thrift.NewTProtocolExceptionReadField(5, "column", p.ThriftName(), err69)
	}
	p.Column = v68
	return err
}

func (p *ColumnPath) ReadFieldColumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField5(iprot)
}

func (p *ColumnPath) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("ColumnPath")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField5(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *ColumnPath) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("column_family", thrift.STRING, 3)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "column_family", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.ColumnFamily))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "column_family", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "column_family", p.ThriftName(), err)
	}
	return err
}

func (p *ColumnPath) WriteFieldColumnFamily(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *ColumnPath) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.SuperColumn != nil {
		if p.IsSetSuperColumn() {
			err = oprot.WriteFieldBegin("super_column", thrift.BINARY, 4)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "super_column", p.ThriftName(), err)
			}
			err = oprot.WriteBinary(p.SuperColumn)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "super_column", p.ThriftName(), err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "super_column", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *ColumnPath) WriteFieldSuperColumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *ColumnPath) WriteField5(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Column != nil {
		if p.IsSetColumn() {
			err = oprot.WriteFieldBegin("column", thrift.BINARY, 5)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(5, "column", p.ThriftName(), err)
			}
			err = oprot.WriteBinary(p.Column)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(5, "column", p.ThriftName(), err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(5, "column", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *ColumnPath) WriteFieldColumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField5(oprot)
}

func (p *ColumnPath) TStructName() string {
	return "ColumnPath"
}

func (p *ColumnPath) ThriftName() string {
	return "ColumnPath"
}

func (p *ColumnPath) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ColumnPath(%+v)", *p)
}

func (p *ColumnPath) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*ColumnPath)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *ColumnPath) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 3:
		return p.ColumnFamily
	case 4:
		return p.SuperColumn
	case 5:
		return p.Column
	}
	return nil
}

func (p *ColumnPath) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("column_family", thrift.STRING, 3),
		thrift.NewTField("super_column", thrift.BINARY, 4),
		thrift.NewTField("column", thrift.BINARY, 5),
	})
}

/**
 * A slice range is a structure that stores basic range, ordering and limit information for a query that will return
 * multiple columns. It could be thought of as Cassandra's version of LIMIT and ORDER BY
 * 
 * @param start. The column name to start the slice with. This attribute is not required, though there is no default value,
 *               and can be safely set to '', i.e., an empty byte array, to start with the first column name. Otherwise, it
 *               must a valid value under the rules of the Comparator defined for the given ColumnFamily.
 * @param finish. The column name to stop the slice at. This attribute is not required, though there is no default value,
 *                and can be safely set to an empty byte array to not stop until 'count' results are seen. Otherwise, it
 *                must also be a valid value to the ColumnFamily Comparator.
 * @param reversed. Whether the results should be ordered in reversed order. Similar to ORDER BY blah DESC in SQL.
 * @param count. How many columns to return. Similar to LIMIT in SQL. May be arbitrarily large, but Thrift will
 *               materialize the whole result into memory before returning it to the client, so be aware that you may
 *               be better served by iterating through slices by passing the last value of one call in as the 'start'
 *               of the next instead of increasing 'count' arbitrarily large.
 * 
 * Attributes:
 *  - Start
 *  - Finish
 *  - Reversed
 *  - Count
 */
type SliceRange struct {
	thrift.TStruct
	Start    []byte "start"    // 1
	Finish   []byte "finish"   // 2
	Reversed bool   "reversed" // 3
	Count    int32  "count"    // 4
}

func NewSliceRange() *SliceRange {
	output := &SliceRange{
		TStruct: thrift.NewTStruct("SliceRange", []thrift.TField{
			thrift.NewTField("start", thrift.BINARY, 1),
			thrift.NewTField("finish", thrift.BINARY, 2),
			thrift.NewTField("reversed", thrift.BOOL, 3),
			thrift.NewTField("count", thrift.I32, 4),
		}),
	}
	{
		output.Reversed = false
		output.Count = 100
	}
	return output
}

func (p *SliceRange) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "start" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "finish" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "reversed" {
			if fieldTypeId == thrift.BOOL {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "count" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *SliceRange) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v70, err71 := iprot.ReadBinary()
	if err71 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "start", p.ThriftName(), err71)
	}
	p.Start = v70
	return err
}

func (p *SliceRange) ReadFieldStart(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *SliceRange) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v72, err73 := iprot.ReadBinary()
	if err73 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "finish", p.ThriftName(), err73)
	}
	p.Finish = v72
	return err
}

func (p *SliceRange) ReadFieldFinish(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *SliceRange) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v74, err75 := iprot.ReadBool()
	if err75 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "reversed", p.ThriftName(), err75)
	}
	p.Reversed = v74
	return err
}

func (p *SliceRange) ReadFieldReversed(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *SliceRange) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v76, err77 := iprot.ReadI32()
	if err77 != nil {
		return thrift.NewTProtocolExceptionReadField(4, "count", p.ThriftName(), err77)
	}
	p.Count = v76
	return err
}

func (p *SliceRange) ReadFieldCount(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *SliceRange) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("SliceRange")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *SliceRange) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Start != nil {
		err = oprot.WriteFieldBegin("start", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "start", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Start)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "start", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "start", p.ThriftName(), err)
		}
	}
	return err
}

func (p *SliceRange) WriteFieldStart(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *SliceRange) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Finish != nil {
		err = oprot.WriteFieldBegin("finish", thrift.BINARY, 2)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "finish", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Finish)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "finish", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "finish", p.ThriftName(), err)
		}
	}
	return err
}

func (p *SliceRange) WriteFieldFinish(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *SliceRange) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("reversed", thrift.BOOL, 3)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "reversed", p.ThriftName(), err)
	}
	err = oprot.WriteBool(bool(p.Reversed))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "reversed", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "reversed", p.ThriftName(), err)
	}
	return err
}

func (p *SliceRange) WriteFieldReversed(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *SliceRange) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("count", thrift.I32, 4)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(4, "count", p.ThriftName(), err)
	}
	err = oprot.WriteI32(int32(p.Count))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(4, "count", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(4, "count", p.ThriftName(), err)
	}
	return err
}

func (p *SliceRange) WriteFieldCount(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *SliceRange) TStructName() string {
	return "SliceRange"
}

func (p *SliceRange) ThriftName() string {
	return "SliceRange"
}

func (p *SliceRange) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SliceRange(%+v)", *p)
}

func (p *SliceRange) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*SliceRange)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *SliceRange) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Start
	case 2:
		return p.Finish
	case 3:
		return p.Reversed
	case 4:
		return p.Count
	}
	return nil
}

func (p *SliceRange) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("start", thrift.BINARY, 1),
		thrift.NewTField("finish", thrift.BINARY, 2),
		thrift.NewTField("reversed", thrift.BOOL, 3),
		thrift.NewTField("count", thrift.I32, 4),
	})
}

/**
 * A SlicePredicate is similar to a mathematic predicate (see http://en.wikipedia.org/wiki/Predicate_(mathematical_logic)),
 * which is described as "a property that the elements of a set have in common."
 * 
 * SlicePredicate's in Cassandra are described with either a list of column_names or a SliceRange.  If column_names is
 * specified, slice_range is ignored.
 * 
 * @param column_name. A list of column names to retrieve. This can be used similar to Memcached's "multi-get" feature
 *                     to fetch N known column names. For instance, if you know you wish to fetch columns 'Joe', 'Jack',
 *                     and 'Jim' you can pass those column names as a list to fetch all three at once.
 * @param slice_range. A SliceRange describing how to range, order, and/or limit the slice.
 * 
 * Attributes:
 *  - ColumnNames
 *  - SliceRange
 */
type SlicePredicate struct {
	thrift.TStruct
	ColumnNames thrift.TList "column_names" // 1
	SliceRange  *SliceRange  "slice_range"  // 2
}

func NewSlicePredicate() *SlicePredicate {
	output := &SlicePredicate{
		TStruct: thrift.NewTStruct("SlicePredicate", []thrift.TField{
			thrift.NewTField("column_names", thrift.LIST, 1),
			thrift.NewTField("slice_range", thrift.STRUCT, 2),
		}),
	}
	{
	}
	return output
}

func (p *SlicePredicate) IsSetColumnNames() bool {
	return p.ColumnNames != nil && p.ColumnNames.Len() > 0
}

func (p *SlicePredicate) IsSetSliceRange() bool {
	return p.SliceRange != nil
}

func (p *SlicePredicate) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "column_names" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "slice_range" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *SlicePredicate) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype83, _size80, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.ColumnNames", "", err)
	}
	p.ColumnNames = thrift.NewTList(_etype83, _size80)
	for _i84 := 0; _i84 < _size80; _i84++ {
		v86, err87 := iprot.ReadBinary()
		if err87 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_elem85", "", err87)
		}
		_elem85 := v86
		p.ColumnNames.Push(_elem85)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *SlicePredicate) ReadFieldColumnNames(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *SlicePredicate) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.SliceRange = NewSliceRange()
	err90 := p.SliceRange.Read(iprot)
	if err90 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.SliceRangeSliceRange", err90)
	}
	return err
}

func (p *SlicePredicate) ReadFieldSliceRange(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *SlicePredicate) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("SlicePredicate")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *SlicePredicate) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.ColumnNames != nil {
		if p.IsSetColumnNames() {
			err = oprot.WriteFieldBegin("column_names", thrift.LIST, 1)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "column_names", p.ThriftName(), err)
			}
			err = oprot.WriteListBegin(thrift.BINARY, p.ColumnNames.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			for Iter91 := range p.ColumnNames.Iter() {
				Iter92 := Iter91.([]byte)
				err = oprot.WriteBinary(Iter92)
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Iter92", "", err)
				}
			}
			err = oprot.WriteListEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "column_names", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *SlicePredicate) WriteFieldColumnNames(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *SlicePredicate) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.SliceRange != nil {
		if p.IsSetSliceRange() {
			err = oprot.WriteFieldBegin("slice_range", thrift.STRUCT, 2)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "slice_range", p.ThriftName(), err)
			}
			err = p.SliceRange.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("SliceRange", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "slice_range", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *SlicePredicate) WriteFieldSliceRange(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *SlicePredicate) TStructName() string {
	return "SlicePredicate"
}

func (p *SlicePredicate) ThriftName() string {
	return "SlicePredicate"
}

func (p *SlicePredicate) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SlicePredicate(%+v)", *p)
}

func (p *SlicePredicate) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*SlicePredicate)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *SlicePredicate) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.ColumnNames
	case 2:
		return p.SliceRange
	}
	return nil
}

func (p *SlicePredicate) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("column_names", thrift.LIST, 1),
		thrift.NewTField("slice_range", thrift.STRUCT, 2),
	})
}

/**
 * Attributes:
 *  - ColumnName
 *  - Op
 *  - Value
 */
type IndexExpression struct {
	thrift.TStruct
	ColumnName []byte        "column_name" // 1
	Op         IndexOperator "op"          // 2
	Value      []byte        "value"       // 3
}

func NewIndexExpression() *IndexExpression {
	output := &IndexExpression{
		TStruct: thrift.NewTStruct("IndexExpression", []thrift.TField{
			thrift.NewTField("column_name", thrift.BINARY, 1),
			thrift.NewTField("op", thrift.I32, 2),
			thrift.NewTField("value", thrift.BINARY, 3),
		}),
	}
	{
		output.Op = math.MinInt32 - 1
	}
	return output
}

func (p *IndexExpression) IsSetOp() bool {
	return int64(p.Op) != math.MinInt32-1
}

func (p *IndexExpression) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "column_name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "op" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "value" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *IndexExpression) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v93, err94 := iprot.ReadBinary()
	if err94 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "column_name", p.ThriftName(), err94)
	}
	p.ColumnName = v93
	return err
}

func (p *IndexExpression) ReadFieldColumnName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *IndexExpression) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v95, err96 := iprot.ReadI32()
	if err96 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "op", p.ThriftName(), err96)
	}
	p.Op = IndexOperator(v95)
	return err
}

func (p *IndexExpression) ReadFieldOp(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *IndexExpression) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v97, err98 := iprot.ReadBinary()
	if err98 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "value", p.ThriftName(), err98)
	}
	p.Value = v97
	return err
}

func (p *IndexExpression) ReadFieldValue(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *IndexExpression) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("IndexExpression")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *IndexExpression) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.ColumnName != nil {
		err = oprot.WriteFieldBegin("column_name", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "column_name", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.ColumnName)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "column_name", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "column_name", p.ThriftName(), err)
		}
	}
	return err
}

func (p *IndexExpression) WriteFieldColumnName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *IndexExpression) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetOp() {
		err = oprot.WriteFieldBegin("op", thrift.I32, 2)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "op", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.Op))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "op", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "op", p.ThriftName(), err)
		}
	}
	return err
}

func (p *IndexExpression) WriteFieldOp(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *IndexExpression) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Value != nil {
		err = oprot.WriteFieldBegin("value", thrift.BINARY, 3)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "value", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Value)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "value", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "value", p.ThriftName(), err)
		}
	}
	return err
}

func (p *IndexExpression) WriteFieldValue(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *IndexExpression) TStructName() string {
	return "IndexExpression"
}

func (p *IndexExpression) ThriftName() string {
	return "IndexExpression"
}

func (p *IndexExpression) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("IndexExpression(%+v)", *p)
}

func (p *IndexExpression) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*IndexExpression)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *IndexExpression) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.ColumnName
	case 2:
		return p.Op
	case 3:
		return p.Value
	}
	return nil
}

func (p *IndexExpression) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("column_name", thrift.BINARY, 1),
		thrift.NewTField("op", thrift.I32, 2),
		thrift.NewTField("value", thrift.BINARY, 3),
	})
}

/**
 * Attributes:
 *  - Expressions
 *  - StartKey
 *  - Count
 */
type IndexClause struct {
	thrift.TStruct
	Expressions thrift.TList "expressions" // 1
	StartKey    []byte       "start_key"   // 2
	Count       int32        "count"       // 3
}

func NewIndexClause() *IndexClause {
	output := &IndexClause{
		TStruct: thrift.NewTStruct("IndexClause", []thrift.TField{
			thrift.NewTField("expressions", thrift.LIST, 1),
			thrift.NewTField("start_key", thrift.BINARY, 2),
			thrift.NewTField("count", thrift.I32, 3),
		}),
	}
	{
		output.Count = 100
	}
	return output
}

func (p *IndexClause) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "expressions" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "start_key" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "count" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *IndexClause) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype104, _size101, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Expressions", "", err)
	}
	p.Expressions = thrift.NewTList(_etype104, _size101)
	for _i105 := 0; _i105 < _size101; _i105++ {
		_elem106 := NewIndexExpression()
		err109 := _elem106.Read(iprot)
		if err109 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem106IndexExpression", err109)
		}
		p.Expressions.Push(_elem106)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *IndexClause) ReadFieldExpressions(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *IndexClause) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v110, err111 := iprot.ReadBinary()
	if err111 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "start_key", p.ThriftName(), err111)
	}
	p.StartKey = v110
	return err
}

func (p *IndexClause) ReadFieldStartKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *IndexClause) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v112, err113 := iprot.ReadI32()
	if err113 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "count", p.ThriftName(), err113)
	}
	p.Count = v112
	return err
}

func (p *IndexClause) ReadFieldCount(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *IndexClause) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("IndexClause")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *IndexClause) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Expressions != nil {
		err = oprot.WriteFieldBegin("expressions", thrift.LIST, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "expressions", p.ThriftName(), err)
		}
		err = oprot.WriteListBegin(thrift.STRUCT, p.Expressions.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		for Iter114 := range p.Expressions.Iter() {
			Iter115 := Iter114.(*IndexExpression)
			err = Iter115.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("IndexExpression", err)
			}
		}
		err = oprot.WriteListEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "expressions", p.ThriftName(), err)
		}
	}
	return err
}

func (p *IndexClause) WriteFieldExpressions(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *IndexClause) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.StartKey != nil {
		err = oprot.WriteFieldBegin("start_key", thrift.BINARY, 2)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "start_key", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.StartKey)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "start_key", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "start_key", p.ThriftName(), err)
		}
	}
	return err
}

func (p *IndexClause) WriteFieldStartKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *IndexClause) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("count", thrift.I32, 3)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "count", p.ThriftName(), err)
	}
	err = oprot.WriteI32(int32(p.Count))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "count", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "count", p.ThriftName(), err)
	}
	return err
}

func (p *IndexClause) WriteFieldCount(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *IndexClause) TStructName() string {
	return "IndexClause"
}

func (p *IndexClause) ThriftName() string {
	return "IndexClause"
}

func (p *IndexClause) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("IndexClause(%+v)", *p)
}

func (p *IndexClause) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*IndexClause)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *IndexClause) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Expressions
	case 2:
		return p.StartKey
	case 3:
		return p.Count
	}
	return nil
}

func (p *IndexClause) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("expressions", thrift.LIST, 1),
		thrift.NewTField("start_key", thrift.BINARY, 2),
		thrift.NewTField("count", thrift.I32, 3),
	})
}

/**
 * The semantics of start keys and tokens are slightly different.
 * Keys are start-inclusive; tokens are start-exclusive.  Token
 * ranges may also wrap -- that is, the end token may be less
 * than the start one.  Thus, a range from keyX to keyX is a
 * one-element range, but a range from tokenY to tokenY is the
 * full ring.
 * 
 * Attributes:
 *  - StartKey
 *  - EndKey
 *  - StartToken
 *  - EndToken
 *  - Count
 */
type KeyRange struct {
	thrift.TStruct
	StartKey   []byte "start_key"   // 1
	EndKey     []byte "end_key"     // 2
	StartToken string "start_token" // 3
	EndToken   string "end_token"   // 4
	Count      int32  "count"       // 5
}

func NewKeyRange() *KeyRange {
	output := &KeyRange{
		TStruct: thrift.NewTStruct("KeyRange", []thrift.TField{
			thrift.NewTField("start_key", thrift.BINARY, 1),
			thrift.NewTField("end_key", thrift.BINARY, 2),
			thrift.NewTField("start_token", thrift.STRING, 3),
			thrift.NewTField("end_token", thrift.STRING, 4),
			thrift.NewTField("count", thrift.I32, 5),
		}),
	}
	{
		output.Count = 100
	}
	return output
}

func (p *KeyRange) IsSetStartKey() bool {
	return p.StartKey != nil
}

func (p *KeyRange) IsSetEndKey() bool {
	return p.EndKey != nil
}

func (p *KeyRange) IsSetStartToken() bool {
	return p.StartToken != ""
}

func (p *KeyRange) IsSetEndToken() bool {
	return p.EndToken != ""
}

func (p *KeyRange) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "start_key" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "end_key" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "start_token" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "end_token" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 5 || fieldName == "count" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *KeyRange) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v116, err117 := iprot.ReadBinary()
	if err117 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "start_key", p.ThriftName(), err117)
	}
	p.StartKey = v116
	return err
}

func (p *KeyRange) ReadFieldStartKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *KeyRange) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v118, err119 := iprot.ReadBinary()
	if err119 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "end_key", p.ThriftName(), err119)
	}
	p.EndKey = v118
	return err
}

func (p *KeyRange) ReadFieldEndKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *KeyRange) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v120, err121 := iprot.ReadString()
	if err121 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "start_token", p.ThriftName(), err121)
	}
	p.StartToken = v120
	return err
}

func (p *KeyRange) ReadFieldStartToken(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *KeyRange) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v122, err123 := iprot.ReadString()
	if err123 != nil {
		return thrift.NewTProtocolExceptionReadField(4, "end_token", p.ThriftName(), err123)
	}
	p.EndToken = v122
	return err
}

func (p *KeyRange) ReadFieldEndToken(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *KeyRange) ReadField5(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v124, err125 := iprot.ReadI32()
	if err125 != nil {
		return thrift.NewTProtocolExceptionReadField(5, "count", p.ThriftName(), err125)
	}
	p.Count = v124
	return err
}

func (p *KeyRange) ReadFieldCount(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField5(iprot)
}

func (p *KeyRange) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("KeyRange")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField5(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *KeyRange) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.StartKey != nil {
		if p.IsSetStartKey() {
			err = oprot.WriteFieldBegin("start_key", thrift.BINARY, 1)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "start_key", p.ThriftName(), err)
			}
			err = oprot.WriteBinary(p.StartKey)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "start_key", p.ThriftName(), err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "start_key", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *KeyRange) WriteFieldStartKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *KeyRange) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.EndKey != nil {
		if p.IsSetEndKey() {
			err = oprot.WriteFieldBegin("end_key", thrift.BINARY, 2)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "end_key", p.ThriftName(), err)
			}
			err = oprot.WriteBinary(p.EndKey)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "end_key", p.ThriftName(), err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "end_key", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *KeyRange) WriteFieldEndKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *KeyRange) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetStartToken() {
		err = oprot.WriteFieldBegin("start_token", thrift.STRING, 3)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "start_token", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.StartToken))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "start_token", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "start_token", p.ThriftName(), err)
		}
	}
	return err
}

func (p *KeyRange) WriteFieldStartToken(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *KeyRange) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetEndToken() {
		err = oprot.WriteFieldBegin("end_token", thrift.STRING, 4)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "end_token", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.EndToken))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "end_token", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "end_token", p.ThriftName(), err)
		}
	}
	return err
}

func (p *KeyRange) WriteFieldEndToken(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *KeyRange) WriteField5(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("count", thrift.I32, 5)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(5, "count", p.ThriftName(), err)
	}
	err = oprot.WriteI32(int32(p.Count))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(5, "count", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(5, "count", p.ThriftName(), err)
	}
	return err
}

func (p *KeyRange) WriteFieldCount(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField5(oprot)
}

func (p *KeyRange) TStructName() string {
	return "KeyRange"
}

func (p *KeyRange) ThriftName() string {
	return "KeyRange"
}

func (p *KeyRange) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KeyRange(%+v)", *p)
}

func (p *KeyRange) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*KeyRange)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *KeyRange) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.StartKey
	case 2:
		return p.EndKey
	case 3:
		return p.StartToken
	case 4:
		return p.EndToken
	case 5:
		return p.Count
	}
	return nil
}

func (p *KeyRange) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("start_key", thrift.BINARY, 1),
		thrift.NewTField("end_key", thrift.BINARY, 2),
		thrift.NewTField("start_token", thrift.STRING, 3),
		thrift.NewTField("end_token", thrift.STRING, 4),
		thrift.NewTField("count", thrift.I32, 5),
	})
}

/**
 * A KeySlice is key followed by the data it maps to. A collection of KeySlice is returned by the get_range_slice operation.
 * 
 * @param key. a row key
 * @param columns. List of data represented by the key. Typically, the list is pared down to only the columns specified by
 *                 a SlicePredicate.
 * 
 * Attributes:
 *  - Key
 *  - Columns
 */
type KeySlice struct {
	thrift.TStruct
	Key     []byte       "key"     // 1
	Columns thrift.TList "columns" // 2
}

func NewKeySlice() *KeySlice {
	output := &KeySlice{
		TStruct: thrift.NewTStruct("KeySlice", []thrift.TField{
			thrift.NewTField("key", thrift.BINARY, 1),
			thrift.NewTField("columns", thrift.LIST, 2),
		}),
	}
	{
	}
	return output
}

func (p *KeySlice) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "key" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "columns" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *KeySlice) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v126, err127 := iprot.ReadBinary()
	if err127 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "key", p.ThriftName(), err127)
	}
	p.Key = v126
	return err
}

func (p *KeySlice) ReadFieldKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *KeySlice) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype133, _size130, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Columns", "", err)
	}
	p.Columns = thrift.NewTList(_etype133, _size130)
	for _i134 := 0; _i134 < _size130; _i134++ {
		_elem135 := NewColumnOrSuperColumn()
		err138 := _elem135.Read(iprot)
		if err138 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem135ColumnOrSuperColumn", err138)
		}
		p.Columns.Push(_elem135)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *KeySlice) ReadFieldColumns(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *KeySlice) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("KeySlice")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *KeySlice) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Key != nil {
		err = oprot.WriteFieldBegin("key", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Key)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
	}
	return err
}

func (p *KeySlice) WriteFieldKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *KeySlice) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Columns != nil {
		err = oprot.WriteFieldBegin("columns", thrift.LIST, 2)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "columns", p.ThriftName(), err)
		}
		err = oprot.WriteListBegin(thrift.STRUCT, p.Columns.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		for Iter139 := range p.Columns.Iter() {
			Iter140 := Iter139.(*ColumnOrSuperColumn)
			err = Iter140.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("ColumnOrSuperColumn", err)
			}
		}
		err = oprot.WriteListEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "columns", p.ThriftName(), err)
		}
	}
	return err
}

func (p *KeySlice) WriteFieldColumns(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *KeySlice) TStructName() string {
	return "KeySlice"
}

func (p *KeySlice) ThriftName() string {
	return "KeySlice"
}

func (p *KeySlice) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KeySlice(%+v)", *p)
}

func (p *KeySlice) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*KeySlice)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *KeySlice) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Key
	case 2:
		return p.Columns
	}
	return nil
}

func (p *KeySlice) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("key", thrift.BINARY, 1),
		thrift.NewTField("columns", thrift.LIST, 2),
	})
}

/**
 * Attributes:
 *  - Key
 *  - Count
 */
type KeyCount struct {
	thrift.TStruct
	Key   []byte "key"   // 1
	Count int32  "count" // 2
}

func NewKeyCount() *KeyCount {
	output := &KeyCount{
		TStruct: thrift.NewTStruct("KeyCount", []thrift.TField{
			thrift.NewTField("key", thrift.BINARY, 1),
			thrift.NewTField("count", thrift.I32, 2),
		}),
	}
	{
	}
	return output
}

func (p *KeyCount) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "key" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "count" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *KeyCount) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v141, err142 := iprot.ReadBinary()
	if err142 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "key", p.ThriftName(), err142)
	}
	p.Key = v141
	return err
}

func (p *KeyCount) ReadFieldKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *KeyCount) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v143, err144 := iprot.ReadI32()
	if err144 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "count", p.ThriftName(), err144)
	}
	p.Count = v143
	return err
}

func (p *KeyCount) ReadFieldCount(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *KeyCount) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("KeyCount")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *KeyCount) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Key != nil {
		err = oprot.WriteFieldBegin("key", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Key)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
	}
	return err
}

func (p *KeyCount) WriteFieldKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *KeyCount) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("count", thrift.I32, 2)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "count", p.ThriftName(), err)
	}
	err = oprot.WriteI32(int32(p.Count))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "count", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "count", p.ThriftName(), err)
	}
	return err
}

func (p *KeyCount) WriteFieldCount(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *KeyCount) TStructName() string {
	return "KeyCount"
}

func (p *KeyCount) ThriftName() string {
	return "KeyCount"
}

func (p *KeyCount) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KeyCount(%+v)", *p)
}

func (p *KeyCount) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*KeyCount)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *KeyCount) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Key
	case 2:
		return p.Count
	}
	return nil
}

func (p *KeyCount) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("key", thrift.BINARY, 1),
		thrift.NewTField("count", thrift.I32, 2),
	})
}

/**
 * Note that the timestamp is only optional in case of counter deletion.
 * 
 * Attributes:
 *  - Timestamp
 *  - SuperColumn
 *  - Predicate
 */
type Deletion struct {
	thrift.TStruct
	Timestamp   int64           "timestamp"    // 1
	SuperColumn []byte          "super_column" // 2
	Predicate   *SlicePredicate "predicate"    // 3
}

func NewDeletion() *Deletion {
	output := &Deletion{
		TStruct: thrift.NewTStruct("Deletion", []thrift.TField{
			thrift.NewTField("timestamp", thrift.I64, 1),
			thrift.NewTField("super_column", thrift.BINARY, 2),
			thrift.NewTField("predicate", thrift.STRUCT, 3),
		}),
	}
	{
	}
	return output
}

func (p *Deletion) IsSetTimestamp() bool {
	return p.Timestamp != 0
}

func (p *Deletion) IsSetSuperColumn() bool {
	return p.SuperColumn != nil
}

func (p *Deletion) IsSetPredicate() bool {
	return p.Predicate != nil
}

func (p *Deletion) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "timestamp" {
			if fieldTypeId == thrift.I64 {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "super_column" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "predicate" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *Deletion) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v145, err146 := iprot.ReadI64()
	if err146 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "timestamp", p.ThriftName(), err146)
	}
	p.Timestamp = v145
	return err
}

func (p *Deletion) ReadFieldTimestamp(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *Deletion) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v147, err148 := iprot.ReadBinary()
	if err148 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "super_column", p.ThriftName(), err148)
	}
	p.SuperColumn = v147
	return err
}

func (p *Deletion) ReadFieldSuperColumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *Deletion) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.Predicate = NewSlicePredicate()
	err151 := p.Predicate.Read(iprot)
	if err151 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.PredicateSlicePredicate", err151)
	}
	return err
}

func (p *Deletion) ReadFieldPredicate(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *Deletion) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("Deletion")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *Deletion) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetTimestamp() {
		err = oprot.WriteFieldBegin("timestamp", thrift.I64, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "timestamp", p.ThriftName(), err)
		}
		err = oprot.WriteI64(int64(p.Timestamp))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "timestamp", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "timestamp", p.ThriftName(), err)
		}
	}
	return err
}

func (p *Deletion) WriteFieldTimestamp(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *Deletion) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.SuperColumn != nil {
		if p.IsSetSuperColumn() {
			err = oprot.WriteFieldBegin("super_column", thrift.BINARY, 2)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "super_column", p.ThriftName(), err)
			}
			err = oprot.WriteBinary(p.SuperColumn)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "super_column", p.ThriftName(), err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "super_column", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *Deletion) WriteFieldSuperColumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *Deletion) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Predicate != nil {
		if p.IsSetPredicate() {
			err = oprot.WriteFieldBegin("predicate", thrift.STRUCT, 3)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(3, "predicate", p.ThriftName(), err)
			}
			err = p.Predicate.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("SlicePredicate", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(3, "predicate", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *Deletion) WriteFieldPredicate(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *Deletion) TStructName() string {
	return "Deletion"
}

func (p *Deletion) ThriftName() string {
	return "Deletion"
}

func (p *Deletion) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Deletion(%+v)", *p)
}

func (p *Deletion) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*Deletion)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *Deletion) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Timestamp
	case 2:
		return p.SuperColumn
	case 3:
		return p.Predicate
	}
	return nil
}

func (p *Deletion) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("timestamp", thrift.I64, 1),
		thrift.NewTField("super_column", thrift.BINARY, 2),
		thrift.NewTField("predicate", thrift.STRUCT, 3),
	})
}

/**
 * A Mutation is either an insert (represented by filling column_or_supercolumn) or a deletion (represented by filling the deletion attribute).
 * @param column_or_supercolumn. An insert to a column or supercolumn (possibly counter column or supercolumn)
 * @param deletion. A deletion of a column or supercolumn
 * 
 * Attributes:
 *  - ColumnOrSupercolumn
 *  - Deletion
 */
type Mutation struct {
	thrift.TStruct
	ColumnOrSupercolumn *ColumnOrSuperColumn "column_or_supercolumn" // 1
	Deletion            *Deletion            "deletion"              // 2
}

func NewMutation() *Mutation {
	output := &Mutation{
		TStruct: thrift.NewTStruct("Mutation", []thrift.TField{
			thrift.NewTField("column_or_supercolumn", thrift.STRUCT, 1),
			thrift.NewTField("deletion", thrift.STRUCT, 2),
		}),
	}
	{
	}
	return output
}

func (p *Mutation) IsSetColumnOrSupercolumn() bool {
	return p.ColumnOrSupercolumn != nil
}

func (p *Mutation) IsSetDeletion() bool {
	return p.Deletion != nil
}

func (p *Mutation) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "column_or_supercolumn" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "deletion" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *Mutation) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.ColumnOrSupercolumn = NewColumnOrSuperColumn()
	err154 := p.ColumnOrSupercolumn.Read(iprot)
	if err154 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.ColumnOrSupercolumnColumnOrSuperColumn", err154)
	}
	return err
}

func (p *Mutation) ReadFieldColumnOrSupercolumn(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *Mutation) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.Deletion = NewDeletion()
	err157 := p.Deletion.Read(iprot)
	if err157 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.DeletionDeletion", err157)
	}
	return err
}

func (p *Mutation) ReadFieldDeletion(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *Mutation) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("Mutation")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *Mutation) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.ColumnOrSupercolumn != nil {
		if p.IsSetColumnOrSupercolumn() {
			err = oprot.WriteFieldBegin("column_or_supercolumn", thrift.STRUCT, 1)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "column_or_supercolumn", p.ThriftName(), err)
			}
			err = p.ColumnOrSupercolumn.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("ColumnOrSuperColumn", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(1, "column_or_supercolumn", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *Mutation) WriteFieldColumnOrSupercolumn(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *Mutation) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Deletion != nil {
		if p.IsSetDeletion() {
			err = oprot.WriteFieldBegin("deletion", thrift.STRUCT, 2)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "deletion", p.ThriftName(), err)
			}
			err = p.Deletion.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("Deletion", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "deletion", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *Mutation) WriteFieldDeletion(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *Mutation) TStructName() string {
	return "Mutation"
}

func (p *Mutation) ThriftName() string {
	return "Mutation"
}

func (p *Mutation) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Mutation(%+v)", *p)
}

func (p *Mutation) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*Mutation)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *Mutation) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.ColumnOrSupercolumn
	case 2:
		return p.Deletion
	}
	return nil
}

func (p *Mutation) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("column_or_supercolumn", thrift.STRUCT, 1),
		thrift.NewTField("deletion", thrift.STRUCT, 2),
	})
}

/**
 * Attributes:
 *  - Host
 *  - Datacenter
 *  - Rack
 */
type EndpointDetails struct {
	thrift.TStruct
	Host       string "host"       // 1
	Datacenter string "datacenter" // 2
	Rack       string "rack"       // 3
}

func NewEndpointDetails() *EndpointDetails {
	output := &EndpointDetails{
		TStruct: thrift.NewTStruct("EndpointDetails", []thrift.TField{
			thrift.NewTField("host", thrift.STRING, 1),
			thrift.NewTField("datacenter", thrift.STRING, 2),
			thrift.NewTField("rack", thrift.STRING, 3),
		}),
	}
	{
	}
	return output
}

func (p *EndpointDetails) IsSetRack() bool {
	return p.Rack != ""
}

func (p *EndpointDetails) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "host" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "datacenter" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "rack" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *EndpointDetails) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v158, err159 := iprot.ReadString()
	if err159 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "host", p.ThriftName(), err159)
	}
	p.Host = v158
	return err
}

func (p *EndpointDetails) ReadFieldHost(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *EndpointDetails) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v160, err161 := iprot.ReadString()
	if err161 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "datacenter", p.ThriftName(), err161)
	}
	p.Datacenter = v160
	return err
}

func (p *EndpointDetails) ReadFieldDatacenter(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *EndpointDetails) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v162, err163 := iprot.ReadString()
	if err163 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "rack", p.ThriftName(), err163)
	}
	p.Rack = v162
	return err
}

func (p *EndpointDetails) ReadFieldRack(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *EndpointDetails) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("EndpointDetails")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *EndpointDetails) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("host", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "host", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Host))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "host", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "host", p.ThriftName(), err)
	}
	return err
}

func (p *EndpointDetails) WriteFieldHost(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *EndpointDetails) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("datacenter", thrift.STRING, 2)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "datacenter", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Datacenter))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "datacenter", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "datacenter", p.ThriftName(), err)
	}
	return err
}

func (p *EndpointDetails) WriteFieldDatacenter(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *EndpointDetails) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetRack() {
		err = oprot.WriteFieldBegin("rack", thrift.STRING, 3)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "rack", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.Rack))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "rack", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "rack", p.ThriftName(), err)
		}
	}
	return err
}

func (p *EndpointDetails) WriteFieldRack(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *EndpointDetails) TStructName() string {
	return "EndpointDetails"
}

func (p *EndpointDetails) ThriftName() string {
	return "EndpointDetails"
}

func (p *EndpointDetails) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("EndpointDetails(%+v)", *p)
}

func (p *EndpointDetails) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*EndpointDetails)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *EndpointDetails) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Host
	case 2:
		return p.Datacenter
	case 3:
		return p.Rack
	}
	return nil
}

func (p *EndpointDetails) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("host", thrift.STRING, 1),
		thrift.NewTField("datacenter", thrift.STRING, 2),
		thrift.NewTField("rack", thrift.STRING, 3),
	})
}

/**
 * A TokenRange describes part of the Cassandra ring, it is a mapping from a range to
 * endpoints responsible for that range.
 * @param start_token The first token in the range
 * @param end_token The last token in the range
 * @param endpoints The endpoints responsible for the range (listed by their configured listen_address)
 * @param rpc_endpoints The endpoints responsible for the range (listed by their configured rpc_address)
 * 
 * Attributes:
 *  - StartToken
 *  - EndToken
 *  - Endpoints
 *  - RpcEndpoints
 *  - EndpointDetails
 */
type TokenRange struct {
	thrift.TStruct
	StartToken      string       "start_token"      // 1
	EndToken        string       "end_token"        // 2
	Endpoints       thrift.TList "endpoints"        // 3
	RpcEndpoints    thrift.TList "rpc_endpoints"    // 4
	EndpointDetails thrift.TList "endpoint_details" // 5
}

func NewTokenRange() *TokenRange {
	output := &TokenRange{
		TStruct: thrift.NewTStruct("TokenRange", []thrift.TField{
			thrift.NewTField("start_token", thrift.STRING, 1),
			thrift.NewTField("end_token", thrift.STRING, 2),
			thrift.NewTField("endpoints", thrift.LIST, 3),
			thrift.NewTField("rpc_endpoints", thrift.LIST, 4),
			thrift.NewTField("endpoint_details", thrift.LIST, 5),
		}),
	}
	{
	}
	return output
}

func (p *TokenRange) IsSetRpcEndpoints() bool {
	return p.RpcEndpoints != nil && p.RpcEndpoints.Len() > 0
}

func (p *TokenRange) IsSetEndpointDetails() bool {
	return p.EndpointDetails != nil && p.EndpointDetails.Len() > 0
}

func (p *TokenRange) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "start_token" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "end_token" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "endpoints" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "rpc_endpoints" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 5 || fieldName == "endpoint_details" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *TokenRange) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v164, err165 := iprot.ReadString()
	if err165 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "start_token", p.ThriftName(), err165)
	}
	p.StartToken = v164
	return err
}

func (p *TokenRange) ReadFieldStartToken(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *TokenRange) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v166, err167 := iprot.ReadString()
	if err167 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "end_token", p.ThriftName(), err167)
	}
	p.EndToken = v166
	return err
}

func (p *TokenRange) ReadFieldEndToken(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *TokenRange) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype173, _size170, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Endpoints", "", err)
	}
	p.Endpoints = thrift.NewTList(_etype173, _size170)
	for _i174 := 0; _i174 < _size170; _i174++ {
		v176, err177 := iprot.ReadString()
		if err177 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_elem175", "", err177)
		}
		_elem175 := v176
		p.Endpoints.Push(_elem175)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *TokenRange) ReadFieldEndpoints(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *TokenRange) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype183, _size180, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.RpcEndpoints", "", err)
	}
	p.RpcEndpoints = thrift.NewTList(_etype183, _size180)
	for _i184 := 0; _i184 < _size180; _i184++ {
		v186, err187 := iprot.ReadString()
		if err187 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_elem185", "", err187)
		}
		_elem185 := v186
		p.RpcEndpoints.Push(_elem185)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *TokenRange) ReadFieldRpcEndpoints(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *TokenRange) ReadField5(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype193, _size190, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.EndpointDetails", "", err)
	}
	p.EndpointDetails = thrift.NewTList(_etype193, _size190)
	for _i194 := 0; _i194 < _size190; _i194++ {
		_elem195 := NewEndpointDetails()
		err198 := _elem195.Read(iprot)
		if err198 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem195EndpointDetails", err198)
		}
		p.EndpointDetails.Push(_elem195)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *TokenRange) ReadFieldEndpointDetails(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField5(iprot)
}

func (p *TokenRange) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("TokenRange")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField5(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *TokenRange) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("start_token", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "start_token", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.StartToken))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "start_token", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "start_token", p.ThriftName(), err)
	}
	return err
}

func (p *TokenRange) WriteFieldStartToken(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *TokenRange) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("end_token", thrift.STRING, 2)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "end_token", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.EndToken))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "end_token", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "end_token", p.ThriftName(), err)
	}
	return err
}

func (p *TokenRange) WriteFieldEndToken(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *TokenRange) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Endpoints != nil {
		err = oprot.WriteFieldBegin("endpoints", thrift.LIST, 3)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "endpoints", p.ThriftName(), err)
		}
		err = oprot.WriteListBegin(thrift.STRING, p.Endpoints.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		for Iter199 := range p.Endpoints.Iter() {
			Iter200 := Iter199.(string)
			err = oprot.WriteString(string(Iter200))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Iter200", "", err)
			}
		}
		err = oprot.WriteListEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "endpoints", p.ThriftName(), err)
		}
	}
	return err
}

func (p *TokenRange) WriteFieldEndpoints(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *TokenRange) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.RpcEndpoints != nil {
		if p.IsSetRpcEndpoints() {
			err = oprot.WriteFieldBegin("rpc_endpoints", thrift.LIST, 4)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "rpc_endpoints", p.ThriftName(), err)
			}
			err = oprot.WriteListBegin(thrift.STRING, p.RpcEndpoints.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			for Iter201 := range p.RpcEndpoints.Iter() {
				Iter202 := Iter201.(string)
				err = oprot.WriteString(string(Iter202))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Iter202", "", err)
				}
			}
			err = oprot.WriteListEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "rpc_endpoints", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *TokenRange) WriteFieldRpcEndpoints(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *TokenRange) WriteField5(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.EndpointDetails != nil {
		if p.IsSetEndpointDetails() {
			err = oprot.WriteFieldBegin("endpoint_details", thrift.LIST, 5)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(5, "endpoint_details", p.ThriftName(), err)
			}
			err = oprot.WriteListBegin(thrift.STRUCT, p.EndpointDetails.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			for Iter203 := range p.EndpointDetails.Iter() {
				Iter204 := Iter203.(*EndpointDetails)
				err = Iter204.Write(oprot)
				if err != nil {
					return thrift.NewTProtocolExceptionWriteStruct("EndpointDetails", err)
				}
			}
			err = oprot.WriteListEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(5, "endpoint_details", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *TokenRange) WriteFieldEndpointDetails(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField5(oprot)
}

func (p *TokenRange) TStructName() string {
	return "TokenRange"
}

func (p *TokenRange) ThriftName() string {
	return "TokenRange"
}

func (p *TokenRange) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TokenRange(%+v)", *p)
}

func (p *TokenRange) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*TokenRange)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *TokenRange) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.StartToken
	case 2:
		return p.EndToken
	case 3:
		return p.Endpoints
	case 4:
		return p.RpcEndpoints
	case 5:
		return p.EndpointDetails
	}
	return nil
}

func (p *TokenRange) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("start_token", thrift.STRING, 1),
		thrift.NewTField("end_token", thrift.STRING, 2),
		thrift.NewTField("endpoints", thrift.LIST, 3),
		thrift.NewTField("rpc_endpoints", thrift.LIST, 4),
		thrift.NewTField("endpoint_details", thrift.LIST, 5),
	})
}

/**
 * Authentication requests can contain any data, dependent on the IAuthenticator used
 * 
 * Attributes:
 *  - Credentials
 */
type AuthenticationRequest struct {
	thrift.TStruct
	Credentials thrift.TMap "credentials" // 1
}

func NewAuthenticationRequest() *AuthenticationRequest {
	output := &AuthenticationRequest{
		TStruct: thrift.NewTStruct("AuthenticationRequest", []thrift.TField{
			thrift.NewTField("credentials", thrift.MAP, 1),
		}),
	}
	{
	}
	return output
}

func (p *AuthenticationRequest) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "credentials" {
			if fieldTypeId == thrift.MAP {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *AuthenticationRequest) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype208, _vtype209, _size207, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Credentials", "", err)
	}
	p.Credentials = thrift.NewTMap(_ktype208, _vtype209, _size207)
	for _i211 := 0; _i211 < _size207; _i211++ {
		v214, err215 := iprot.ReadString()
		if err215 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key212", "", err215)
		}
		_key212 := v214
		v216, err217 := iprot.ReadString()
		if err217 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val213", "", err217)
		}
		_val213 := v216
		p.Credentials.Set(_key212, _val213)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *AuthenticationRequest) ReadFieldCredentials(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *AuthenticationRequest) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("AuthenticationRequest")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *AuthenticationRequest) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Credentials != nil {
		err = oprot.WriteFieldBegin("credentials", thrift.MAP, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "credentials", p.ThriftName(), err)
		}
		err = oprot.WriteMapBegin(thrift.STRING, thrift.STRING, p.Credentials.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		for Miter218 := range p.Credentials.Iter() {
			Kiter219, Viter220 := Miter218.Key().(string), Miter218.Value().(string)
			err = oprot.WriteString(string(Kiter219))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Kiter219", "", err)
			}
			err = oprot.WriteString(string(Viter220))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Viter220", "", err)
			}
		}
		err = oprot.WriteMapEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "credentials", p.ThriftName(), err)
		}
	}
	return err
}

func (p *AuthenticationRequest) WriteFieldCredentials(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *AuthenticationRequest) TStructName() string {
	return "AuthenticationRequest"
}

func (p *AuthenticationRequest) ThriftName() string {
	return "AuthenticationRequest"
}

func (p *AuthenticationRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AuthenticationRequest(%+v)", *p)
}

func (p *AuthenticationRequest) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*AuthenticationRequest)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *AuthenticationRequest) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Credentials
	}
	return nil
}

func (p *AuthenticationRequest) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("credentials", thrift.MAP, 1),
	})
}

/**
 * Attributes:
 *  - Name
 *  - ValidationClass
 *  - IndexType
 *  - IndexName
 *  - IndexOptions
 */
type ColumnDef struct {
	thrift.TStruct
	Name            []byte      "name"             // 1
	ValidationClass string      "validation_class" // 2
	IndexType       IndexType   "index_type"       // 3
	IndexName       string      "index_name"       // 4
	IndexOptions    thrift.TMap "index_options"    // 5
}

func NewColumnDef() *ColumnDef {
	output := &ColumnDef{
		TStruct: thrift.NewTStruct("ColumnDef", []thrift.TField{
			thrift.NewTField("name", thrift.BINARY, 1),
			thrift.NewTField("validation_class", thrift.STRING, 2),
			thrift.NewTField("index_type", thrift.I32, 3),
			thrift.NewTField("index_name", thrift.STRING, 4),
			thrift.NewTField("index_options", thrift.MAP, 5),
		}),
	}
	{
		output.IndexType = math.MinInt32 - 1
	}
	return output
}

func (p *ColumnDef) IsSetIndexType() bool {
	return int64(p.IndexType) != math.MinInt32-1
}

func (p *ColumnDef) IsSetIndexName() bool {
	return p.IndexName != ""
}

func (p *ColumnDef) IsSetIndexOptions() bool {
	return p.IndexOptions != nil && p.IndexOptions.Len() > 0
}

func (p *ColumnDef) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "validation_class" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "index_type" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "index_name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 5 || fieldName == "index_options" {
			if fieldTypeId == thrift.MAP {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *ColumnDef) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v221, err222 := iprot.ReadBinary()
	if err222 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "name", p.ThriftName(), err222)
	}
	p.Name = v221
	return err
}

func (p *ColumnDef) ReadFieldName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *ColumnDef) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v223, err224 := iprot.ReadString()
	if err224 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "validation_class", p.ThriftName(), err224)
	}
	p.ValidationClass = v223
	return err
}

func (p *ColumnDef) ReadFieldValidationClass(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *ColumnDef) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v225, err226 := iprot.ReadI32()
	if err226 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "index_type", p.ThriftName(), err226)
	}
	p.IndexType = IndexType(v225)
	return err
}

func (p *ColumnDef) ReadFieldIndexType(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *ColumnDef) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v227, err228 := iprot.ReadString()
	if err228 != nil {
		return thrift.NewTProtocolExceptionReadField(4, "index_name", p.ThriftName(), err228)
	}
	p.IndexName = v227
	return err
}

func (p *ColumnDef) ReadFieldIndexName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *ColumnDef) ReadField5(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype232, _vtype233, _size231, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.IndexOptions", "", err)
	}
	p.IndexOptions = thrift.NewTMap(_ktype232, _vtype233, _size231)
	for _i235 := 0; _i235 < _size231; _i235++ {
		v238, err239 := iprot.ReadString()
		if err239 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key236", "", err239)
		}
		_key236 := v238
		v240, err241 := iprot.ReadString()
		if err241 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val237", "", err241)
		}
		_val237 := v240
		p.IndexOptions.Set(_key236, _val237)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *ColumnDef) ReadFieldIndexOptions(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField5(iprot)
}

func (p *ColumnDef) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("ColumnDef")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField5(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *ColumnDef) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Name != nil {
		err = oprot.WriteFieldBegin("name", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Name)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
		}
	}
	return err
}

func (p *ColumnDef) WriteFieldName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *ColumnDef) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("validation_class", thrift.STRING, 2)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "validation_class", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.ValidationClass))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "validation_class", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "validation_class", p.ThriftName(), err)
	}
	return err
}

func (p *ColumnDef) WriteFieldValidationClass(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *ColumnDef) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetIndexType() {
		err = oprot.WriteFieldBegin("index_type", thrift.I32, 3)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "index_type", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.IndexType))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "index_type", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "index_type", p.ThriftName(), err)
		}
	}
	return err
}

func (p *ColumnDef) WriteFieldIndexType(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *ColumnDef) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetIndexName() {
		err = oprot.WriteFieldBegin("index_name", thrift.STRING, 4)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "index_name", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.IndexName))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "index_name", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "index_name", p.ThriftName(), err)
		}
	}
	return err
}

func (p *ColumnDef) WriteFieldIndexName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *ColumnDef) WriteField5(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IndexOptions != nil {
		if p.IsSetIndexOptions() {
			err = oprot.WriteFieldBegin("index_options", thrift.MAP, 5)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(5, "index_options", p.ThriftName(), err)
			}
			err = oprot.WriteMapBegin(thrift.STRING, thrift.STRING, p.IndexOptions.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
			}
			for Miter242 := range p.IndexOptions.Iter() {
				Kiter243, Viter244 := Miter242.Key().(string), Miter242.Value().(string)
				err = oprot.WriteString(string(Kiter243))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Kiter243", "", err)
				}
				err = oprot.WriteString(string(Viter244))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Viter244", "", err)
				}
			}
			err = oprot.WriteMapEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(5, "index_options", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *ColumnDef) WriteFieldIndexOptions(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField5(oprot)
}

func (p *ColumnDef) TStructName() string {
	return "ColumnDef"
}

func (p *ColumnDef) ThriftName() string {
	return "ColumnDef"
}

func (p *ColumnDef) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ColumnDef(%+v)", *p)
}

func (p *ColumnDef) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*ColumnDef)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *ColumnDef) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Name
	case 2:
		return p.ValidationClass
	case 3:
		return p.IndexType
	case 4:
		return p.IndexName
	case 5:
		return p.IndexOptions
	}
	return nil
}

func (p *ColumnDef) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("name", thrift.BINARY, 1),
		thrift.NewTField("validation_class", thrift.STRING, 2),
		thrift.NewTField("index_type", thrift.I32, 3),
		thrift.NewTField("index_name", thrift.STRING, 4),
		thrift.NewTField("index_options", thrift.MAP, 5),
	})
}

/**
 * Attributes:
 *  - Keyspace
 *  - Name
 *  - ColumnType
 *  - ComparatorType
 *  - SubcomparatorType
 *  - Comment
 *  - RowCacheSize
 *  - KeyCacheSize
 *  - ReadRepairChance
 *  - ColumnMetadata
 *  - GcGraceSeconds
 *  - DefaultValidationClass
 *  - Id
 *  - MinCompactionThreshold
 *  - MaxCompactionThreshold
 *  - RowCacheSavePeriodInSeconds
 *  - KeyCacheSavePeriodInSeconds
 *  - ReplicateOnWrite
 *  - MergeShardsChance
 *  - KeyValidationClass
 *  - RowCacheProvider
 *  - KeyAlias
 *  - CompactionStrategy
 *  - CompactionStrategyOptions
 *  - RowCacheKeysToSave
 *  - CompressionOptions
 *  - BloomFilterFpChance
 */
type CfDef struct {
	thrift.TStruct
	Keyspace                    string       "keyspace"    // 1
	Name                        string       "name"        // 2
	ColumnType                  string       "column_type" // 3
	_                           interface{}  // nil # 4
	ComparatorType              string       "comparator_type"    // 5
	SubcomparatorType           string       "subcomparator_type" // 6
	_                           interface{}  // nil # 7
	Comment                     string       "comment"        // 8
	RowCacheSize                float64      "row_cache_size" // 9
	_                           interface{}  // nil # 10
	KeyCacheSize                float64      "key_cache_size"                   // 11
	ReadRepairChance            float64      "read_repair_chance"               // 12
	ColumnMetadata              thrift.TList "column_metadata"                  // 13
	GcGraceSeconds              int32        "gc_grace_seconds"                 // 14
	DefaultValidationClass      string       "default_validation_class"         // 15
	Id                          int32        "id"                               // 16
	MinCompactionThreshold      int32        "min_compaction_threshold"         // 17
	MaxCompactionThreshold      int32        "max_compaction_threshold"         // 18
	RowCacheSavePeriodInSeconds int32        "row_cache_save_period_in_seconds" // 19
	KeyCacheSavePeriodInSeconds int32        "key_cache_save_period_in_seconds" // 20
	_                           interface{}  // nil # 21
	_                           interface{}  // nil # 22
	_                           interface{}  // nil # 23
	ReplicateOnWrite            bool         "replicate_on_write"          // 24
	MergeShardsChance           float64      "merge_shards_chance"         // 25
	KeyValidationClass          string       "key_validation_class"        // 26
	RowCacheProvider            string       "row_cache_provider"          // 27
	KeyAlias                    []byte       "key_alias"                   // 28
	CompactionStrategy          string       "compaction_strategy"         // 29
	CompactionStrategyOptions   thrift.TMap  "compaction_strategy_options" // 30
	RowCacheKeysToSave          int32        "row_cache_keys_to_save"      // 31
	CompressionOptions          thrift.TMap  "compression_options"         // 32
	BloomFilterFpChance         float64      "bloom_filter_fp_chance"      // 33
}

func NewCfDef() *CfDef {
	output := &CfDef{
		TStruct: thrift.NewTStruct("CfDef", []thrift.TField{
			thrift.NewTField("keyspace", thrift.STRING, 1),
			thrift.NewTField("name", thrift.STRING, 2),
			thrift.NewTField("column_type", thrift.STRING, 3),
			thrift.NewTField("comparator_type", thrift.STRING, 5),
			thrift.NewTField("subcomparator_type", thrift.STRING, 6),
			thrift.NewTField("comment", thrift.STRING, 8),
			thrift.NewTField("row_cache_size", thrift.DOUBLE, 9),
			thrift.NewTField("key_cache_size", thrift.DOUBLE, 11),
			thrift.NewTField("read_repair_chance", thrift.DOUBLE, 12),
			thrift.NewTField("column_metadata", thrift.LIST, 13),
			thrift.NewTField("gc_grace_seconds", thrift.I32, 14),
			thrift.NewTField("default_validation_class", thrift.STRING, 15),
			thrift.NewTField("id", thrift.I32, 16),
			thrift.NewTField("min_compaction_threshold", thrift.I32, 17),
			thrift.NewTField("max_compaction_threshold", thrift.I32, 18),
			thrift.NewTField("row_cache_save_period_in_seconds", thrift.I32, 19),
			thrift.NewTField("key_cache_save_period_in_seconds", thrift.I32, 20),
			thrift.NewTField("replicate_on_write", thrift.BOOL, 24),
			thrift.NewTField("merge_shards_chance", thrift.DOUBLE, 25),
			thrift.NewTField("key_validation_class", thrift.STRING, 26),
			thrift.NewTField("row_cache_provider", thrift.STRING, 27),
			thrift.NewTField("key_alias", thrift.BINARY, 28),
			thrift.NewTField("compaction_strategy", thrift.STRING, 29),
			thrift.NewTField("compaction_strategy_options", thrift.MAP, 30),
			thrift.NewTField("row_cache_keys_to_save", thrift.I32, 31),
			thrift.NewTField("compression_options", thrift.MAP, 32),
			thrift.NewTField("bloom_filter_fp_chance", thrift.DOUBLE, 33),
		}),
	}
	{
		output.ColumnType = "Standard"
		output.ComparatorType = "BytesType"
		output.RowCacheSize = 0
		output.KeyCacheSize = 200000
		output.ReadRepairChance = 1
	}
	return output
}

func (p *CfDef) IsSetColumnType() bool {
	return p.ColumnType != "Standard"
}

func (p *CfDef) IsSetComparatorType() bool {
	return p.ComparatorType != "BytesType"
}

func (p *CfDef) IsSetSubcomparatorType() bool {
	return p.SubcomparatorType != ""
}

func (p *CfDef) IsSetComment() bool {
	return p.Comment != ""
}

func (p *CfDef) IsSetRowCacheSize() bool {
	return p.RowCacheSize != 0
}

func (p *CfDef) IsSetKeyCacheSize() bool {
	return p.KeyCacheSize != 0
}

func (p *CfDef) IsSetReadRepairChance() bool {
	return p.ReadRepairChance != 1
}

func (p *CfDef) IsSetColumnMetadata() bool {
	return p.ColumnMetadata != nil && p.ColumnMetadata.Len() > 0
}

func (p *CfDef) IsSetGcGraceSeconds() bool {
	return p.GcGraceSeconds != 0
}

func (p *CfDef) IsSetDefaultValidationClass() bool {
	return p.DefaultValidationClass != ""
}

func (p *CfDef) IsSetId() bool {
	return p.Id != 0
}

func (p *CfDef) IsSetMinCompactionThreshold() bool {
	return p.MinCompactionThreshold != 0
}

func (p *CfDef) IsSetMaxCompactionThreshold() bool {
	return p.MaxCompactionThreshold != 0
}

func (p *CfDef) IsSetRowCacheSavePeriodInSeconds() bool {
	return p.RowCacheSavePeriodInSeconds != 0
}

func (p *CfDef) IsSetKeyCacheSavePeriodInSeconds() bool {
	return p.KeyCacheSavePeriodInSeconds != 0
}

func (p *CfDef) IsSetReplicateOnWrite() bool {
	return p.ReplicateOnWrite != false
}

func (p *CfDef) IsSetMergeShardsChance() bool {
	return p.MergeShardsChance != 0
}

func (p *CfDef) IsSetKeyValidationClass() bool {
	return p.KeyValidationClass != ""
}

func (p *CfDef) IsSetRowCacheProvider() bool {
	return p.RowCacheProvider != ""
}

func (p *CfDef) IsSetKeyAlias() bool {
	return p.KeyAlias != nil
}

func (p *CfDef) IsSetCompactionStrategy() bool {
	return p.CompactionStrategy != ""
}

func (p *CfDef) IsSetCompactionStrategyOptions() bool {
	return p.CompactionStrategyOptions != nil && p.CompactionStrategyOptions.Len() > 0
}

func (p *CfDef) IsSetRowCacheKeysToSave() bool {
	return p.RowCacheKeysToSave != 0
}

func (p *CfDef) IsSetCompressionOptions() bool {
	return p.CompressionOptions != nil && p.CompressionOptions.Len() > 0
}

func (p *CfDef) IsSetBloomFilterFpChance() bool {
	return p.BloomFilterFpChance != 0
}

func (p *CfDef) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "keyspace" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "column_type" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 5 || fieldName == "comparator_type" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 6 || fieldName == "subcomparator_type" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField6(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField6(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 8 || fieldName == "comment" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField8(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField8(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 9 || fieldName == "row_cache_size" {
			if fieldTypeId == thrift.DOUBLE {
				err = p.ReadField9(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField9(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 11 || fieldName == "key_cache_size" {
			if fieldTypeId == thrift.DOUBLE {
				err = p.ReadField11(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField11(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 12 || fieldName == "read_repair_chance" {
			if fieldTypeId == thrift.DOUBLE {
				err = p.ReadField12(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField12(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 13 || fieldName == "column_metadata" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField13(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField13(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 14 || fieldName == "gc_grace_seconds" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField14(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField14(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 15 || fieldName == "default_validation_class" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField15(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField15(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 16 || fieldName == "id" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField16(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField16(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 17 || fieldName == "min_compaction_threshold" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField17(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField17(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 18 || fieldName == "max_compaction_threshold" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField18(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField18(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 19 || fieldName == "row_cache_save_period_in_seconds" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField19(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField19(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 20 || fieldName == "key_cache_save_period_in_seconds" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField20(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField20(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 24 || fieldName == "replicate_on_write" {
			if fieldTypeId == thrift.BOOL {
				err = p.ReadField24(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField24(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 25 || fieldName == "merge_shards_chance" {
			if fieldTypeId == thrift.DOUBLE {
				err = p.ReadField25(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField25(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 26 || fieldName == "key_validation_class" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField26(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField26(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 27 || fieldName == "row_cache_provider" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField27(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField27(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 28 || fieldName == "key_alias" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField28(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField28(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 29 || fieldName == "compaction_strategy" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField29(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField29(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 30 || fieldName == "compaction_strategy_options" {
			if fieldTypeId == thrift.MAP {
				err = p.ReadField30(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField30(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 31 || fieldName == "row_cache_keys_to_save" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField31(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField31(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 32 || fieldName == "compression_options" {
			if fieldTypeId == thrift.MAP {
				err = p.ReadField32(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField32(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 33 || fieldName == "bloom_filter_fp_chance" {
			if fieldTypeId == thrift.DOUBLE {
				err = p.ReadField33(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField33(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CfDef) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v245, err246 := iprot.ReadString()
	if err246 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "keyspace", p.ThriftName(), err246)
	}
	p.Keyspace = v245
	return err
}

func (p *CfDef) ReadFieldKeyspace(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *CfDef) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v247, err248 := iprot.ReadString()
	if err248 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "name", p.ThriftName(), err248)
	}
	p.Name = v247
	return err
}

func (p *CfDef) ReadFieldName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *CfDef) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v249, err250 := iprot.ReadString()
	if err250 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "column_type", p.ThriftName(), err250)
	}
	p.ColumnType = v249
	return err
}

func (p *CfDef) ReadFieldColumnType(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *CfDef) ReadField5(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v251, err252 := iprot.ReadString()
	if err252 != nil {
		return thrift.NewTProtocolExceptionReadField(5, "comparator_type", p.ThriftName(), err252)
	}
	p.ComparatorType = v251
	return err
}

func (p *CfDef) ReadFieldComparatorType(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField5(iprot)
}

func (p *CfDef) ReadField6(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v253, err254 := iprot.ReadString()
	if err254 != nil {
		return thrift.NewTProtocolExceptionReadField(6, "subcomparator_type", p.ThriftName(), err254)
	}
	p.SubcomparatorType = v253
	return err
}

func (p *CfDef) ReadFieldSubcomparatorType(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField6(iprot)
}

func (p *CfDef) ReadField8(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v255, err256 := iprot.ReadString()
	if err256 != nil {
		return thrift.NewTProtocolExceptionReadField(8, "comment", p.ThriftName(), err256)
	}
	p.Comment = v255
	return err
}

func (p *CfDef) ReadFieldComment(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField8(iprot)
}

func (p *CfDef) ReadField9(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v257, err258 := iprot.ReadDouble()
	if err258 != nil {
		return thrift.NewTProtocolExceptionReadField(9, "row_cache_size", p.ThriftName(), err258)
	}
	p.RowCacheSize = v257
	return err
}

func (p *CfDef) ReadFieldRowCacheSize(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField9(iprot)
}

func (p *CfDef) ReadField11(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v259, err260 := iprot.ReadDouble()
	if err260 != nil {
		return thrift.NewTProtocolExceptionReadField(11, "key_cache_size", p.ThriftName(), err260)
	}
	p.KeyCacheSize = v259
	return err
}

func (p *CfDef) ReadFieldKeyCacheSize(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField11(iprot)
}

func (p *CfDef) ReadField12(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v261, err262 := iprot.ReadDouble()
	if err262 != nil {
		return thrift.NewTProtocolExceptionReadField(12, "read_repair_chance", p.ThriftName(), err262)
	}
	p.ReadRepairChance = v261
	return err
}

func (p *CfDef) ReadFieldReadRepairChance(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField12(iprot)
}

func (p *CfDef) ReadField13(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype268, _size265, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.ColumnMetadata", "", err)
	}
	p.ColumnMetadata = thrift.NewTList(_etype268, _size265)
	for _i269 := 0; _i269 < _size265; _i269++ {
		_elem270 := NewColumnDef()
		err273 := _elem270.Read(iprot)
		if err273 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem270ColumnDef", err273)
		}
		p.ColumnMetadata.Push(_elem270)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *CfDef) ReadFieldColumnMetadata(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField13(iprot)
}

func (p *CfDef) ReadField14(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v274, err275 := iprot.ReadI32()
	if err275 != nil {
		return thrift.NewTProtocolExceptionReadField(14, "gc_grace_seconds", p.ThriftName(), err275)
	}
	p.GcGraceSeconds = v274
	return err
}

func (p *CfDef) ReadFieldGcGraceSeconds(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField14(iprot)
}

func (p *CfDef) ReadField15(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v276, err277 := iprot.ReadString()
	if err277 != nil {
		return thrift.NewTProtocolExceptionReadField(15, "default_validation_class", p.ThriftName(), err277)
	}
	p.DefaultValidationClass = v276
	return err
}

func (p *CfDef) ReadFieldDefaultValidationClass(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField15(iprot)
}

func (p *CfDef) ReadField16(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v278, err279 := iprot.ReadI32()
	if err279 != nil {
		return thrift.NewTProtocolExceptionReadField(16, "id", p.ThriftName(), err279)
	}
	p.Id = v278
	return err
}

func (p *CfDef) ReadFieldId(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField16(iprot)
}

func (p *CfDef) ReadField17(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v280, err281 := iprot.ReadI32()
	if err281 != nil {
		return thrift.NewTProtocolExceptionReadField(17, "min_compaction_threshold", p.ThriftName(), err281)
	}
	p.MinCompactionThreshold = v280
	return err
}

func (p *CfDef) ReadFieldMinCompactionThreshold(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField17(iprot)
}

func (p *CfDef) ReadField18(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v282, err283 := iprot.ReadI32()
	if err283 != nil {
		return thrift.NewTProtocolExceptionReadField(18, "max_compaction_threshold", p.ThriftName(), err283)
	}
	p.MaxCompactionThreshold = v282
	return err
}

func (p *CfDef) ReadFieldMaxCompactionThreshold(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField18(iprot)
}

func (p *CfDef) ReadField19(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v284, err285 := iprot.ReadI32()
	if err285 != nil {
		return thrift.NewTProtocolExceptionReadField(19, "row_cache_save_period_in_seconds", p.ThriftName(), err285)
	}
	p.RowCacheSavePeriodInSeconds = v284
	return err
}

func (p *CfDef) ReadFieldRowCacheSavePeriodInSeconds(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField19(iprot)
}

func (p *CfDef) ReadField20(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v286, err287 := iprot.ReadI32()
	if err287 != nil {
		return thrift.NewTProtocolExceptionReadField(20, "key_cache_save_period_in_seconds", p.ThriftName(), err287)
	}
	p.KeyCacheSavePeriodInSeconds = v286
	return err
}

func (p *CfDef) ReadFieldKeyCacheSavePeriodInSeconds(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField20(iprot)
}

func (p *CfDef) ReadField24(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v288, err289 := iprot.ReadBool()
	if err289 != nil {
		return thrift.NewTProtocolExceptionReadField(24, "replicate_on_write", p.ThriftName(), err289)
	}
	p.ReplicateOnWrite = v288
	return err
}

func (p *CfDef) ReadFieldReplicateOnWrite(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField24(iprot)
}

func (p *CfDef) ReadField25(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v290, err291 := iprot.ReadDouble()
	if err291 != nil {
		return thrift.NewTProtocolExceptionReadField(25, "merge_shards_chance", p.ThriftName(), err291)
	}
	p.MergeShardsChance = v290
	return err
}

func (p *CfDef) ReadFieldMergeShardsChance(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField25(iprot)
}

func (p *CfDef) ReadField26(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v292, err293 := iprot.ReadString()
	if err293 != nil {
		return thrift.NewTProtocolExceptionReadField(26, "key_validation_class", p.ThriftName(), err293)
	}
	p.KeyValidationClass = v292
	return err
}

func (p *CfDef) ReadFieldKeyValidationClass(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField26(iprot)
}

func (p *CfDef) ReadField27(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v294, err295 := iprot.ReadString()
	if err295 != nil {
		return thrift.NewTProtocolExceptionReadField(27, "row_cache_provider", p.ThriftName(), err295)
	}
	p.RowCacheProvider = v294
	return err
}

func (p *CfDef) ReadFieldRowCacheProvider(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField27(iprot)
}

func (p *CfDef) ReadField28(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v296, err297 := iprot.ReadBinary()
	if err297 != nil {
		return thrift.NewTProtocolExceptionReadField(28, "key_alias", p.ThriftName(), err297)
	}
	p.KeyAlias = v296
	return err
}

func (p *CfDef) ReadFieldKeyAlias(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField28(iprot)
}

func (p *CfDef) ReadField29(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v298, err299 := iprot.ReadString()
	if err299 != nil {
		return thrift.NewTProtocolExceptionReadField(29, "compaction_strategy", p.ThriftName(), err299)
	}
	p.CompactionStrategy = v298
	return err
}

func (p *CfDef) ReadFieldCompactionStrategy(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField29(iprot)
}

func (p *CfDef) ReadField30(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype303, _vtype304, _size302, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.CompactionStrategyOptions", "", err)
	}
	p.CompactionStrategyOptions = thrift.NewTMap(_ktype303, _vtype304, _size302)
	for _i306 := 0; _i306 < _size302; _i306++ {
		v309, err310 := iprot.ReadString()
		if err310 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key307", "", err310)
		}
		_key307 := v309
		v311, err312 := iprot.ReadString()
		if err312 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val308", "", err312)
		}
		_val308 := v311
		p.CompactionStrategyOptions.Set(_key307, _val308)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *CfDef) ReadFieldCompactionStrategyOptions(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField30(iprot)
}

func (p *CfDef) ReadField31(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v313, err314 := iprot.ReadI32()
	if err314 != nil {
		return thrift.NewTProtocolExceptionReadField(31, "row_cache_keys_to_save", p.ThriftName(), err314)
	}
	p.RowCacheKeysToSave = v313
	return err
}

func (p *CfDef) ReadFieldRowCacheKeysToSave(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField31(iprot)
}

func (p *CfDef) ReadField32(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype318, _vtype319, _size317, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.CompressionOptions", "", err)
	}
	p.CompressionOptions = thrift.NewTMap(_ktype318, _vtype319, _size317)
	for _i321 := 0; _i321 < _size317; _i321++ {
		v324, err325 := iprot.ReadString()
		if err325 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key322", "", err325)
		}
		_key322 := v324
		v326, err327 := iprot.ReadString()
		if err327 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val323", "", err327)
		}
		_val323 := v326
		p.CompressionOptions.Set(_key322, _val323)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *CfDef) ReadFieldCompressionOptions(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField32(iprot)
}

func (p *CfDef) ReadField33(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v328, err329 := iprot.ReadDouble()
	if err329 != nil {
		return thrift.NewTProtocolExceptionReadField(33, "bloom_filter_fp_chance", p.ThriftName(), err329)
	}
	p.BloomFilterFpChance = v328
	return err
}

func (p *CfDef) ReadFieldBloomFilterFpChance(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField33(iprot)
}

func (p *CfDef) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("CfDef")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField5(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField6(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField8(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField9(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField11(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField12(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField13(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField14(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField15(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField16(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField17(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField18(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField19(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField20(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField24(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField25(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField26(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField27(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField28(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField29(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField30(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField31(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField32(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField33(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CfDef) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("keyspace", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "keyspace", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Keyspace))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "keyspace", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "keyspace", p.ThriftName(), err)
	}
	return err
}

func (p *CfDef) WriteFieldKeyspace(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *CfDef) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("name", thrift.STRING, 2)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "name", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Name))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "name", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "name", p.ThriftName(), err)
	}
	return err
}

func (p *CfDef) WriteFieldName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *CfDef) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetColumnType() {
		err = oprot.WriteFieldBegin("column_type", thrift.STRING, 3)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "column_type", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.ColumnType))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "column_type", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "column_type", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldColumnType(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *CfDef) WriteField5(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetComparatorType() {
		err = oprot.WriteFieldBegin("comparator_type", thrift.STRING, 5)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(5, "comparator_type", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.ComparatorType))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(5, "comparator_type", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(5, "comparator_type", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldComparatorType(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField5(oprot)
}

func (p *CfDef) WriteField6(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetSubcomparatorType() {
		err = oprot.WriteFieldBegin("subcomparator_type", thrift.STRING, 6)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(6, "subcomparator_type", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.SubcomparatorType))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(6, "subcomparator_type", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(6, "subcomparator_type", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldSubcomparatorType(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField6(oprot)
}

func (p *CfDef) WriteField8(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetComment() {
		err = oprot.WriteFieldBegin("comment", thrift.STRING, 8)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(8, "comment", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.Comment))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(8, "comment", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(8, "comment", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldComment(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField8(oprot)
}

func (p *CfDef) WriteField9(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetRowCacheSize() {
		err = oprot.WriteFieldBegin("row_cache_size", thrift.DOUBLE, 9)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(9, "row_cache_size", p.ThriftName(), err)
		}
		err = oprot.WriteDouble(float64(p.RowCacheSize))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(9, "row_cache_size", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(9, "row_cache_size", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldRowCacheSize(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField9(oprot)
}

func (p *CfDef) WriteField11(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetKeyCacheSize() {
		err = oprot.WriteFieldBegin("key_cache_size", thrift.DOUBLE, 11)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(11, "key_cache_size", p.ThriftName(), err)
		}
		err = oprot.WriteDouble(float64(p.KeyCacheSize))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(11, "key_cache_size", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(11, "key_cache_size", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldKeyCacheSize(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField11(oprot)
}

func (p *CfDef) WriteField12(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetReadRepairChance() {
		err = oprot.WriteFieldBegin("read_repair_chance", thrift.DOUBLE, 12)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(12, "read_repair_chance", p.ThriftName(), err)
		}
		err = oprot.WriteDouble(float64(p.ReadRepairChance))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(12, "read_repair_chance", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(12, "read_repair_chance", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldReadRepairChance(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField12(oprot)
}

func (p *CfDef) WriteField13(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.ColumnMetadata != nil {
		if p.IsSetColumnMetadata() {
			err = oprot.WriteFieldBegin("column_metadata", thrift.LIST, 13)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(13, "column_metadata", p.ThriftName(), err)
			}
			err = oprot.WriteListBegin(thrift.STRUCT, p.ColumnMetadata.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			for Iter330 := range p.ColumnMetadata.Iter() {
				Iter331 := Iter330.(*ColumnDef)
				err = Iter331.Write(oprot)
				if err != nil {
					return thrift.NewTProtocolExceptionWriteStruct("ColumnDef", err)
				}
			}
			err = oprot.WriteListEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(13, "column_metadata", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *CfDef) WriteFieldColumnMetadata(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField13(oprot)
}

func (p *CfDef) WriteField14(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetGcGraceSeconds() {
		err = oprot.WriteFieldBegin("gc_grace_seconds", thrift.I32, 14)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(14, "gc_grace_seconds", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.GcGraceSeconds))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(14, "gc_grace_seconds", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(14, "gc_grace_seconds", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldGcGraceSeconds(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField14(oprot)
}

func (p *CfDef) WriteField15(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetDefaultValidationClass() {
		err = oprot.WriteFieldBegin("default_validation_class", thrift.STRING, 15)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(15, "default_validation_class", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.DefaultValidationClass))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(15, "default_validation_class", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(15, "default_validation_class", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldDefaultValidationClass(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField15(oprot)
}

func (p *CfDef) WriteField16(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetId() {
		err = oprot.WriteFieldBegin("id", thrift.I32, 16)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(16, "id", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.Id))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(16, "id", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(16, "id", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldId(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField16(oprot)
}

func (p *CfDef) WriteField17(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetMinCompactionThreshold() {
		err = oprot.WriteFieldBegin("min_compaction_threshold", thrift.I32, 17)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(17, "min_compaction_threshold", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.MinCompactionThreshold))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(17, "min_compaction_threshold", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(17, "min_compaction_threshold", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldMinCompactionThreshold(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField17(oprot)
}

func (p *CfDef) WriteField18(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetMaxCompactionThreshold() {
		err = oprot.WriteFieldBegin("max_compaction_threshold", thrift.I32, 18)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(18, "max_compaction_threshold", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.MaxCompactionThreshold))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(18, "max_compaction_threshold", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(18, "max_compaction_threshold", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldMaxCompactionThreshold(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField18(oprot)
}

func (p *CfDef) WriteField19(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetRowCacheSavePeriodInSeconds() {
		err = oprot.WriteFieldBegin("row_cache_save_period_in_seconds", thrift.I32, 19)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(19, "row_cache_save_period_in_seconds", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.RowCacheSavePeriodInSeconds))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(19, "row_cache_save_period_in_seconds", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(19, "row_cache_save_period_in_seconds", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldRowCacheSavePeriodInSeconds(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField19(oprot)
}

func (p *CfDef) WriteField20(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetKeyCacheSavePeriodInSeconds() {
		err = oprot.WriteFieldBegin("key_cache_save_period_in_seconds", thrift.I32, 20)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(20, "key_cache_save_period_in_seconds", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.KeyCacheSavePeriodInSeconds))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(20, "key_cache_save_period_in_seconds", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(20, "key_cache_save_period_in_seconds", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldKeyCacheSavePeriodInSeconds(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField20(oprot)
}

func (p *CfDef) WriteField24(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetReplicateOnWrite() {
		err = oprot.WriteFieldBegin("replicate_on_write", thrift.BOOL, 24)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(24, "replicate_on_write", p.ThriftName(), err)
		}
		err = oprot.WriteBool(bool(p.ReplicateOnWrite))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(24, "replicate_on_write", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(24, "replicate_on_write", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldReplicateOnWrite(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField24(oprot)
}

func (p *CfDef) WriteField25(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetMergeShardsChance() {
		err = oprot.WriteFieldBegin("merge_shards_chance", thrift.DOUBLE, 25)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(25, "merge_shards_chance", p.ThriftName(), err)
		}
		err = oprot.WriteDouble(float64(p.MergeShardsChance))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(25, "merge_shards_chance", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(25, "merge_shards_chance", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldMergeShardsChance(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField25(oprot)
}

func (p *CfDef) WriteField26(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetKeyValidationClass() {
		err = oprot.WriteFieldBegin("key_validation_class", thrift.STRING, 26)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(26, "key_validation_class", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.KeyValidationClass))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(26, "key_validation_class", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(26, "key_validation_class", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldKeyValidationClass(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField26(oprot)
}

func (p *CfDef) WriteField27(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetRowCacheProvider() {
		err = oprot.WriteFieldBegin("row_cache_provider", thrift.STRING, 27)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(27, "row_cache_provider", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.RowCacheProvider))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(27, "row_cache_provider", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(27, "row_cache_provider", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldRowCacheProvider(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField27(oprot)
}

func (p *CfDef) WriteField28(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.KeyAlias != nil {
		if p.IsSetKeyAlias() {
			err = oprot.WriteFieldBegin("key_alias", thrift.BINARY, 28)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(28, "key_alias", p.ThriftName(), err)
			}
			err = oprot.WriteBinary(p.KeyAlias)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(28, "key_alias", p.ThriftName(), err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(28, "key_alias", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *CfDef) WriteFieldKeyAlias(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField28(oprot)
}

func (p *CfDef) WriteField29(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetCompactionStrategy() {
		err = oprot.WriteFieldBegin("compaction_strategy", thrift.STRING, 29)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(29, "compaction_strategy", p.ThriftName(), err)
		}
		err = oprot.WriteString(string(p.CompactionStrategy))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(29, "compaction_strategy", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(29, "compaction_strategy", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldCompactionStrategy(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField29(oprot)
}

func (p *CfDef) WriteField30(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.CompactionStrategyOptions != nil {
		if p.IsSetCompactionStrategyOptions() {
			err = oprot.WriteFieldBegin("compaction_strategy_options", thrift.MAP, 30)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(30, "compaction_strategy_options", p.ThriftName(), err)
			}
			err = oprot.WriteMapBegin(thrift.STRING, thrift.STRING, p.CompactionStrategyOptions.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
			}
			for Miter332 := range p.CompactionStrategyOptions.Iter() {
				Kiter333, Viter334 := Miter332.Key().(string), Miter332.Value().(string)
				err = oprot.WriteString(string(Kiter333))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Kiter333", "", err)
				}
				err = oprot.WriteString(string(Viter334))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Viter334", "", err)
				}
			}
			err = oprot.WriteMapEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(30, "compaction_strategy_options", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *CfDef) WriteFieldCompactionStrategyOptions(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField30(oprot)
}

func (p *CfDef) WriteField31(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetRowCacheKeysToSave() {
		err = oprot.WriteFieldBegin("row_cache_keys_to_save", thrift.I32, 31)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(31, "row_cache_keys_to_save", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.RowCacheKeysToSave))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(31, "row_cache_keys_to_save", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(31, "row_cache_keys_to_save", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldRowCacheKeysToSave(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField31(oprot)
}

func (p *CfDef) WriteField32(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.CompressionOptions != nil {
		if p.IsSetCompressionOptions() {
			err = oprot.WriteFieldBegin("compression_options", thrift.MAP, 32)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(32, "compression_options", p.ThriftName(), err)
			}
			err = oprot.WriteMapBegin(thrift.STRING, thrift.STRING, p.CompressionOptions.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
			}
			for Miter335 := range p.CompressionOptions.Iter() {
				Kiter336, Viter337 := Miter335.Key().(string), Miter335.Value().(string)
				err = oprot.WriteString(string(Kiter336))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Kiter336", "", err)
				}
				err = oprot.WriteString(string(Viter337))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Viter337", "", err)
				}
			}
			err = oprot.WriteMapEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(32, "compression_options", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *CfDef) WriteFieldCompressionOptions(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField32(oprot)
}

func (p *CfDef) WriteField33(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetBloomFilterFpChance() {
		err = oprot.WriteFieldBegin("bloom_filter_fp_chance", thrift.DOUBLE, 33)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(33, "bloom_filter_fp_chance", p.ThriftName(), err)
		}
		err = oprot.WriteDouble(float64(p.BloomFilterFpChance))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(33, "bloom_filter_fp_chance", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(33, "bloom_filter_fp_chance", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CfDef) WriteFieldBloomFilterFpChance(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField33(oprot)
}

func (p *CfDef) TStructName() string {
	return "CfDef"
}

func (p *CfDef) ThriftName() string {
	return "CfDef"
}

func (p *CfDef) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CfDef(%+v)", *p)
}

func (p *CfDef) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*CfDef)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *CfDef) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Keyspace
	case 2:
		return p.Name
	case 3:
		return p.ColumnType
	case 5:
		return p.ComparatorType
	case 6:
		return p.SubcomparatorType
	case 8:
		return p.Comment
	case 9:
		return p.RowCacheSize
	case 11:
		return p.KeyCacheSize
	case 12:
		return p.ReadRepairChance
	case 13:
		return p.ColumnMetadata
	case 14:
		return p.GcGraceSeconds
	case 15:
		return p.DefaultValidationClass
	case 16:
		return p.Id
	case 17:
		return p.MinCompactionThreshold
	case 18:
		return p.MaxCompactionThreshold
	case 19:
		return p.RowCacheSavePeriodInSeconds
	case 20:
		return p.KeyCacheSavePeriodInSeconds
	case 24:
		return p.ReplicateOnWrite
	case 25:
		return p.MergeShardsChance
	case 26:
		return p.KeyValidationClass
	case 27:
		return p.RowCacheProvider
	case 28:
		return p.KeyAlias
	case 29:
		return p.CompactionStrategy
	case 30:
		return p.CompactionStrategyOptions
	case 31:
		return p.RowCacheKeysToSave
	case 32:
		return p.CompressionOptions
	case 33:
		return p.BloomFilterFpChance
	}
	return nil
}

func (p *CfDef) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("keyspace", thrift.STRING, 1),
		thrift.NewTField("name", thrift.STRING, 2),
		thrift.NewTField("column_type", thrift.STRING, 3),
		thrift.NewTField("comparator_type", thrift.STRING, 5),
		thrift.NewTField("subcomparator_type", thrift.STRING, 6),
		thrift.NewTField("comment", thrift.STRING, 8),
		thrift.NewTField("row_cache_size", thrift.DOUBLE, 9),
		thrift.NewTField("key_cache_size", thrift.DOUBLE, 11),
		thrift.NewTField("read_repair_chance", thrift.DOUBLE, 12),
		thrift.NewTField("column_metadata", thrift.LIST, 13),
		thrift.NewTField("gc_grace_seconds", thrift.I32, 14),
		thrift.NewTField("default_validation_class", thrift.STRING, 15),
		thrift.NewTField("id", thrift.I32, 16),
		thrift.NewTField("min_compaction_threshold", thrift.I32, 17),
		thrift.NewTField("max_compaction_threshold", thrift.I32, 18),
		thrift.NewTField("row_cache_save_period_in_seconds", thrift.I32, 19),
		thrift.NewTField("key_cache_save_period_in_seconds", thrift.I32, 20),
		thrift.NewTField("replicate_on_write", thrift.BOOL, 24),
		thrift.NewTField("merge_shards_chance", thrift.DOUBLE, 25),
		thrift.NewTField("key_validation_class", thrift.STRING, 26),
		thrift.NewTField("row_cache_provider", thrift.STRING, 27),
		thrift.NewTField("key_alias", thrift.BINARY, 28),
		thrift.NewTField("compaction_strategy", thrift.STRING, 29),
		thrift.NewTField("compaction_strategy_options", thrift.MAP, 30),
		thrift.NewTField("row_cache_keys_to_save", thrift.I32, 31),
		thrift.NewTField("compression_options", thrift.MAP, 32),
		thrift.NewTField("bloom_filter_fp_chance", thrift.DOUBLE, 33),
	})
}

/**
 * Attributes:
 *  - Name
 *  - StrategyClass
 *  - StrategyOptions
 *  - ReplicationFactor: @deprecated
 *  - CfDefs
 *  - DurableWrites
 */
type KsDef struct {
	thrift.TStruct
	Name              string       "name"               // 1
	StrategyClass     string       "strategy_class"     // 2
	StrategyOptions   thrift.TMap  "strategy_options"   // 3
	ReplicationFactor int32        "replication_factor" // 4
	CfDefs            thrift.TList "cf_defs"            // 5
	DurableWrites     bool         "durable_writes"     // 6
}

func NewKsDef() *KsDef {
	output := &KsDef{
		TStruct: thrift.NewTStruct("KsDef", []thrift.TField{
			thrift.NewTField("name", thrift.STRING, 1),
			thrift.NewTField("strategy_class", thrift.STRING, 2),
			thrift.NewTField("strategy_options", thrift.MAP, 3),
			thrift.NewTField("replication_factor", thrift.I32, 4),
			thrift.NewTField("cf_defs", thrift.LIST, 5),
			thrift.NewTField("durable_writes", thrift.BOOL, 6),
		}),
	}
	{
		output.DurableWrites = true
	}
	return output
}

func (p *KsDef) IsSetStrategyOptions() bool {
	return p.StrategyOptions != nil && p.StrategyOptions.Len() > 0
}

func (p *KsDef) IsSetReplicationFactor() bool {
	return p.ReplicationFactor != 0
}

func (p *KsDef) IsSetDurableWrites() bool {
	return p.DurableWrites != true
}

func (p *KsDef) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "name" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "strategy_class" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "strategy_options" {
			if fieldTypeId == thrift.MAP {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "replication_factor" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 5 || fieldName == "cf_defs" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField5(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 6 || fieldName == "durable_writes" {
			if fieldTypeId == thrift.BOOL {
				err = p.ReadField6(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField6(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *KsDef) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v338, err339 := iprot.ReadString()
	if err339 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "name", p.ThriftName(), err339)
	}
	p.Name = v338
	return err
}

func (p *KsDef) ReadFieldName(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *KsDef) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v340, err341 := iprot.ReadString()
	if err341 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "strategy_class", p.ThriftName(), err341)
	}
	p.StrategyClass = v340
	return err
}

func (p *KsDef) ReadFieldStrategyClass(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *KsDef) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype345, _vtype346, _size344, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.StrategyOptions", "", err)
	}
	p.StrategyOptions = thrift.NewTMap(_ktype345, _vtype346, _size344)
	for _i348 := 0; _i348 < _size344; _i348++ {
		v351, err352 := iprot.ReadString()
		if err352 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key349", "", err352)
		}
		_key349 := v351
		v353, err354 := iprot.ReadString()
		if err354 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val350", "", err354)
		}
		_val350 := v353
		p.StrategyOptions.Set(_key349, _val350)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *KsDef) ReadFieldStrategyOptions(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *KsDef) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v355, err356 := iprot.ReadI32()
	if err356 != nil {
		return thrift.NewTProtocolExceptionReadField(4, "replication_factor", p.ThriftName(), err356)
	}
	p.ReplicationFactor = v355
	return err
}

func (p *KsDef) ReadFieldReplicationFactor(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *KsDef) ReadField5(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype362, _size359, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.CfDefs", "", err)
	}
	p.CfDefs = thrift.NewTList(_etype362, _size359)
	for _i363 := 0; _i363 < _size359; _i363++ {
		_elem364 := NewCfDef()
		err367 := _elem364.Read(iprot)
		if err367 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem364CfDef", err367)
		}
		p.CfDefs.Push(_elem364)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *KsDef) ReadFieldCfDefs(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField5(iprot)
}

func (p *KsDef) ReadField6(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v368, err369 := iprot.ReadBool()
	if err369 != nil {
		return thrift.NewTProtocolExceptionReadField(6, "durable_writes", p.ThriftName(), err369)
	}
	p.DurableWrites = v368
	return err
}

func (p *KsDef) ReadFieldDurableWrites(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField6(iprot)
}

func (p *KsDef) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("KsDef")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField5(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField6(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *KsDef) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("name", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Name))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "name", p.ThriftName(), err)
	}
	return err
}

func (p *KsDef) WriteFieldName(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *KsDef) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("strategy_class", thrift.STRING, 2)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "strategy_class", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.StrategyClass))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "strategy_class", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "strategy_class", p.ThriftName(), err)
	}
	return err
}

func (p *KsDef) WriteFieldStrategyClass(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *KsDef) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.StrategyOptions != nil {
		if p.IsSetStrategyOptions() {
			err = oprot.WriteFieldBegin("strategy_options", thrift.MAP, 3)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(3, "strategy_options", p.ThriftName(), err)
			}
			err = oprot.WriteMapBegin(thrift.STRING, thrift.STRING, p.StrategyOptions.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
			}
			for Miter370 := range p.StrategyOptions.Iter() {
				Kiter371, Viter372 := Miter370.Key().(string), Miter370.Value().(string)
				err = oprot.WriteString(string(Kiter371))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Kiter371", "", err)
				}
				err = oprot.WriteString(string(Viter372))
				if err != nil {
					return thrift.NewTProtocolExceptionWriteField(0, "Viter372", "", err)
				}
			}
			err = oprot.WriteMapEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(3, "strategy_options", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *KsDef) WriteFieldStrategyOptions(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *KsDef) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetReplicationFactor() {
		err = oprot.WriteFieldBegin("replication_factor", thrift.I32, 4)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "replication_factor", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.ReplicationFactor))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "replication_factor", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(4, "replication_factor", p.ThriftName(), err)
		}
	}
	return err
}

func (p *KsDef) WriteFieldReplicationFactor(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *KsDef) WriteField5(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.CfDefs != nil {
		err = oprot.WriteFieldBegin("cf_defs", thrift.LIST, 5)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(5, "cf_defs", p.ThriftName(), err)
		}
		err = oprot.WriteListBegin(thrift.STRUCT, p.CfDefs.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		for Iter373 := range p.CfDefs.Iter() {
			Iter374 := Iter373.(*CfDef)
			err = Iter374.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("CfDef", err)
			}
		}
		err = oprot.WriteListEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(5, "cf_defs", p.ThriftName(), err)
		}
	}
	return err
}

func (p *KsDef) WriteFieldCfDefs(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField5(oprot)
}

func (p *KsDef) WriteField6(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetDurableWrites() {
		err = oprot.WriteFieldBegin("durable_writes", thrift.BOOL, 6)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(6, "durable_writes", p.ThriftName(), err)
		}
		err = oprot.WriteBool(bool(p.DurableWrites))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(6, "durable_writes", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(6, "durable_writes", p.ThriftName(), err)
		}
	}
	return err
}

func (p *KsDef) WriteFieldDurableWrites(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField6(oprot)
}

func (p *KsDef) TStructName() string {
	return "KsDef"
}

func (p *KsDef) ThriftName() string {
	return "KsDef"
}

func (p *KsDef) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KsDef(%+v)", *p)
}

func (p *KsDef) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*KsDef)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *KsDef) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Name
	case 2:
		return p.StrategyClass
	case 3:
		return p.StrategyOptions
	case 4:
		return p.ReplicationFactor
	case 5:
		return p.CfDefs
	case 6:
		return p.DurableWrites
	}
	return nil
}

func (p *KsDef) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("name", thrift.STRING, 1),
		thrift.NewTField("strategy_class", thrift.STRING, 2),
		thrift.NewTField("strategy_options", thrift.MAP, 3),
		thrift.NewTField("replication_factor", thrift.I32, 4),
		thrift.NewTField("cf_defs", thrift.LIST, 5),
		thrift.NewTField("durable_writes", thrift.BOOL, 6),
	})
}

/**
 * Row returned from a CQL query
 * 
 * Attributes:
 *  - Key
 *  - Columns
 */
type CqlRow struct {
	thrift.TStruct
	Key     []byte       "key"     // 1
	Columns thrift.TList "columns" // 2
}

func NewCqlRow() *CqlRow {
	output := &CqlRow{
		TStruct: thrift.NewTStruct("CqlRow", []thrift.TField{
			thrift.NewTField("key", thrift.BINARY, 1),
			thrift.NewTField("columns", thrift.LIST, 2),
		}),
	}
	{
	}
	return output
}

func (p *CqlRow) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "key" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "columns" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CqlRow) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v375, err376 := iprot.ReadBinary()
	if err376 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "key", p.ThriftName(), err376)
	}
	p.Key = v375
	return err
}

func (p *CqlRow) ReadFieldKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *CqlRow) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype382, _size379, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Columns", "", err)
	}
	p.Columns = thrift.NewTList(_etype382, _size379)
	for _i383 := 0; _i383 < _size379; _i383++ {
		_elem384 := NewColumn()
		err387 := _elem384.Read(iprot)
		if err387 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem384Column", err387)
		}
		p.Columns.Push(_elem384)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *CqlRow) ReadFieldColumns(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *CqlRow) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("CqlRow")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CqlRow) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Key != nil {
		err = oprot.WriteFieldBegin("key", thrift.BINARY, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
		err = oprot.WriteBinary(p.Key)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CqlRow) WriteFieldKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *CqlRow) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Columns != nil {
		err = oprot.WriteFieldBegin("columns", thrift.LIST, 2)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "columns", p.ThriftName(), err)
		}
		err = oprot.WriteListBegin(thrift.STRUCT, p.Columns.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		for Iter388 := range p.Columns.Iter() {
			Iter389 := Iter388.(*Column)
			err = Iter389.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("Column", err)
			}
		}
		err = oprot.WriteListEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "columns", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CqlRow) WriteFieldColumns(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *CqlRow) TStructName() string {
	return "CqlRow"
}

func (p *CqlRow) ThriftName() string {
	return "CqlRow"
}

func (p *CqlRow) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CqlRow(%+v)", *p)
}

func (p *CqlRow) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*CqlRow)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *CqlRow) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Key
	case 2:
		return p.Columns
	}
	return nil
}

func (p *CqlRow) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("key", thrift.BINARY, 1),
		thrift.NewTField("columns", thrift.LIST, 2),
	})
}

/**
 * Attributes:
 *  - NameTypes
 *  - ValueTypes
 *  - DefaultNameType
 *  - DefaultValueType
 */
type CqlMetadata struct {
	thrift.TStruct
	NameTypes        thrift.TMap "name_types"         // 1
	ValueTypes       thrift.TMap "value_types"        // 2
	DefaultNameType  string      "default_name_type"  // 3
	DefaultValueType string      "default_value_type" // 4
}

func NewCqlMetadata() *CqlMetadata {
	output := &CqlMetadata{
		TStruct: thrift.NewTStruct("CqlMetadata", []thrift.TField{
			thrift.NewTField("name_types", thrift.MAP, 1),
			thrift.NewTField("value_types", thrift.MAP, 2),
			thrift.NewTField("default_name_type", thrift.STRING, 3),
			thrift.NewTField("default_value_type", thrift.STRING, 4),
		}),
	}
	{
	}
	return output
}

func (p *CqlMetadata) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "name_types" {
			if fieldTypeId == thrift.MAP {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "value_types" {
			if fieldTypeId == thrift.MAP {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "default_name_type" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "default_value_type" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CqlMetadata) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype393, _vtype394, _size392, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.NameTypes", "", err)
	}
	p.NameTypes = thrift.NewTMap(_ktype393, _vtype394, _size392)
	for _i396 := 0; _i396 < _size392; _i396++ {
		v399, err400 := iprot.ReadBinary()
		if err400 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key397", "", err400)
		}
		_key397 := v399
		v401, err402 := iprot.ReadString()
		if err402 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val398", "", err402)
		}
		_val398 := v401
		p.NameTypes.Set(_key397, _val398)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *CqlMetadata) ReadFieldNameTypes(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *CqlMetadata) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype406, _vtype407, _size405, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.ValueTypes", "", err)
	}
	p.ValueTypes = thrift.NewTMap(_ktype406, _vtype407, _size405)
	for _i409 := 0; _i409 < _size405; _i409++ {
		v412, err413 := iprot.ReadBinary()
		if err413 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key410", "", err413)
		}
		_key410 := v412
		v414, err415 := iprot.ReadString()
		if err415 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val411", "", err415)
		}
		_val411 := v414
		p.ValueTypes.Set(_key410, _val411)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *CqlMetadata) ReadFieldValueTypes(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *CqlMetadata) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v416, err417 := iprot.ReadString()
	if err417 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "default_name_type", p.ThriftName(), err417)
	}
	p.DefaultNameType = v416
	return err
}

func (p *CqlMetadata) ReadFieldDefaultNameType(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *CqlMetadata) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v418, err419 := iprot.ReadString()
	if err419 != nil {
		return thrift.NewTProtocolExceptionReadField(4, "default_value_type", p.ThriftName(), err419)
	}
	p.DefaultValueType = v418
	return err
}

func (p *CqlMetadata) ReadFieldDefaultValueType(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *CqlMetadata) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("CqlMetadata")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CqlMetadata) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.NameTypes != nil {
		err = oprot.WriteFieldBegin("name_types", thrift.MAP, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name_types", p.ThriftName(), err)
		}
		err = oprot.WriteMapBegin(thrift.BINARY, thrift.STRING, p.NameTypes.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		for Miter420 := range p.NameTypes.Iter() {
			Kiter421, Viter422 := Miter420.Key().([]byte), Miter420.Value().(string)
			err = oprot.WriteBinary(Kiter421)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Kiter421", "", err)
			}
			err = oprot.WriteString(string(Viter422))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Viter422", "", err)
			}
		}
		err = oprot.WriteMapEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "name_types", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CqlMetadata) WriteFieldNameTypes(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *CqlMetadata) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.ValueTypes != nil {
		err = oprot.WriteFieldBegin("value_types", thrift.MAP, 2)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "value_types", p.ThriftName(), err)
		}
		err = oprot.WriteMapBegin(thrift.BINARY, thrift.STRING, p.ValueTypes.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		for Miter423 := range p.ValueTypes.Iter() {
			Kiter424, Viter425 := Miter423.Key().([]byte), Miter423.Value().(string)
			err = oprot.WriteBinary(Kiter424)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Kiter424", "", err)
			}
			err = oprot.WriteString(string(Viter425))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Viter425", "", err)
			}
		}
		err = oprot.WriteMapEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(2, "value_types", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CqlMetadata) WriteFieldValueTypes(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *CqlMetadata) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("default_name_type", thrift.STRING, 3)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "default_name_type", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.DefaultNameType))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "default_name_type", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(3, "default_name_type", p.ThriftName(), err)
	}
	return err
}

func (p *CqlMetadata) WriteFieldDefaultNameType(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *CqlMetadata) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("default_value_type", thrift.STRING, 4)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(4, "default_value_type", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.DefaultValueType))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(4, "default_value_type", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(4, "default_value_type", p.ThriftName(), err)
	}
	return err
}

func (p *CqlMetadata) WriteFieldDefaultValueType(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *CqlMetadata) TStructName() string {
	return "CqlMetadata"
}

func (p *CqlMetadata) ThriftName() string {
	return "CqlMetadata"
}

func (p *CqlMetadata) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CqlMetadata(%+v)", *p)
}

func (p *CqlMetadata) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*CqlMetadata)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *CqlMetadata) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.NameTypes
	case 2:
		return p.ValueTypes
	case 3:
		return p.DefaultNameType
	case 4:
		return p.DefaultValueType
	}
	return nil
}

func (p *CqlMetadata) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("name_types", thrift.MAP, 1),
		thrift.NewTField("value_types", thrift.MAP, 2),
		thrift.NewTField("default_name_type", thrift.STRING, 3),
		thrift.NewTField("default_value_type", thrift.STRING, 4),
	})
}

/**
 * Attributes:
 *  - TypeA1
 *  - Rows
 *  - Num
 *  - Schema
 */
type CqlResult struct {
	thrift.TStruct
	TypeA1 CqlResultType "type"   // 1
	Rows   thrift.TList  "rows"   // 2
	Num    int32         "num"    // 3
	Schema *CqlMetadata  "schema" // 4
}

func NewCqlResult() *CqlResult {
	output := &CqlResult{
		TStruct: thrift.NewTStruct("CqlResult", []thrift.TField{
			thrift.NewTField("type", thrift.I32, 1),
			thrift.NewTField("rows", thrift.LIST, 2),
			thrift.NewTField("num", thrift.I32, 3),
			thrift.NewTField("schema", thrift.STRUCT, 4),
		}),
	}
	{
		output.TypeA1 = math.MinInt32 - 1
	}
	return output
}

func (p *CqlResult) IsSetTypeA1() bool {
	return int64(p.TypeA1) != math.MinInt32-1
}

func (p *CqlResult) IsSetRows() bool {
	return p.Rows != nil && p.Rows.Len() > 0
}

func (p *CqlResult) IsSetNum() bool {
	return p.Num != 0
}

func (p *CqlResult) IsSetSchema() bool {
	return p.Schema != nil
}

func (p *CqlResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "type" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 2 || fieldName == "rows" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 3 || fieldName == "num" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField3(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else if fieldId == 4 || fieldName == "schema" {
			if fieldTypeId == thrift.STRUCT {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField4(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CqlResult) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v426, err427 := iprot.ReadI32()
	if err427 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "type", p.ThriftName(), err427)
	}
	p.TypeA1 = CqlResultType(v426)
	return err
}

func (p *CqlResult) ReadFieldType(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *CqlResult) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype433, _size430, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Rows", "", err)
	}
	p.Rows = thrift.NewTList(_etype433, _size430)
	for _i434 := 0; _i434 < _size430; _i434++ {
		_elem435 := NewCqlRow()
		err438 := _elem435.Read(iprot)
		if err438 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem435CqlRow", err438)
		}
		p.Rows.Push(_elem435)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *CqlResult) ReadFieldRows(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *CqlResult) ReadField3(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v439, err440 := iprot.ReadI32()
	if err440 != nil {
		return thrift.NewTProtocolExceptionReadField(3, "num", p.ThriftName(), err440)
	}
	p.Num = v439
	return err
}

func (p *CqlResult) ReadFieldNum(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField3(iprot)
}

func (p *CqlResult) ReadField4(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	p.Schema = NewCqlMetadata()
	err443 := p.Schema.Read(iprot)
	if err443 != nil {
		return thrift.NewTProtocolExceptionReadStruct("p.SchemaCqlMetadata", err443)
	}
	return err
}

func (p *CqlResult) ReadFieldSchema(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField4(iprot)
}

func (p *CqlResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("CqlResult")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField3(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField4(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *CqlResult) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetTypeA1() {
		err = oprot.WriteFieldBegin("type", thrift.I32, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "type", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.TypeA1))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "type", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "type", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CqlResult) WriteFieldType(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *CqlResult) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Rows != nil {
		if p.IsSetRows() {
			err = oprot.WriteFieldBegin("rows", thrift.LIST, 2)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "rows", p.ThriftName(), err)
			}
			err = oprot.WriteListBegin(thrift.STRUCT, p.Rows.Len())
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			for Iter444 := range p.Rows.Iter() {
				Iter445 := Iter444.(*CqlRow)
				err = Iter445.Write(oprot)
				if err != nil {
					return thrift.NewTProtocolExceptionWriteStruct("CqlRow", err)
				}
			}
			err = oprot.WriteListEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(2, "rows", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *CqlResult) WriteFieldRows(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *CqlResult) WriteField3(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.IsSetNum() {
		err = oprot.WriteFieldBegin("num", thrift.I32, 3)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "num", p.ThriftName(), err)
		}
		err = oprot.WriteI32(int32(p.Num))
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "num", p.ThriftName(), err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(3, "num", p.ThriftName(), err)
		}
	}
	return err
}

func (p *CqlResult) WriteFieldNum(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField3(oprot)
}

func (p *CqlResult) WriteField4(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Schema != nil {
		if p.IsSetSchema() {
			err = oprot.WriteFieldBegin("schema", thrift.STRUCT, 4)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "schema", p.ThriftName(), err)
			}
			err = p.Schema.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("CqlMetadata", err)
			}
			err = oprot.WriteFieldEnd()
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(4, "schema", p.ThriftName(), err)
			}
		}
	}
	return err
}

func (p *CqlResult) WriteFieldSchema(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField4(oprot)
}

func (p *CqlResult) TStructName() string {
	return "CqlResult"
}

func (p *CqlResult) ThriftName() string {
	return "CqlResult"
}

func (p *CqlResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CqlResult(%+v)", *p)
}

func (p *CqlResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*CqlResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *CqlResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.TypeA1
	case 2:
		return p.Rows
	case 3:
		return p.Num
	case 4:
		return p.Schema
	}
	return nil
}

func (p *CqlResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("type", thrift.I32, 1),
		thrift.NewTField("rows", thrift.LIST, 2),
		thrift.NewTField("num", thrift.I32, 3),
		thrift.NewTField("schema", thrift.STRUCT, 4),
	})
}

const VERSION = "19.20.0"

func init() {
}
