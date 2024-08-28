package helpers

import (
	"GUI-GO/Apis"
	"GUI-GO/Server/Handlers"
	"context"
	"errors"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// CustomButton is a button with a custom background color.
type CustomButton struct {
	widget.Button
	backgroundColor color.Color
}

// NewCustomButton creates a new CustomButton with the specified label, tapped function, and background color.
func NewCustomButton(label string, tapped func(), backgroundColor color.Color) *CustomButton {
	button := &CustomButton{
		Button:          *widget.NewButton(label, tapped),
		backgroundColor: backgroundColor,
	}
	button.ExtendBaseWidget(button)
	return button
}

// CreateRenderer creates a custom renderer for the button.
func (b *CustomButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(b.backgroundColor)
	label := canvas.NewText(b.Text, theme.ForegroundColor())
	label.Alignment = fyne.TextAlignCenter

	return &customButtonRenderer{
		button:  b,
		bg:      bg,
		label:   label,
		objects: []fyne.CanvasObject{bg, label},
	}
}

// customButtonRenderer is a custom renderer for the CustomButton.
type customButtonRenderer struct {
	button  *CustomButton
	bg      *canvas.Rectangle
	label   *canvas.Text
	objects []fyne.CanvasObject
}

// Layout arranges the components of the button.
func (r *customButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.label.Resize(size)
}

// MinSize returns the minimum size of the button.
func (r *customButtonRenderer) MinSize() fyne.Size {
	return r.label.MinSize()
}

// Refresh refreshes the button to apply any updates.
func (r *customButtonRenderer) Refresh() {
	r.bg.FillColor = r.button.backgroundColor
	r.bg.Refresh()
	r.label.Text = r.button.Text
	r.label.Color = theme.ForegroundColor()
	r.label.Refresh()
}

// Objects returns the objects to be drawn.
func (r *customButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy cleans up the renderer.
func (r *customButtonRenderer) Destroy() {
	// No resources to clean up
}

// HandleDrawerSelection updates the right pane content based on the drawer item selection
func HandleDrawerSelection(client *mongo.Client, item, drawerItemText string, rightPane *fyne.Container) {

	serviceCollection := client.Database("test").Collection("services")
	rappCollection := client.Database("test").Collection("rapps")
	// subscriptionsCollection := client.Database("test").Collection("subscriptions")
	// subscribersCollection := client.Database("test").Collection("subscribers")
	dataTypeProdCapsCollection := client.Database("test").Collection("dataTypeProdCaps")
	dataJobsCollection := client.Database("test").Collection("dataJobs")

	// For testing purpose: insert a few documents into the rapps collection
	newRapps := []interface{}{
		Apis.Rapp{ApfId: "testrapp1", IsAuthorized: true, AuthorizedServices: []string{}},
		Apis.Rapp{ApfId: "testrapp2", IsAuthorized: false, AuthorizedServices: []string{}},
	}
	_, err := rappCollection.InsertMany(context.TODO(), newRapps)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	if item == "Discover Service API" && drawerItemText == "Get API" {
		// Create form fields
		apiInvokerID := widget.NewEntry()
		apiInvokerID.SetPlaceHolder("api-invoker-id")
		apiName := widget.NewEntry()
		apiName.SetPlaceHolder("api-name (optional)")
		apiVersion := widget.NewEntry()
		apiVersion.SetPlaceHolder("api-version (optional)")

		// Custom validation for api-invoker-id (required field)
		apiInvokerID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("api-invoker-id is required")
			}
			return nil
		}

		// Create labels
		labelInvokerID := widget.NewLabelWithStyle("api-invoker-id *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		labelName := widget.NewLabel("api-name")
		labelVersion := widget.NewLabel("api-version")

		// Description labels
		descriptionInvokerID := widget.NewLabel("String identifying the API invoker assigned by the CAPIF core function.")
		descriptionName := widget.NewLabel("API name, it is set as {apiName} part of the URI structure.")
		descriptionVersion := widget.NewLabel("API major version of the URI (e.g., v1).")

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := apiInvokerID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill apiInvokerID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				r.HandleFunc("/allServiceAPIs", Handlers.ServiceDiscoveryHandler(serviceCollection, rappCollection, &apiInvokerID.Text,
					&apiName.Text, &apiVersion.Text)).Methods("GET")
				fmt.Println("Button clicked!")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelInvokerID, descriptionInvokerID, apiInvokerID,
			labelName, descriptionName, apiName,
			labelVersion, descriptionVersion, apiVersion,
			buttonContainer,
		)
		highlightedLabel := canvas.NewText("API Name: Discover Service API - Get API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: This API is used to discover all the available services."),
			widget.NewLabel("Details: Calling this API will provide you with a list of available services."),
			form,
		}
	} else if item == "Register Service API" && drawerItemText == "Post API" {
		// Create form fields for Post API
		apfID := widget.NewEntry()
		apfID.SetPlaceHolder("apfID")

		// Custom validation for apiID (required field)
		apfID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("apfID is required")
			}
			return nil
		}

		// Create labels
		labelApfID := widget.NewLabelWithStyle("apfID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Rich text box with the provided JSON
		jsonData, err := loadJSONFromFile("JSONFiles/Register_Service_Post.json") // Update with the correct file path
		if err != nil {
			log.Printf("Error loading JSON: %v", err)
			return
		}

		richTextBox := widget.NewMultiLineEntry()
		richTextBox.SetText(jsonData)

		// Increase the height of the rich text box to double its current size
		richTextBoxContainer := container.NewVScroll(richTextBox)
		richTextBoxContainer.SetMinSize(fyne.NewSize(richTextBox.MinSize().Width, richTextBox.MinSize().Height*5))

		// Label the text box as "Request Body"
		labelRequestBody := widget.NewLabel("Request Body")

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := apfID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill apfID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+apfID.Text+"/service-apis", Handlers.PublishServiceHandler(serviceCollection, rappCollection)).Methods("POST")
				fmt.Println("/" + apfID.Text + "/service-apis")
				// You can add more submit logic here if needed
			}
		})

		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelApfID, apfID,
			labelRequestBody, richTextBoxContainer,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Register Service API - Post API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: Publish a new API."),
			form,
		}
	} else if item == "Register Service API" && drawerItemText == "Get API" {
		// Create form fields for Get API
		apfID := widget.NewEntry()
		apfID.SetPlaceHolder("apfID")

		// Custom validation for apfID (required field)
		apfID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("apfID is required")
			}
			return nil
		}

		richTextBox := widget.NewMultiLineEntry()

		// Increase the height of the rich text box to double its current size
		richTextBoxContainer := container.NewVScroll(richTextBox)
		richTextBoxContainer.SetMinSize(fyne.NewSize(richTextBox.MinSize().Width, richTextBox.MinSize().Height*5))

		// Create labels
		labelApfID := widget.NewLabelWithStyle("apfID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := apfID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill apfID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				// r.HandleFunc("/"+apfID.Text+"/service-apis/", Handlers.GetSpecificServiceAPIHandler(serviceCollection)).Methods("GET")
				r.HandleFunc("/"+apfID.Text+"/service-apis/", func(w http.ResponseWriter, r *http.Request) {
					// Get the JSON data by calling the modified handler function
					jsonData, err := Handlers.GetServiceAPIsHandler(serviceCollection, rappCollection)(r)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					// Optionally, you can process the jsonData before sending it in the response
					// For example, logging or modifying the JSON data

					// Set the response header and write the JSON data
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonData)
					richTextBox.SetText(string(jsonData))
				}).Methods("GET")

			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelApfID, apfID,
			buttonContainer,
			richTextBoxContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Register Service API - Get API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: Retrieve all published APIs."),
			form,
		}
	} else if item == "Register Service API" && drawerItemText == "Get Specific API" {
		// Create form fields for Get Specific API
		apfID := widget.NewEntry()
		apfID.SetPlaceHolder("apfID")
		serviceAPIID := widget.NewEntry()
		serviceAPIID.SetPlaceHolder("serviceAPIID")

		// Custom validation for apfID (required field)
		apfID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("apfID is required")
			}
			return nil
		}

		// Custom validation for serviceAPIID (required field)
		serviceAPIID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("serviceAPIID is required")
			}
			return nil
		}

		// Create labels
		labelApfID := widget.NewLabelWithStyle("apfID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		labelServiceAPIID := widget.NewLabelWithStyle("serviceAPIID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := apfID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill apfID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			}
			if err1 := serviceAPIID.Validate(); err1 != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Service API ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+apfID.Text+"/service-apis/"+serviceAPIID.Text, Handlers.GetSpecificServiceAPIHandler(serviceCollection)).Methods("GET")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelApfID, apfID,
			labelServiceAPIID, serviceAPIID,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Register Service API - Get Specific API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: Retrieve a published service API."),
			form,
		}
	} else if item == "Register Service API" && drawerItemText == "Delete API" {
		// Create form fields for Delete API
		apfID := widget.NewEntry()
		apfID.SetPlaceHolder("apfID")
		serviceAPIID := widget.NewEntry()
		serviceAPIID.SetPlaceHolder("serviceAPIID")

		// Custom validation for apfID (required field)
		apfID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("apfID is required")
			}
			return nil
		}

		// Custom validation for serviceAPIID (required field)
		serviceAPIID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("serviceAPIID is required")
			}
			return nil
		}

		// Create labels
		labelApfID := widget.NewLabelWithStyle("apfID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		labelServiceAPIID := widget.NewLabelWithStyle("serviceAPIID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := apfID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill apfID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			}
			if err1 := serviceAPIID.Validate(); err1 != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Service API ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+apfID.Text+"/service-apis/"+serviceAPIID.Text, Handlers.DeleteServiceAPIHandler(serviceCollection, rappCollection)).Methods("DELETE")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelApfID, apfID,
			labelServiceAPIID, serviceAPIID,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Register Service API - Delete API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: Unpublish a published service API."),
			form,
		}
	} else if item == "Register Service API" && drawerItemText == "Put API" {
		// Create form fields for Put API
		apfID := widget.NewEntry()
		apfID.SetPlaceHolder("apfID")
		serviceAPIID := widget.NewEntry()
		serviceAPIID.SetPlaceHolder("serviceAPIID")

		// Custom validation for apfID (required field)
		apfID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("apfID is required")
			}
			return nil
		}

		// Custom validation for serviceAPIID (required field)
		serviceAPIID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("serviceAPIID is required")
			}
			return nil
		}

		// Create labels
		labelApfID := widget.NewLabelWithStyle("apfID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		labelServiceAPIID := widget.NewLabelWithStyle("serviceAPIID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Rich text box with the provided JSON
		jsonText, err := loadJSONFromFile("JSONFiles/Register_Service_Post.json") // Update with the correct file path
		if err != nil {
			log.Printf("Error loading JSON: %v", err)
			return
		}

		richTextBox := widget.NewMultiLineEntry()
		richTextBox.SetText(jsonText)

		// Increase the height of the rich text box to double its current size
		richTextBoxContainer := container.NewVScroll(richTextBox)
		richTextBoxContainer.SetMinSize(fyne.NewSize(richTextBox.MinSize().Width, richTextBox.MinSize().Height*5))

		// Label the text box as "Request Body"
		labelRequestBody := widget.NewLabel("Request Body")

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := apfID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill apfID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			}
			if err1 := serviceAPIID.Validate(); err1 != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Service API ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+apfID.Text+"/service-apis/"+serviceAPIID.Text, Handlers.UpdateServiceAPIHandler(serviceCollection, rappCollection)).Methods("PUT")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelApfID, apfID,
			labelServiceAPIID, serviceAPIID,
			labelRequestBody, richTextBoxContainer,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Register Service API - Put API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: Update a published service API."),
			form,
		}
	} else if item == "Register Service API" && drawerItemText == "Patch API" {
		// Create form fields for Patch API
		apfID := widget.NewEntry()
		apfID.SetPlaceHolder("apfID")
		serviceAPIID := widget.NewEntry()
		serviceAPIID.SetPlaceHolder("serviceAPIID")

		// Custom validation for apfID (required field)
		apfID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("apfID is required")
			}
			return nil
		}

		// Custom validation for serviceAPIID (required field)
		serviceAPIID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("serviceAPIID is required")
			}
			return nil
		}

		// Create labels
		labelApfID := widget.NewLabelWithStyle("apfID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		labelServiceAPIID := widget.NewLabelWithStyle("serviceAPIID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Rich text box with the provided JSON
		jsonText, err := loadJSONFromFile("JSONFiles/Register_Service_Patch.json") // Update with the correct file path
		if err != nil {
			log.Printf("Error loading JSON: %v", err)
			return
		}

		richTextBox := widget.NewMultiLineEntry()
		richTextBox.SetText(jsonText)

		// Increase the height of the rich text box to double its current size
		richTextBoxContainer := container.NewVScroll(richTextBox)
		richTextBoxContainer.SetMinSize(fyne.NewSize(richTextBox.MinSize().Width, richTextBox.MinSize().Height*5))

		// Label the text box as "Request Body"
		labelRequestBody := widget.NewLabel("Request Body")

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := apfID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill apfID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			}
			if err1 := serviceAPIID.Validate(); err1 != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Service API ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+apfID.Text+"/service-apis/"+serviceAPIID.Text, Handlers.PatchServiceAPIHandler(serviceCollection, rappCollection)).Methods("PATCH")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelApfID, apfID,
			labelServiceAPIID, serviceAPIID,
			labelRequestBody, richTextBoxContainer,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Register Service API - Patch API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: Modify an existing published service API."),
			form,
		}
	} else if item == "Data Registration API" && drawerItemText == "Delete API" {
		// Create the RegistrationID text field
		registrationID := widget.NewEntry()
		registrationID.SetPlaceHolder("RegistrationID")

		// Custom validation for RegistrationID (required field)
		registrationID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("RegistrationID is required")
			}
			return nil
		}

		// Create labels
		labelRegistrationID := widget.NewLabelWithStyle("RegistrationID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Create the RegistrationID text field
		rAppID := widget.NewEntry()
		rAppID.SetPlaceHolder("rAppID")

		// Custom validation for RegistrationID (required field)
		rAppID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("rAppID is required")
			}
			return nil
		}

		// Create labels
		labelRAppID := widget.NewLabelWithStyle("rAppID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := registrationID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Registration ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			}
			if err1 := rAppID.Validate(); err1 != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill rApp ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/rapps/"+rAppID.Text+"/datatypeprodcaps/"+registrationID.Text, Handlers.DeregisterDmeTypeProdCapHandler(rappCollection, dataTypeProdCapsCollection)).Methods("DELETE")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelRegistrationID, registrationID,
			labelRAppID, rAppID,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Data Registration API - Delete", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: This API is used to delete a registered data entry."),
			form,
		}
	} else if item == "Data Registration API" && drawerItemText == "Post API" {
		// Create the rAppID text field
		rAppID := widget.NewEntry()
		rAppID.SetPlaceHolder("rAppID")

		// Custom validation for rAppID (required field)
		rAppID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("rAppID is required")
			}
			return nil
		}

		// JSON representation based on the provided structure
		jsonData, err := loadJSONFromFile("JSONFiles/Post_Data_Registration.json") // Update with the correct file path
		if err != nil {
			log.Printf("Error loading JSON: %v", err)
			return
		}

		// Rich text box for JSON
		jsonTextBox := widget.NewMultiLineEntry()
		jsonTextBox.SetText(jsonData)
		jsonTextBoxContainer := container.NewVScroll(jsonTextBox)
		jsonTextBoxContainer.SetMinSize(fyne.NewSize(500, 400))

		// Create labels
		labelrAppID := widget.NewLabelWithStyle("rAppID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := rAppID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill rApp ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/rapps/"+rAppID.Text+"/datatypeprodcaps", Handlers.RegisterDmeTypeProdCapHandler(rappCollection, dataTypeProdCapsCollection)).Methods("POST")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelrAppID, rAppID,
			widget.NewLabel("Editable JSON:"),
			jsonTextBoxContainer, buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Data Registration API - Register DME Type Prod CAP API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: This API is used to register a DME type production capability."),
			form,
		}
	} else if item == "Data Discovery API" && drawerItemText == "Get All API" {
		// Create the datatypes text field
		datatypes := widget.NewEntry()
		datatypes.SetPlaceHolder("datatypes")

		// Custom validation for datatypes (required field)
		datatypes.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("datatypes is required")
			}
			return nil
		}

		// Create labels
		labelDatatypes := widget.NewLabelWithStyle("datatypes *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := datatypes.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Data Types.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+datatypes.Text, Handlers.GetAllDataTypesHandler(dataTypeProdCapsCollection)).Methods("GET")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelDatatypes, datatypes,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Data Discovery API - Get All Data Types API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: This API retrieves all data types."),
			form,
		}
	} else if item == "Data Discovery API" && drawerItemText == "Get By ID API" {
		// Create the dataTypes and dataTypeId text fields
		dataTypes := widget.NewEntry()
		dataTypes.SetPlaceHolder("dataTypes")
		dataTypeId := widget.NewEntry()
		dataTypeId.SetPlaceHolder("dataTypeId")

		// Custom validation for dataTypes (required field)
		dataTypes.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("dataTypes is required")
			}
			return nil
		}

		// Custom validation for dataTypeId (required field)
		dataTypeId.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("dataTypeId is required")
			}
			return nil
		}

		// Create labels
		labelDataTypes := widget.NewLabelWithStyle("dataTypes *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		labelDataTypeId := widget.NewLabelWithStyle("dataTypeId *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := dataTypes.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Data Types.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			}
			if err := dataTypeId.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Data Types ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+dataTypes.Text+"/"+dataTypeId.Text, Handlers.GetDataTypeByIdHandler(dataTypeProdCapsCollection)).Methods("GET")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelDataTypes, dataTypes,
			labelDataTypeId, dataTypeId,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Data Discovery API - Get Data Type By ID API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: This API retrieves a specific data type by ID."),
			form,
		}
	} else if item == "Data Access API" && drawerItemText == "Post API" {
		// Create the consumerId text field
		consumerId := widget.NewEntry()
		consumerId.SetPlaceHolder("consumerId")

		// Custom validation for consumerId (required field)
		consumerId.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("consumerId is required")
			}
			return nil
		}

		// JSON representation based on the provided structure
		jsonData, err := loadJSONFromFile("JSONFiles/Data_Access_Post.json") // Update with the correct file path
		if err != nil {
			log.Printf("Error loading JSON: %v", err)
			return
		}

		// Rich text box for JSON
		jsonTextBox := widget.NewMultiLineEntry()
		jsonTextBox.SetText(jsonData)
		jsonTextBoxContainer := container.NewVScroll(jsonTextBox)
		jsonTextBoxContainer.SetMinSize(fyne.NewSize(500, 400))

		// Create labels
		labelConsumerID := widget.NewLabelWithStyle("consumerId *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := consumerId.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Consumer ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+consumerId.Text+"/dataJobs", Handlers.CreateDataJobHandler(dataJobsCollection)).Methods("POST")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelConsumerID, consumerId,
			widget.NewLabel("Editable JSON:"),
			jsonTextBoxContainer,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Data Access API - Post API For DataJob", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: This API is used to create a job described by the DataJobInfo."),
			form,
		}
	} else if item == "Data Access API" && drawerItemText == "Notify API" {
		// Create the dataAvailabilityNotificationURI text field
		dataAvailabilityNotificationURI := widget.NewEntry()
		dataAvailabilityNotificationURI.SetPlaceHolder("dataAvailabilityNotificationURI")

		// Custom validation for dataAvailabilityNotificationURI (required field)
		dataAvailabilityNotificationURI.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("dataAvailabilityNotificationURI is required")
			}
			return nil
		}

		// JSON representation based on the provided structure
		jsonData, err := loadJSONFromFile("JSONFiles/DataAvailabilityNotification.json") // Update with the correct file path
		if err != nil {
			log.Printf("Error loading JSON: %v", err)
			return
		}

		// Rich text box for JSON
		jsonTextBox := widget.NewMultiLineEntry()
		jsonTextBox.SetText(jsonData)
		jsonTextBoxContainer := container.NewVScroll(jsonTextBox)
		jsonTextBoxContainer.SetMinSize(fyne.NewSize(500, 400))

		// Create labels
		labeldataAvailabilityNotificationURI := widget.NewLabelWithStyle("dataAvailabilityNotificationURI *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := dataAvailabilityNotificationURI.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Data Availability Notification URI.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+dataAvailabilityNotificationURI.Text, Handlers.NotifyDataAvailabilityHandler(dataJobsCollection)).Methods("POST")
			}

		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labeldataAvailabilityNotificationURI, dataAvailabilityNotificationURI,
			widget.NewLabel("Editable JSON:"),
			jsonTextBoxContainer, buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Data Access API - Notify Data Availability(POST)", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: This API is used to notify consumer of the data availability."),
			form,
		}
	} else if item == "Data Access API" && drawerItemText == "Delete API" {
		// Create the mandatory text fields
		consumerID := widget.NewEntry()
		consumerID.SetPlaceHolder("consumerID")
		dataJobs := widget.NewEntry()
		dataJobs.SetPlaceHolder("dataJobs")
		dataJobId := widget.NewEntry()
		dataJobId.SetPlaceHolder("dataJobId")

		// Custom validation for each field (all are required)
		consumerID.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("consumerID is required")
			}
			return nil
		}

		dataJobs.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("dataJobs is required")
			}
			return nil
		}

		dataJobId.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("dataJobId is required")
			}
			return nil
		}

		// Create labels
		labelConsumerID := widget.NewLabelWithStyle("consumerID *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		labelDataJobs := widget.NewLabelWithStyle("dataJobs *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		labelDataJobId := widget.NewLabelWithStyle("dataJobId *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := consumerID.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Consumer ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			}
			if err1 := dataJobs.Validate(); err1 != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Data Jobs.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			}
			if err2 := dataJobId.Validate(); err2 != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Data Job ID.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/"+consumerID.Text+"/"+dataJobs.Text+"/"+dataJobId.Text, Handlers.DeleteDataJobHandler(dataJobsCollection)).Methods("DELETE")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labelConsumerID, consumerID,
			labelDataJobs, dataJobs,
			labelDataJobId, dataJobId,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Delete Data JOB API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: This API is used to delete a specific data job."),
			form,
		}
	} else if item == "Push Data API" && drawerItemText == "Post API" {
		// Create form fields for Post API
		dataPushURI := widget.NewEntry()
		dataPushURI.SetPlaceHolder("dataPushURI")

		// Custom validation for dataPushURI (required field)
		dataPushURI.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("dataPushURI is required")
			}
			return nil
		}

		// Create labels
		labeldataPushURI := widget.NewLabelWithStyle("dataPushURI *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Rich text box
		richTextBox := widget.NewMultiLineEntry()
		richTextBox.SetText("")

		// Increase the height of the rich text box to double its current size
		richTextBoxContainer := container.NewVScroll(richTextBox)
		richTextBoxContainer.SetMinSize(fyne.NewSize(richTextBox.MinSize().Width, richTextBox.MinSize().Height*5))

		// Label the text box as "Request Body"
		labelRequestBody := widget.NewLabel("Request Body")

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := dataPushURI.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Data Push URI.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/api/v1/push-data", Handlers.PushDataHandler()).Methods("POST")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labeldataPushURI, dataPushURI,
			labelRequestBody, richTextBoxContainer,
			buttonContainer,
		)
		highlightedLabel := canvas.NewText("API Name: Push Data API - Post API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel,
			widget.NewLabel("Description: Post the payload to specified URI."),
			form,
		}
	} else if item == "Pull Data API" && drawerItemText == "Get API" {
		// Create form fields for Post API
		dataPullURI := widget.NewEntry()
		dataPullURI.SetPlaceHolder("dataPullURI")

		// Custom validation for dataPullURI (required field)
		dataPullURI.Validator = func(s string) error {
			if len(s) == 0 {
				return errors.New("dataPullURI is required")
			}
			return nil
		}

		// Create labels
		labeldataPullURI := widget.NewLabelWithStyle("dataPullURI *", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Add a button at the bottom middle of the right pane
		submitButton := widget.NewButton("Submit", func() {
			if err := dataPullURI.Validate(); err != nil {
				// If the apiID field is empty, show a pop-up dialog using dialog.NewInformation
				dialog := dialog.NewInformation(
					"Validation Error",
					"Please fill Data Pull URI.",
					fyne.CurrentApp().Driver().AllWindows()[0],
				)
				dialog.Show()
			} else {
				fmt.Println("Button clicked!")
				r.HandleFunc("/api/v1/pull-data", Handlers.PullDataHandler()).Methods("GET")
			}
		})
		submitButton.Importance = widget.HighImportance
		buttonContainer := container.NewHBox(
			layout.NewSpacer(), // Spacer that pushes the button to the right
			submitButton,
		)

		// Create a form-like structure
		form := container.NewVBox(
			labeldataPullURI, dataPullURI,
			buttonContainer,
		)

		highlightedLabel := canvas.NewText("API Name: Pull Data API - Get API", theme.PrimaryColor())
		highlightedLabel.TextStyle = fyne.TextStyle{Bold: true}
		highlightedLabel.TextSize = theme.TextHeadingSize()
		highlightedLabel.Alignment = fyne.TextAlignLeading

		// Update the right pane with the form
		rightPane.Objects = []fyne.CanvasObject{
			highlightedLabel, widget.NewLabel("Description: Gets data from the specified."),
			form,
		}
	} else {
		rightPane.Objects = []fyne.CanvasObject{
			widget.NewLabel("You selected: " + drawerItemText),
			widget.NewLabel("Here is some random content for " + drawerItemText + "."),
		}
	}
	rightPane.Refresh()
}

// loadJSONFromFile loads JSON data from a file
func loadJSONFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open JSON file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("could not read JSON file: %w", err)
	}

	return string(bytes), nil
}
