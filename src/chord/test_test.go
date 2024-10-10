package chord

import (
	"fmt"
	"testing"
	"time"
)

func naiveHash(val int) int {
	return ((val * 1919810) / 114514) % 65536
}

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
	cfg := make_config(t, 5)
	time.Sleep(time.Duration(500)*time.Millisecond)
	if cfg.chords[0].hash("hello") != cfg.chords[1].hash("hello"){
		t.Errorf("hash expected to be the same, but got different")
	}
}

func TestPutGet1 (t *testing.T) {
	cfg := make_config(t, 5)
	defer cfg.cleanup()
	time.Sleep(time.Duration(500)*time.Millisecond)
	cfg.chords[4].put("apple","red")
	cfg.chords[1].put("banana","yellow")
	cfg.chords[3].put("pear","green")
	a1, err1 := cfg.chords[2].get("apple")
	a2, err2 := cfg.chords[0].get("banana")
	a3, err3 := cfg.chords[1].get("pear")
	if (err1 != nil || err2 != nil || err3 != nil) {
		t.Errorf("get failed")
	}

	if a1 != "red" {
		t.Errorf("expect 'red', but got '%s'", a1)
	}
	if a2 != "yellow" {
		t.Errorf("expect 'yellow', but got '%s'", a2)
	}
	if a3 != "green" {
		t.Errorf("expect 'green', but got '%s'", a3)
	}
	_, err := cfg.chords[0].get("nokey")
	if (err == nil) {
		t.Errorf("expect return a error.")
	}
	time.Sleep(time.Duration(500)*time.Millisecond)

}

func TestFindRightNode(t *testing.T) {
	cfg := make_config(t, 5)


}