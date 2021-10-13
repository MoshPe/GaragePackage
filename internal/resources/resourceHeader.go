package resources

import "container/list"

type Resource struct {
	Name            string
	AmountAvailable int
	whenAvailable   []list.List
}
