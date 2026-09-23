package booking

import (
	"context"
	"testing"
	"time"
)

func TestCreateBooking_Validation(t *testing.T) {
	service := NewService(nil)
	ctx := context.Background()

	now := time.Now()
	pastTime := now.Add(-24 * time.Hour).Format(time.RFC3339)
	futureTime1 := now.Add(24 * time.Hour).Format(time.RFC3339)
	futureTime2 := now.Add(26 * time.Hour).Format(time.RFC3339)

	tests := []struct {
		name          string
		startTimeStr  string
		endTimeStr    string
		expectedError string
	}{
		{
			name:          "Ошибка: неверный формат времени",
			startTimeStr:  "просто текст",
			endTimeStr:    futureTime2,
			expectedError: "invalid start_time format",
		},
		{
			name:          "Ошибка: бронирование в прошлом",
			startTimeStr:  pastTime,
			endTimeStr:    futureTime1,
			expectedError: "cannot book in the past",
		},
		{
			name:          "Ошибка: конец раньше начала",
			startTimeStr:  futureTime2,
			endTimeStr:    futureTime1,
			expectedError: "end_time must be after start_time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.CreateBooking(ctx, "user1", "room1", tt.startTimeStr, tt.endTimeStr)

			if err == nil {
				t.Errorf("Ожидали ошибку '%s', но ошибки не было", tt.expectedError)
				return
			}

			if err.Error() != tt.expectedError {
				t.Errorf("Ожидали ошибку '%s', а получили '%s'", tt.expectedError, err.Error())
			}
		})
	}
}

