
package stream

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/looplab/fsm"
)

type FSM struct {
	Clock  func() time.Time
	FSM    *fsm.FSM
	Target chan string
}

func (s *Stream) PushEvent(ctx context.Context, str string) {
	log := logger.Entry(ctx)
	log.Tracef("pushEvent: %s", str)
	select {
	case s.fsm.Target <- str:
	case <-time.After(time.Second):
		log.Warn("pushEvent hung")
	}
	log.Tracef("pushEvent done: %s", str)
}

//go:generate sh -c "cd ../../ && go run ./... -dump-fsm | dot -Nmargin=0.8 -s144 -Tsvg /dev/stdin -o fsm.svg"
func (s *Stream) GetFSM() *fsm.FSM {
	return s.fsm.FSM
}

func (s *Stream) newFSM(ctx context.Context) *FSM {
	log := logger.Entry(ctx)
	f := FSM{
		FSM: fsm.NewFSM(
			"undefined",
			fsm.Events{
				{Name: "steady", Src: []string{"undefined", "down"}, Dst: "down"},
				{Name: "steady", Src: []string{"up"}, Dst: "up"},
				{Name: "steady", Src: []string{"going_up"}, Dst: "going_up"},
				{Name: "steady", Src: []string{"going_down"}, Dst: "going_down"},
				{Name: "steady_timer", Src: []string{"up"}, Dst: "going_down"},
				{Name: "steady_timer", Src: []string{"going_down", "down", "going_up"}, Dst: "down"},
				{Name: "unsteady", Src: []string{"undefined", "down", "going_up"}, Dst: "going_up"},
				{Name: "unsteady", Src: []string{"up", "going_down"}, Dst: "up"},
				{Name: "unsteady_timer", Src: []string{"going_up", "going_down", "up"}, Dst: "up"},
				{Name: "no_data", Src: []string{"undefined", "down"}, Dst: "undefined"},
				{Name: "no_data", Src: []string{"going_up"}, Dst: "going_up"},
				{Name: "no_data", Src: []string{"going_down", "up"}, Dst: "going_down"},
				{Name: "no_data_timer", Src: []string{"undefined", "down", "going_up", "going_down", "up"}, Dst: "undefined"},
			},
			fsm.Callbacks{
				"enter_up": func(e *fsm.Event) {