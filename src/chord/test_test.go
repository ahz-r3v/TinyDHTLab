package chord

import (
	"fmt"
	"testing"
	"time"
)

func TestMakeNode (t *testing.T) {
	num := 5
	cfg := make_config(t, num)
	if  cfg.chords[0].Tellself(0,0) != cfg.chords[1].Tellothers(0,0)[0] {
		t.Errorf("chords[0] and chords[2].Tell[0] expected to be the same, but got different")
	}
	time.Sleep(time.Duration(1000)*time.Millisecond)
	for i:=0;i<num;i++{
		fmt.Println(cfg.chords[i].Location)
		fmt.Println(cfg.chords[i].PrevLocation)
		fmt.Println(cfg.chords[i].SuccLocation)
	}
	
}

func TestLocationRPC (t *testing.T) {
	cfg := make_config(t, 10)
	argsl := &GetLocationArgs{}
	replyl := &GetLocationReply{}
	cfg.chords[2].Peers[0].Call("Chord.GetLocation",argsl,replyl)
	l1 := replyl.Location
	cfg.chords[2].SendGetLocation(0,argsl,replyl)
	l2 := replyl.Location
	if cfg.chords[0].Tellself(0,0) != cfg.chords[2].Tellothers(0,0)[0] {
		t.Errorf("chords[0] and chords[2].Tell[0] expected to be the same, but got different")
	}
	if l1 != l2{
		t.Errorf("results of RPC should be the same but not")
	}
	fmt.Println(l1, l2)
	Time := 500
	time.Sleep(time.Duration(Time)*time.Millisecond)
}

func TestHash (t *testing.T) {
	if BKDRHash("hello") != BKDRHash("hello"){
		t.Errorf("hash expected to be the same, but got different")
	}
	fmt.Println(BKDRHash("hello"))
	fmt.Println(naiveHash(BKDRHash("hello")))
	fmt.Println(BKDRHash("hello1"))
	fmt.Println(naiveHash(BKDRHash("hello1")))
	fmt.Println(BKDRHash("...999"))
}

func TestPut1 (t *testing.T) {
	cfg := make_config(t, 5)
	defer cfg.cleanup()
	time.Sleep(time.Duration(500)*time.Millisecond)
	cfg.chords[4].put("apple","red")
	fmt.Println(cfg.chords[2].get("apple"))
	cfg.chords[1].put("banana","yellow")
	fmt.Println(cfg.chords[0].get("banana"))
	cfg.chords[3].put("pear","green")
	fmt.Println(cfg.chords[4].get("peer"))
	fmt.Println(cfg.chords[0].get("apple"))
	fmt.Println(cfg.chords[2].get("banana"))
	fmt.Println(cfg.chords[1].get("pear"))
	time.Sleep(time.Duration(500)*time.Millisecond)

}