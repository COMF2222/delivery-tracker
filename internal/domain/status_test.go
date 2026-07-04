package domain

import (
	"testing"
)

func TestCanChangeStatus(t *testing.T) {
	tests := []struct {
		name    string
		current Status
		next    Status
		want    bool
	}{
		{name: "created -> purchased", current: StatusCreated, next: StatusPurchased, want: true},
		{name: "created -> warehouse", current: StatusCreated, next: StatusWarehouse, want: false},
		{name: "purchased -> created", current: StatusPurchased, next: StatusCreated, want: false},
		{name: "delivered -> arrived", current: StatusDelivered, next: StatusArrived, want: false},
		{name: "arrived -> delivered", current: StatusArrived, next: StatusDelivered, want: true},
		{name: "created -> created", current: StatusCreated, next: StatusCreated, want: false},
		{name: "customs -> delivered", current: StatusCustoms, next: StatusDelivered, want: false},
		{name: "warehouse -> in_transit", current: StatusWarehouse, next: StatusInTransit, want: true},
		{name: "in_transit -> customs", current: StatusInTransit, next: StatusCustoms, want: true},
		{name: "delivered -> delivered", current: StatusDelivered, next: StatusDelivered, want: false},
		{name: "arrived -> customs", current: StatusArrived, next: StatusCustoms, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanChangeStatus(tt.current, tt.next)
			if result != tt.want {
				t.Errorf("Change status(%s, %s) = %v, want %v", tt.current, tt.next, result, tt.want)
			}
		})
	}
}
