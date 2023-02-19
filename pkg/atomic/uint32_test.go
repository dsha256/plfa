package atomic

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewValue(t *testing.T) {
	tests := []struct {
		name string
		want Value
	}{
		{
			name: "New value must be equal to 0",
			want: 0,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewValue(); got != tt.want {
				t.Errorf("NewValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_Decrement(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want uint32
	}{
		{
			name: "Concurrency safety on Increment",
			v:    0,
			want: 1000000,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			i := 0
			for i < int(tt.want) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					tt.v.Increment()
				}()
				i++
			}
			wg.Wait()
			require.Equal(t, tt.want, tt.v.Get())
		})
	}
}

func TestValue_Get(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want uint32
	}{
		{
			name: "Get",
			v:    0,
			want: 0,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.Increment()
			tt.v.Decrement()
			if got := tt.v.Get(); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_Increment(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want uint32
	}{
		{
			name: "Concurrency safety on Decrement",
			v:    0,
			want: 1000000,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			i := 0
			for i < int(tt.want)*2 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					tt.v.Increment()
				}()
				i++
			}

			i = 0
			for i < int(tt.want) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					tt.v.Decrement()
				}()
				i++
			}
			wg.Wait()

			require.Equal(t, tt.want, tt.v.Get())
		})
	}
}
