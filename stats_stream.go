package netstatd

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
)

// statsStreamFactory implements tcpassembly.StreamFactory
type statsStreamFactory struct {
	*NetStats
}

// statsStream will handle the actual decoding of stats requests.
type statsStream struct {
	net, transport                      gopacket.Flow
	bytes, packets, outOfOrder, skipped int64
	start, end                          time.Time
	sawStart, sawEnd                    bool
	streams                             int64
	*NetStats
}

func (factory *statsStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	s := &statsStream{
		net:       net,
		transport: transport,
		start:     time.Now(),
		NetStats:  factory.NetStats,
	}
	s.end = s.start
	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return s
}

func (s *statsStream) Reassembled(reassemblies []tcpassembly.Reassembly) {
	for _, reassembly := range reassemblies {
		if reassembly.Seen.Before(s.end) {
			s.outOfOrder++
		} else {
			s.end = reassembly.Seen
		}
		s.bytes += int64(len(reassembly.Bytes))
		s.packets += 1
		if reassembly.Skip > 0 {
			s.skipped += int64(reassembly.Skip)
		}
		s.sawStart = s.sawStart || reassembly.Start
		s.sawEnd = s.sawEnd || reassembly.End

		s.Mark(reassembly.Bytes)
	}
}

func (s *statsStream) ReassemblyComplete() {
	//do nothing
}
