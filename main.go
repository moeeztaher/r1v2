package main

import (
	"GUI-GO/ui"
	"context"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	a := app.New()
	win := a.NewWindow("R1 UI Tool")
	win.Resize(fyne.NewSize(1200, 800))

	// Initialize the left and right panes
	leftPane, rightPane := ui.SetupUI(client)

	// Set a fixed width for the left pane with padding
	paddedLeftPane := container.NewPadded(leftPane)
	fixedWidthLeftPane := container.NewVBox(paddedLeftPane)
	fixedWidthLeftPane.Resize(fyne.NewSize(300, 800)) // Setting the width to 300 pixels

	// Create a spacer for padding on the left side of the right pane
	leftPadding := canvas.NewRectangle(nil)
	leftPadding.SetMinSize(fyne.NewSize(10, 0)) // 40px padding on the left side

	// Right pane with padding, covering the remaining space
	paddedRightPane := container.NewPadded(rightPane)
	expandableRightPane := container.NewBorder(nil, nil, leftPadding, nil, paddedRightPane)

	// Create the main content area using Border layout
	content := container.NewBorder(nil, nil, fixedWidthLeftPane, nil, expandableRightPane)

	// Create the overall layout
	layout := container.NewBorder(nil, nil, nil, nil, content)

	// Set the window content
	win.SetContent(layout)

	// Show and run the application
	win.ShowAndRun()

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")
}
