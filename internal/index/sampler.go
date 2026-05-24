package index

import (
	"time"
)

// SampleEntry represents a single sampled index entry with its byte offset
// and parsed timestamp, used for fast approximate seeking.
type SampleEntry struct {
	Offset    int64
	Timestamp time.Time
}

// Sampler holds a downsampled view of an index for quick binary search.
type Sampler struct {
	samples []SampleEntry
	step    int
}

// NewSampler creates a Sampler by picking every nth entry from entries.
// A step of 0 or 1 includes every entry.
func NewSampler(entries []Entry, step int) *Sampler {
	if step < 1 {
		step = 1
	}
	samples := make([]SampleEntry, 0, (len(entries)/step)+1)
	for i, e := range entries {
		if i%step == 0 {
			samples = append(samples, SampleEntry{
				Offset:    e.Offset,
				Timestamp: e.Timestamp,
			})
		}
	}
	return &Sampler{samples: samples, step: step}
}

// Len returns the number of samples held by the Sampler.
func (s *Sampler) Len() int { return len(s.samples) }

// Step returns the sampling step used when building the Sampler.
func (s *Sampler) Step() int { return s.step }

// NearestOffset returns the byte offset of the sample whose timestamp is
// closest to (but not after) t. If no sample precedes t, the first sample's
// offset is returned. Returns -1 when the Sampler is empty.
func (s *Sampler) NearestOffset(t time.Time) int64 {
	if len(s.samples) == 0 {
		return -1
	}
	best := s.samples[0].Offset
	for _, se := range s.samples {
		if se.Timestamp.After(t) {
			break
		}
		best = se.Offset
	}
	return best
}

// Samples returns a copy of the internal sample slice.
func (s *Sampler) Samples() []SampleEntry {
	out := make([]SampleEntry, len(s.samples))
	copy(out, s.samples)
	return out
}
