package netstatd

import (
	"bufio"
	"io"
	"log"
	"net/http"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

// statsStreamFactory implements tcpassembly.StreamFactory
type statsStreamFactory struct {
	netStat *NetStat
}

// statsStream will handle the actual decoding of stats requests.
type statsStream struct {
	net, transport gopacket.Flow
	netStat        *NetStat
	r              tcpreader.ReaderStream
}

func (factory *statsStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	s := &statsStream{
		net:       net,
		transport: transport,
		netStat:   factory.netStat,
		r:         tcpreader.NewReaderStream(),
	}

	go s.parseHttpRequest()
	return s
}

func (s *statsStream) Reassembled(reassemblies []tcpassembly.Reassembly) {
	for _, reassembly := range reassemblies {
		s.netStat.MarkTcp(reassembly.Bytes)
	}

	s.r.Reassembled(reassemblies)
}

func (s *statsStream) ReassemblyComplete() {
	s.r.ReassemblyComplete()
}

func (s *statsStream) parseHttpRequest() {
	buf := bufio.NewReader(&s.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			log.Println("Error reading stream", s.net, s.transport, ":", err)
		} else {
			s.netStat.MarkHttpRequest(req)
		}
	}
}
