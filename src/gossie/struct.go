package gossie

import (
    "reflect"
    "strings"
    "sync"
    "os"
    "fmt"
)

/*

this is very much WIP but past 50% completion

todo:

    error checking and passing in mapField
    composite marshaling

---

in CQL:
CREATE TABLE timeline (
    user_id varchar,
    tweet_id uuid,
    author varchar,
    body varchar,
    PRIMARY KEY (user_id, tweet_id)
);

represents a single column in a single row
type TimelineTweetAtt struct  {
    UserId      string      "cf:Timelines key:UserId col:TweetId,Att val:Content"
    TweetId     UUID
    Att         string
    Content     string
}

represents a range of columns with single values for the first composite and the row key, and a range over the second comp,
storing only the comp name, not the col value
type TimelineTweetAttNames struct  {
    UserId      string      "cf:Timelines key:UserId col:TweetId,Atts val:"
    TweetId     UUID
    Atts        []string
}

same but with a second slice holding the values
type TimelineTweetAttNames struct  {
    UserId      string      "cf:Timelines key:UserId col:TweetId,AttNames val:AttValues"
    TweetId     UUID
    AttNames    []string
    AttValues   []string
}



another way for mapping a single tweet with struct fields instead of slices
*name means "any field name inside the struct that is not mentioned in the mapping for other purpose"
*value means "the value of the struct field for the selected *name mapping"
type TimelineTweet struct  {
    UserId      string      "cf:Timelines key:UserId col:TweetId,*name val:*value"
    TweetId     UUID
    Author      string
    Body        string
}



type Timeline struct  {
    UserId      string      "cf:Timeline key:UserId col:Section,TweetId,*name val:*value"
    Section     int
    TweetId     UUID
    Author      string
    Body        string
}


mapField(source interface{}, row *Row, mapping *map, component int, composite []byte, value []byte)

    if components left

        switch type of field named by component
            case base type
                set value of the current composite field to the field value
                mapField(source, row, map, component+1, composite, value)

            case base type slice
                iterate slice
                    clone composite
                    set value of the current composite field to the current slice value
                    mapField(source, row, map, component+1, composite, value)

            case *name
                iterate over non-key/col/val-referenced struct fields
                    clone composite
                    set value of the current composite field to the struct field name
                    value := field value
                    mapField(source, row, map, component+1, composite, value)

    else

        if val: is *value
            use passed value
            emit column inside row with the composite plus value

        else if val: is a base type slice struct field??????? <-- very useful try to support for index-like CFs
            get the index used for the composite part slice and use it here???

        else is empty
            value is the logical zero for the CF default validator
            emit column inside row with the composite plus value

        else is constant value???
            value is constant read from mapping?
            emit column inside row with the composite plus value

*/

const (
    _   = iota
    baseTypeField
    baseTypeSliceField
    starNameField
    starValueField
)

// mapping stores how to map from/to a struct
type fieldMapping struct {
    fieldKind     int
    position      int
    name          string
    goType        reflect.Type
    cassandraType TypeDesc
}
type structMapping struct {
    cf      string
    key     *fieldMapping
    columns []*fieldMapping
    value   *fieldMapping
    others  []*fieldMapping
}

var mapCache map[reflect.Type]*structMapping
var mapCacheMutex *sync.Mutex = new(sync.Mutex)

func isMarshalable(k reflect.Kind) bool {
    return k == reflect.Bool ||
        k == reflect.Int ||
        k == reflect.Int8 ||
        k == reflect.Int16 ||
        k == reflect.Int32 ||
        k == reflect.Int64 ||
        k == reflect.Float32 ||
        k == reflect.Float64 ||
        k == reflect.Array ||
        k == reflect.Slice ||
        k == reflect.String
}

func defaultCassandraType(t reflect.Type) (TypeDesc, int) {
    switch t.Kind() {
    case reflect.Bool:
        return BooleanType, baseTypeField
    case reflect.String:
        return UTF8Type, baseTypeField
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return LongType, baseTypeField
    case reflect.Float32:
        return FloatType, baseTypeField
    case reflect.Float64:
        return DoubleType, baseTypeField
    case reflect.Array:
        if t.Name() == "UUID" && t.Size() == 16 {
            return UUIDType, baseTypeField
        }
        return UnknownType, baseTypeField
    case reflect.Slice:
        if et := t.Elem(); et.Kind() == reflect.Int8 {
            return BytesType, baseTypeField
        } else {
            if subTD, subKind := defaultCassandraType(et); subTD != UnknownType && subKind == baseTypeField {
                return subTD, baseTypeSliceField
            }
            return UnknownType, baseTypeField
        }
        return UnknownType, baseTypeField
    }
    return UnknownType, baseTypeField
}

func newFieldMapping(pos int, sf reflect.StructField) *fieldMapping {
    fm := &fieldMapping{}
    fm.cassandraType, fm.fieldKind = defaultCassandraType(sf.Type)
    fm.position = pos
    fm.name = sf.Name
    fm.goType = sf.Type
    return fm
}

func newStructMapping(t reflect.Type) (*structMapping, os.Error) {
    sm := &structMapping{}
    n := t.NumField()
    found := false
    // globally recognized meta fields in the struct tag
    meta := map[string]string{
        "cf":  "",
        "key": "",
        "col": "",
        "val": "",
    }
    // hold a field mapping for every candidate field
    fields := make(map[string]*fieldMapping)

    // pass 1: gather field metadata
    for i := 0; i < n; i++ {
        sf := t.Field(i)
        // find the field tags
        for key, _ := range meta {
            if tagValue := sf.Tag.Get(key); tagValue != "" {
                meta[key] = tagValue
            }
        }
        // build a field mapping for all non-anon named fields with a suitable Kind
        if sf.Name != "" && !sf.Anonymous && isMarshalable(sf.Type.Kind()) {
            fields[sf.Name] = newFieldMapping(i, sf)
        }
    }

    // pass 2: struct data for each meta field
    if name := meta["cf"]; name != "" {
        sm.cf = name
    }
    if name := meta["key"]; name != "" {
        if sm.key, found = fields[name]; !found {
            return nil, os.NewError(fmt.Sprint("Referenced key field ", name, " does not exist in struct ", t.Name()))
        }
        if sm.key.fieldKind != baseTypeField {
            return nil, os.NewError(fmt.Sprint("Referenced key field ", name, " in struct ", t.Name(), " has invalid type"))
        }
        fields[name] = nil, false
    } else {
        return nil, os.NewError(fmt.Sprint("No key field in struct ", t.Name()))
    }
    if name := meta["val"]; (name != "") || (name == "*value") {
        if name == "*value" {
            sm.value = &fieldMapping{fieldKind: starValueField}
        } else if sm.value, found = fields[name]; !found {
            return nil, os.NewError(fmt.Sprint("Referenced value field ", name, " does not exist in struct ", t.Name()))
        }
        fields[name] = nil, false
    } else {
        return nil, os.NewError(fmt.Sprint("No val field in struct ", t.Name()))
    }
    if meta["col"] != "" {
        colNames := strings.Split(meta["col"], ",")
        for _, name := range colNames {
            var fm *fieldMapping
            if name == "*name" {
                fm = &fieldMapping{fieldKind: starNameField}
            } else if fm, found = fields[name]; !found {
                return nil, os.NewError(fmt.Sprint("Referenced column field ", name, " does not exist in struct ", t.Name()))
            }
            fields[name] = nil, false
            sm.columns = append(sm.columns, fm)
        }
        for _, fm := range fields {
            sm.others = append(sm.others, fm)
        }
    } else {
        return nil, os.NewError(fmt.Sprint("No col field in struct ", t.Name()))
    }
    return sm, nil
}

func getMapping(v reflect.Value) (*structMapping, os.Error) {
    var sm *structMapping
    var err os.Error
    found := false
    t := v.Type()
    mapCacheMutex.Lock()
    if sm, found = mapCache[t]; !found {
        sm, err = newStructMapping(t)
    }
    mapCacheMutex.Unlock()
    return sm, err
}

func Map(source interface{}) *Row {

    // always work with a pointer to struct
    vp := reflect.ValueOf(source)
    if vp.Kind() != reflect.Ptr {
        return nil // error???
    }
    if vp.IsNil() {
        return nil // error???
    }
    v := reflect.Indirect(vp)
    if v.Kind() != reflect.Struct {
        return nil // error???
    }

    sm, _ := getMapping(v)
    row := &Row{}
    // TODO: marshal key
    mapField(v, row, sm, 0, make([]byte, 0), make([]byte, 0), 0)
    return row
}

func mapField(source reflect.Value, row *Row, sm *structMapping, component int, composite []byte, value []byte, valueIndex int) {

    // TODO error checking! marshal calls etc

    // check if there are components left
    if component < len(sm.columns) {

        fm := sm.columns[component]

        // switch type of field named by component
        switch fm.fieldKind {

        // base type
        case baseTypeField:
            // set value of the current composite field to the field value
            v := source.Field(fm.position)
            b, _ := Marshal(v.Interface(), fm.cassandraType)
            // !!!!!!!!!!!!!!!!!!!!!!!!!!!!
            // TODO: actual composite serialization goes here
            composite = b
            // !!!!!!!!!!!!!!!!!!!!!!!!!!!!
            mapField(source, row, sm, component+1, composite, value, valueIndex)

        // slice of base type
        /*
           case baseTypeSliceField:
        */

        // *name
        case starNameField:
            // iterate over non-key/col/val-referenced struct fields
            for i, fm := range sm.others {
                // set value of the current composite field to the field value
                b, _ := Marshal(fm.name, UTF8Type)
                // !!!!!!!!!!!!!!!!!!!!!!!!!!!!
                // TODO: actual composite serialization goes here
                subComposite := b
                // !!!!!!!!!!!!!!!!!!!!!!!!!!!!
                v := source.Field(fm.position)
                b, _ = Marshal(v.Interface(), fm.cassandraType)
                mapField(source, row, sm, component+1, subComposite, b, i)
            }
        }

    } else {
        // no components left, emit column

        fm := sm.value

        // switch type of value field
        switch fm.fieldKind {

        case starValueField:
            // use passed value
            row.Columns = append(row.Columns, &Column{Name: composite, Value: value})

        case baseTypeSliceField:
            // set value to the passed value index in this slice
            vs := source.Field(fm.position)
            v := vs.Index(valueIndex)
            b, _ := Marshal(v.Interface(), fm.cassandraType)
            row.Columns = append(row.Columns, &Column{Name: composite, Value: b})

            // support non-slice field case?

            // support literal case?

            // support zero, non-set case?

        }
    }
}

/*

mapField(source interface{}, row *Row, mapping *map, component int, composite []byte, value []byte)

    if components left

        switch type of field named by component
            case base type
                set value of the current composite field to the field value
                mapField(source, row, map, component+1, composite, value)

            case base type slice
                iterate slice
                    clone composite
                    set value of the current composite field to the current slice value
                    mapField(source, row, map, component+1, composite, value)

            case *name
                iterate over non-key/col/val-referenced struct fields
                    clone composite
                    set value of the current composite field to the struct field name
                    value := field value
                    mapField(source, row, map, component+1, composite, value)

    else

        if val: is *value
            use passed value
            emit column inside row with the composite plus value

        else if val: is a base type slice struct field??????? <-- very useful try to support for index-like CFs
            get the index used for the composite part slice and use it here???

        else is empty
            value is the logical zero for the CF default validator
            emit column inside row with the composite plus value

        else is constant value???
            value is constant read from mapping?
            emit column inside row with the composite plus value

*/
