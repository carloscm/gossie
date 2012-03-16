package gossie

import (
    "testing"
    "reflect"
    "os"
)

/*

todo:

    basically everything. real unit testing for all struct funcs

*/

type errNoMeta struct {
    a int
}
type errNoMetaKeyColVal struct {
    a int `cf:"cfname"`
}
type errNoMetaColVal struct {
    a int `cf:"cfname" key:"a"`
}
type errNoMetaVal struct {
    a   int `cf:"cfname" key:"a" col:"b"`
    b   int
}
type errInvKey struct {
    a   int `cf:"cfname" key:"z" col:"b" val:"c"`
    b   int
    c   int
}
type errInvCol struct {
    a   int `cf:"cfname" key:"a" col:"z" val:"c"`
    b   int
    c   int
}
type errInvVal struct {
    a   int `cf:"cfname" key:"a" col:"b" val:"z"`
    b   int
    c   int
}
type noErrA struct {
    a   int `cf:"cfname" key:"a" col:"b" val:"c"`
    b   int
    c   int
}
type noErrB struct {
    a   int `cf:"cfname" key:"a" col:"*name" val:"*value"`
    b   int
    c   int
}
type noErrC struct {
    a   int `cf:"cfname" key:"a" col:"b,*name" val:"*value"`
    b   int
    c   int
}

func buildMappingFromPtr(instance interface{}) (*structMapping, os.Error) {
    valuePtr := reflect.ValueOf(instance)
    value := reflect.Indirect(valuePtr)
    typ := value.Type()
    return newStructMapping(typ)
}

func structMapMustError(t *testing.T, instance interface{}) {
    _, err := buildMappingFromPtr(instance)
    if err == nil {
        t.Error("Expected error calling newStructMapping, got none")
    }
}

func TestStructMapping(t *testing.T) {
    var sampleInt int
    sampleIntT := reflect.TypeOf(sampleInt)

    structMapMustError(t, &errNoMeta{})
    structMapMustError(t, &errNoMetaKeyColVal{})
    structMapMustError(t, &errNoMetaColVal{})
    structMapMustError(t, &errNoMetaVal{})
    structMapMustError(t, &errInvKey{})
    structMapMustError(t, &errInvCol{})
    structMapMustError(t, &errInvVal{})

    mapA, _ := buildMappingFromPtr(&noErrA{1, 2, 3})
    goodA := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "a", goType: sampleIntT, cassandraType: LongType},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: baseTypeField, position: 1, name: "b", goType: sampleIntT, cassandraType: LongType},
        },
        value:             &fieldMapping{fieldKind: baseTypeField, position: 2, name: "c", goType: sampleIntT, cassandraType: LongType},
        others:            nil,
        isCompositeColumn: false,
    }
    if !reflect.DeepEqual(mapA, goodA) {
        t.Error("Mapping for struct sample A does not match expected output, ", mapA, " vs ", goodA)
    }

    mapB, _ := buildMappingFromPtr(&noErrB{1, 2, 3})
    goodB := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "a", goType: sampleIntT, cassandraType: LongType},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: starNameField, position: 0, name: "", goType: nil, cassandraType: 0},
        },
        value: &fieldMapping{fieldKind: starValueField, position: 0, name: "", goType: nil, cassandraType: 0},
        others: map[string]*fieldMapping{
            "b": &fieldMapping{fieldKind: baseTypeField, position: 1, name: "b", goType: sampleIntT, cassandraType: LongType},
            "c": &fieldMapping{fieldKind: baseTypeField, position: 2, name: "c", goType: sampleIntT, cassandraType: LongType},
        },
        isCompositeColumn: false,
    }
    if !reflect.DeepEqual(mapB, goodB) {
        t.Error("Mapping for struct sample B does not match expected output, ", mapB, " vs ", goodB)
    }

    mapC, _ := buildMappingFromPtr(&noErrC{1, 2, 3})
    goodC := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "a", goType: sampleIntT, cassandraType: LongType},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: baseTypeField, position: 1, name: "b", goType: sampleIntT, cassandraType: LongType},
            &fieldMapping{fieldKind: starNameField, position: 0, name: "", goType: nil, cassandraType: 0},
        },
        value: &fieldMapping{fieldKind: starValueField, position: 0, name: "", goType: nil, cassandraType: 0},
        others: map[string]*fieldMapping{
            "c": &fieldMapping{fieldKind: baseTypeField, position: 2, name: "c", goType: sampleIntT, cassandraType: LongType},
        },
        isCompositeColumn: true,
    }
    if !reflect.DeepEqual(mapC, goodC) {
        t.Error("Mapping for struct sample C does not match expected output, ", mapC, " vs ", goodC)
    }
}

type timeline struct {
    UserId  string `cf:"Timelines" key:"UserId" col:"TweetId,*name" val:"*value"`
    TweetId int
    Author  string
    Body    string
}

func TestMap(t *testing.T) {

    tweet := &timeline{UserId: "abc", TweetId: 3, Author: "xyz", Body: "hello world"}

    row, _ := Map(tweet)

    if len(row.Columns) != 2 {
        t.Error("Expected number of columns is 2, got ", len(row.Columns))
    }

    t.Log(row.Columns[0].Name)
    t.Log(row.Columns[0].Value)
    t.Log(row.Columns[1].Name)
    t.Log(row.Columns[1].Value)

    //t.Fatal("heh")

}
