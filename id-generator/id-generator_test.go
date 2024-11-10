package idgenerator

import (
	"log"
	"os"
	"testing"
)

func Test_Simple(t *testing.T) {
	ig, err := CreateIdGenerator(os.Getenv("TEST_DB"))
	if nil != err {
		t.Fatal(err)
	}
	ig.SetPrefix("TT").SetLength(8)
	id, err := ig.GetId()
	if nil != err {
		t.Fatal(err)
	}
	t.Log(id)

	id, err = ig.GetDatetimeId()
	if nil != err {
		t.Fatal(err)
	}
	t.Log(id)
}

func Test_MultipleDatetime(t *testing.T) {
	ig, err := CreateIdGenerator(os.Getenv("TEST_DB"))
	if nil != err {
		t.Fatal(err)
	}
	ig.SetPrefix("TT").SetLength(8)

	ch := make(chan int)
	for i := 0; i < 10; i++ {
		go func() {
			id, err := ig.GetDatetimeId()
			if nil != err {
				log.Println("go routine get failed, ", err)
				ch <- -1
			} else {
				ch <- 0
				t.Log(id)
			}
		}()
	}

	for i := 0; i < 10; i++ {
		<-ch
	}
}

func Benchmark_Simple(b *testing.B) {
	ig, err := CreateIdGenerator(os.Getenv("TEST_DB"))
	if nil != err {
		b.Fatal(err)
	}
	ig.SetPrefix("TT").SetLength(8)
	b.ResetTimer()
	for i := 0; i < 400; i++ {
		_, err = ig.GetId()
		if nil != err {
			b.Fatal(err)
		}
	}
}
