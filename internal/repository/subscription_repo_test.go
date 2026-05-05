package repository

import (
	"testing"
	"time"
)

func TestMonthsBetween(t *testing.T) {
	tests := []struct {
		name     string
		a        time.Time
		b        time.Time
		expected int
	}{
		{
			name:     "same month",
			a:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			b:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "one month difference",
			a:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			b:        time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "twelve months",
			a:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			b:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 12,
		},
		{
			name:     "reverse order",
			a:        time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			b:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := monthsBetween(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("monthsBetween(%v, %v) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMonthsIntersection(t *testing.T) {
	jan2025 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	feb2025 := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	mar2025 := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	apr2025 := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	dec2025 := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		start1   time.Time
		end1     time.Time
		start2   time.Time
		end2     *time.Time
		expected int
	}{
		{
			name:     "exact match",
			start1:   jan2025,
			end1:     mar2025,
			start2:   jan2025,
			end2:     &mar2025,
			expected: 3, // январь, февраль, март
		},
		{
			name:     "partial overlap",
			start1:   jan2025,
			end1:     mar2025,
			start2:   feb2025,
			end2:     &apr2025,
			expected: 2, // февраль, март
		},
		{
			name:     "no overlap",
			start1:   jan2025,
			end1:     feb2025,
			start2:   mar2025,
			end2:     &apr2025,
			expected: 0,
		},
		{
			name:     "subscription infinite",
			start1:   jan2025,
			end1:     mar2025,
			start2:   feb2025,
			end2:     nil,
			expected: 2, // февраль, март
		},
		{
			name:     "request infinite",
			start1:   jan2025,
			end1:     mar2025,
			start2:   feb2025,
			end2:     &dec2025,
			expected: 2, // февраль, март (ограничено end1)
		},
		{
			name:     "single month",
			start1:   jan2025,
			end1:     jan2025,
			start2:   jan2025,
			end2:     &jan2025,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := monthsIntersection(tt.start1, tt.end1, tt.start2, tt.end2)
			if result != tt.expected {
				t.Errorf("monthsIntersection(%v, %v, %v, %v) = %d, expected %d",
					tt.start1, tt.end1, tt.start2, tt.end2, result, tt.expected)
			}
		})
	}
}