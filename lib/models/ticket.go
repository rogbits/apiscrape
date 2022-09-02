package models

type Ticket struct {
	Id           int64  `json:"id"`
	TicketNumber int64  `json:"ticket_number"`
	OpenedAt     int64  `json:"opened_at"`
	ClosedAt     int64  `json:"closed_at"`
	FireDate     int64  `json:"fire_date"`
	FireTime     int64  `json:"fire_time"`
	ReadyDate    int64  `json:"ready_date"`
	ReadyTime    int64  `json:"ready_time"`
	Name         string `json:"name"`
	Open         bool   `json:"open"`
	Void         bool   `json:"void"`

	Embedded TicketEmbed `json:"_embedded"`
}

type TicketEmbed struct {
	Items []*Item `json:"items"`
}

func NewTicket() *Ticket {
	return new(Ticket)
}

type TicketList struct {
	Tickets []*Ticket `json:"tickets"`
}

func NewTicketList() *TicketList {
	return new(TicketList)
}

func (tl *TicketList) AddTicket(ticket *Ticket) {
	tl.Tickets = append(tl.Tickets, ticket)
}

func (tl *TicketList) GetTicket(i uint64) *Ticket {
	return tl.Tickets[i]
}
