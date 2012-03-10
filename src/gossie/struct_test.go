package gossie

import (
    "testing"
)

/*

todo:

    basically everything. real unit testing for all struct funcs

*/

type Timeline struct {
    UserId  string `cf:"Timelines" key:"UserId" col:"TweetId,*name" val:"*value"`
    TweetId int
    Author  string
    Body    string
}

func TestMap(t *testing.T) {

    tweet := &Timeline{UserId: "abc", TweetId: 3, Author: "xyz", Body: "hello world"}

    row := Map(tweet)

    if len(row.Columns) != 2 {
        t.Error("Expected number of columns is 2, got ", len(row.Columns))
    }

    t.Log(row.Columns[0].Name)
    t.Log(row.Columns[0].Value)
    t.Log(row.Columns[1].Name)
    t.Log(row.Columns[1].Value)

    t.Fatal("heh")

}
