package domain

type Status string

const (
	StatusCreated   Status = "CREATED"
	StatusPurchased Status = "PURCHASED"
	StatusWarehouse Status = "WAREHOUSE"
	StatusInTransit Status = "IN_TRANSIT"
	StatusCustoms   Status = "CUSTOMS"
	StatusArrived   Status = "ARRIVED"
	StatusDelivered Status = "DELIVERED"
)

var allowedTransitions = map[Status]Status{
	StatusCreated:   StatusPurchased,
	StatusPurchased: StatusWarehouse,
	StatusWarehouse: StatusInTransit,
	StatusInTransit: StatusCustoms,
	StatusCustoms:   StatusArrived,
	StatusArrived:   StatusDelivered,
}

func CanChangeStatus(current, next Status) bool {

	allowedNext := allowedTransitions[current]

	return allowedNext == next
}
