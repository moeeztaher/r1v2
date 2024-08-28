package ui

import (
	helpers "GUI-GO/ui-helpers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.mongodb.org/mongo-driver/mongo"
)

var leftPaneItems = []string{
	"Discover Service API", "Register Service API", "Data Access API", "Data Registration API",
	"Data Discovery API", "Pull Data API", "Push Data API", "Subscription API",
}

var drawerItems = [][]string{
	{"Get API"},
	{"Post API", "Get Specific API", "Get API", "Put API", "Delete API", "Patch API"},
	{"Post API", "Delete API", "Notify API"},
	{"Post API", "Delete API"},
	{"Get All API", "Get By ID API"},
	{"Get API"},
	{"Post API"},
	{"Entry Option 1", "Entry Option 2", "Entry Option 3"},
}

// SetupUI initializes the UI components
func SetupUI(client *mongo.Client) (fyne.CanvasObject, fyne.CanvasObject) {
	// Container to hold the items in the left pane including the drawers
	leftPane := container.NewVBox()

	// Container to hold the right-side content
	rightPane := container.NewVBox(widget.NewLabel("Select an option from the drawer"))

	// Track which item is currently expanded
	var expandedItemIndex = -1

	// Populate the left pane with buttons and drawers
	for i, item := range leftPaneItems {
		index := i // capture the index in the loop

		// Define the button outside the callback function to allow capturing it in the closure
		var button *widget.Button

		button = widget.NewButtonWithIcon(item, theme.Icon(theme.IconNameArrowDropDown), func() {
			if expandedItemIndex == index {
				// If this item is already expanded, collapse it
				leftPane.Objects[index*2+1] = container.NewWithoutLayout() // remove the drawer content
				button.SetIcon(theme.Icon(theme.IconNameArrowDropDown))
				expandedItemIndex = -1
			} else {
				// Expand this item
				if expandedItemIndex != -1 {
					// If another item is expanded, collapse it first
					leftPane.Objects[expandedItemIndex*2+1] = container.NewWithoutLayout()
					// Reset the icon of the previously expanded button
					prevButton, ok := leftPane.Objects[expandedItemIndex*2].(*widget.Button)
					if ok {
						prevButton.SetIcon(theme.Icon(theme.IconNameArrowDropDown))
					}
				}

				// Create the drawer content
				drawer := container.NewVBox()
				for _, drawerItem := range drawerItems[index] {
					drawerItemText := drawerItem // capture drawer item text

					// Set the color of drawerItemButton to 50% less bright compared to leftPaneButton
					drawerItemButton := helpers.NewCustomButton(drawerItemText, func() {
						helpers.HandleDrawerSelection(client, item, drawerItemText, rightPane)
					}, theme.PrimaryColor())

					// Wrap drawerItemButton in a container to set its width smaller than the left pane items
					drawer.Add(container.NewPadded(container.NewMax(drawerItemButton)))
				}

				// Add the drawer under the clicked item
				leftPane.Objects[index*2+1] = drawer
				expandedItemIndex = index
				button.SetIcon(theme.Icon(theme.IconNameArrowDropUp))
			}

			// Refresh the container to show the changes
			leftPane.Refresh()
		})

		// Add the button and an empty container (placeholder for drawer) to the left pane
		leftPane.Add(button)
		leftPane.Add(container.NewWithoutLayout()) // no layout to avoid extra spacing
	}

	return leftPane, rightPane
}
