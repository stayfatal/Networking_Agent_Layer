package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/stayfatal/Networking_Agent_Layer/internal/domain"
	"testing"
)

type generatorStub struct {
	names []string
	err   error
}

func (g generatorStub) Generate(domain.Request) ([]string, error) { return g.names, g.err }

type packet struct {
	kind         string
	index, total int
	data         []byte
}
type transportStub struct {
	packets []packet
	failure error
}

func (t *transportStub) Send(_ domain.Request, kind string, index, total int, data []byte) error {
	t.packets = append(t.packets, packet{kind, index, total, data})
	return nil
}
func (t *transportStub) Failure(_ domain.Request, err error) error { t.failure = err; return nil }
func TestProcessChunks(t *testing.T) {
	for n := 5; n <= 10; n++ {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			names := make([]string, n)
			for i := range names {
				names[i] = fmt.Sprint(i)
			}
			tr := &transportStub{}
			New(generatorStub{names: names}, tr).Process(domain.Request{ProductDescription: "товар"})
			total := (n + 2) / 3
			if len(tr.packets) != total+1 || tr.packets[0].kind != "description" {
				t.Fatalf("packets: %+v", tr.packets)
			}
			var got []string
			for i, p := range tr.packets[1:] {
				var chunk []string
				if err := json.Unmarshal(p.data, &chunk); err != nil {
					t.Fatal(err)
				}
				if p.kind != "names" || p.index != i+1 || p.total != total || len(chunk) < 2 || len(chunk) > 3 {
					t.Fatalf("invalid chunk: %+v %v", p, chunk)
				}
				got = append(got, chunk...)
			}
			for i := range names {
				if got[i] != names[i] {
					t.Fatalf("got %v", got)
				}
			}
		})
	}
}
func TestProcessReportsGenerationFailure(t *testing.T) {
	tr := &transportStub{}
	err := fmt.Errorf("unavailable")
	New(generatorStub{err: err}, tr).Process(domain.Request{})
	if tr.failure != err || len(tr.packets) != 1 {
		t.Fatalf("failure=%v packets=%v", tr.failure, tr.packets)
	}
}
