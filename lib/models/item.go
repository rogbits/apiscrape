package models

type Item struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
	Sent     bool   `json:"sent"`
	Split    int    `json:"split"`
}

func NewItem() *Item {
	return new(Item)
}

func CopyItem(orig *Item) *Item {
	i := NewItem()
	i.Id = orig.Id
	i.Name = orig.Name
	i.Quantity = orig.Quantity
	i.Sent = orig.Sent
	i.Split = orig.Split

	return i
}

type ItemList struct {
	Items []*Item `json:"items"`
}

func (il *ItemList) AddItem(item *Item) {
	il.Items = append(il.Items, item)
}

func (il *ItemList) GetItem(i int64) *Item {
	return il.Items[i]
}
