package handlers

var Default = map[any]Handler{
	"/start":            Start,
	&createGroupButton:  Test,
	&manageGroupsButton: Test,
	&managePlantsButton: Test,
	&addFlowerButton:    Test,
}
