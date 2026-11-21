package usecase

import (
	"encoding/json"
	"github.com/stayfatal/Networking_Agent_Layer/internal/domain"
	"log"
)

type Generator interface {
	Generate(domain.Request) ([]string, error)
}
type Transport interface {
	Send(domain.Request, string, int, int, []byte) error
	Failure(domain.Request, error) error
}
type Service struct {
	generator Generator
	transport Transport
}

func New(generator Generator, transport Transport) *Service {
	return &Service{generator: generator, transport: transport}
}
func (s *Service) Process(req domain.Request) {
	err := s.transport.Send(req, "description", 0, 0, []byte(req.ProductDescription))
	if err == nil {
		var names []string
		names, err = s.generator.Generate(req)
		if err == nil {
			total := (len(names) + 2) / 3
			for i, start := 0, 0; start < len(names); i++ {
				end := start + 3
				if end > len(names) {
					end = len(names)
				}
				if len(names)-end == 1 {
					end--
				}
				b, _ := json.Marshal(names[start:end])
				if err = s.transport.Send(req, "names", i+1, total, b); err != nil {
					break
				}
				start = end
			}
		}
	}
	if err != nil {
		log.Print("process: ", err)
		e := s.transport.Failure(req, err)
		if e != nil {
			log.Print("failure delivery: ", e)
		}
	}
}
