package service

import (
	"context"
	"testing"
)

func TestService_Process(t *testing.T) {

	type args struct {
		ctx   context.Context
		batch Batch
		tasks chan *Item
	}

	tests := []struct {
		name    string
		args    args
		prepare func(c context.CancelFunc)
		wantErr bool
	}{
		{
			name: "SuccessfulProcessing",
			args: args{
				ctx:   context.Background(),
				batch: []*Item{nil, nil, nil, nil},
				tasks: make(chan *Item, 4),
			},
		},
		{
			name: "ErrorCancelContext",
			args: args{
				ctx:   context.Background(),
				batch: []*Item{nil, nil, nil, nil},
				tasks: make(chan *Item, 4),
			},
			prepare: func(c context.CancelFunc) {
				c()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &Service{
				Count:       10,
				IntervalSec: 5,
			}

			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()

			if tt.prepare != nil {
				tt.prepare(cancel)
			}

			err := s.Process(ctx, &tt.args.batch, tt.args.tasks)

			if (err != nil) != tt.wantErr {
				t.Errorf("SuccessfulProcessing() error: %v, wantErr: %v", err, tt.wantErr)
			}

		})
	}

}
