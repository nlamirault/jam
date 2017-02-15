package waveout

import (
	"errors"
	"syscall"
	"time"
	"unsafe"
)

var (
	// ErrLessChunks means shortage of chunks.
	ErrLessChunks = errors.New("less chunks")
)

// Player is PCM player
type Player struct {
	h         syscall.Handle
	d         time.Duration
	chunks    []*chunk
	nextChunk int
}

type chunk struct {
	wh  WaveHdr
	buf []byte
}

// New creates a new Player instance.
func New(channels, samplesPerSec, bitsPerSample uint) (p *Player, err error) {
	ba := channels * bitsPerSample / 8
	f := WaveFormatEx{
		FormatTag:      WAVE_FORMAT_PCM,
		Channels:       uint16(channels),
		SamplesPerSec:  uint32(samplesPerSec),
		BitsPerSample:  uint16(bitsPerSample),
		BlockAlign:     uint16(ba),
		AvgBytesPerSec: uint32(samplesPerSec * ba),
	}
	var h syscall.Handle
	r := Open(&h, WAVE_MAPPER, &f, 0, 0, CALLBACK_NULL)
	if r != 0 {
		return nil, r
	}
	return &Player{
		h: h,
		d: time.Second / time.Duration(samplesPerSec * ba),
	}, nil
}

// NewWithBuffers creates a new Player instance with buffers.
func NewWithBuffers(channels, samplesPerSec, bitsPerSample uint, num, size int) (p *Player, err error) {
	p, err = New(channels, samplesPerSec, bitsPerSample)
	if err != nil {
		return nil, err
	}
	err = p.AddBuffers(num, size)
	if err != nil {
		p.Close()
		return nil, err
	}
	return p, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AddBuffer adds a buffer to the Player.
func (p *Player) AddBuffer(size int) error {
	buf := make([]byte, size)
	c := &chunk{
		wh: WaveHdr{
			Data:         &buf[0],
			BufferLength: uint32(size),
			Flags:        WHDR_BEGINLOOP | WHDR_ENDLOOP,
			Loops:        1,
		},
		buf: buf,
	}
	r := PrepareHeader(p.h, &c.wh, uint32(unsafe.Sizeof(c.wh)))
	if r != 0 {
		return r
	}
	p.chunks = append(p.chunks, c)
	return nil
}

// AddBuffers adds multiple buffers to the player.
func (p *Player) AddBuffers(num, size int) error {
	for i := 0; i < num; i++ {
		err := p.AddBuffer(size)
		if err != nil {
			return err
		}
	}
	return nil
}

// Write outputs PCM sound.
func (p *Player) Write(b []byte) (n int, err error) {
	for len(b) > 0 {
		c, err := p.getNextChunk()
		if err != nil {
			return n, err
		}
		// copy data to buffer of a chunk.
		l := min(len(b), len(c.buf))
		copy(c.buf[:l], b[:l])
		c.wh.BufferLength = uint32(l)
		b = b[l:]
		// queue a chunk to output as sound.
		r := Write(p.h, &c.wh, uint32(unsafe.Sizeof(c.wh)))
		if r != 0 {
			return n, r
		}
		n += l
	}
	return n, nil
}

func (p *Player) getNextChunk() (*chunk, error) {
	if len(p.chunks) == 0 {
		return nil, ErrLessChunks
	}
	c := p.chunks[p.nextChunk]
	p.nextChunk++
	if p.nextChunk >= len(p.chunks) {
		p.nextChunk = 0
	}
	d := p.d
	for {
		if c.wh.Flags&WHDR_INQUEUE == 0 {
			break
		}
		// auto adjust.
		d /= 2
		if d < time.Millisecond {
			d = time.Millisecond
		}
		time.Sleep(d)
	}
	return c, nil
}

// Close closes a Player.
func (p *Player) Close() (err error) {
	if p == nil {
		return nil
	}
	if p.h == 0 {
		return nil
	}
	Reset(p.h)
	for _, c := range p.chunks {
		r := UnprepareHeader(p.h, &c.wh, uint32(unsafe.Sizeof(c.wh)))
		if err != nil && r != 0 {
			err = r
		}
	}
	p.chunks = nil
	r := Close(p.h)
	if err != nil && r != 0 {
		err = r
	}
	p.h = 0
	return err
}

// SetVolume changes volume of the player.
func (p *Player) SetVolume(left, right uint16) error {
	v := uint32(right)<<16 + uint32(left)
	r := SetVolume(p.h, v)
	if r != 0 {
		return r
	}
	return nil
}

// Volume obtains volume of the player.
func (p *Player) Volume() (left, right uint16, err error) {
	var v uint32
	r := GetVolume(p.h, &v)
	if r != 0 {
		return 0, 0, r
	}
	return uint16(v), uint16(v >> 16), nil
}

// Pause pauses sound output.  It can be resumed with Restart().
func (p *Player) Pause() error {
	r := Pause(p.h)
	if r != 0 {
		return r
	}
	return nil
}

// Restart resumes sound output.
func (p *Player) Restart() error {
	r := Restart(p.h)
	if r != 0 {
		return r
	}
	return nil
}

// Reset stops sound.
func (p *Player) Reset() error {
	r := Reset(p.h)
	if r != 0 {
		return r
	}
	return nil
}

// Wait waits end of sound.
func (p *Player) Wait() error {
	if len(p.chunks) == 0 {
		return nil
	}
	for {
		// count inqueued chunks.
		n := 0
		for _, c := range p.chunks {
			if c.wh.Flags&WHDR_INQUEUE != 0 {
				n++
				break
			}
		}
		if n == 0 {
			return nil
		}
		time.Sleep(1 * time.Millisecond)
	}
}
