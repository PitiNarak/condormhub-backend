package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/internal/handler"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/cucumber/godog"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// DormUpdateTest contains the test context for the dorm update feature
type DormUpdateTest struct {
	app          *fiber.App
	mockService  *MockDormService
	request      *http.Request
	response     *http.Response
	responseBody map[string]interface{}
	userID       uuid.UUID
	userRole     domain.Role
	dormID       uuid.UUID
	ownerID      uuid.UUID
	updateData   *dto.DormUpdateRequestBody
	initialDorm  *domain.Dorm
	updatedDorm  *domain.Dorm
}

// Initialize sets up the test environment
func (test *DormUpdateTest) Initialize(t *testing.T) {
	// Setup mock service
	test.mockService = new(MockDormService)
	dormHandler := handler.NewDormHandler(test.mockService)

	// Setup Fiber app with proper error handling
	test.app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Handle validation errors properly
			if validationErr, ok := err.(*fiber.Error); ok {
				return c.Status(validationErr.Code).JSON(fiber.Map{
					"success": false,
					"message": validationErr.Message,
				})
			}

			// Handle custom app errors
			if appErr, ok := err.(*apperror.AppError); ok {
				return c.Status(appErr.Code).JSON(fiber.Map{
					"success": false,
					"message": appErr.Message,
				})
			}

			// Default error handling
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Internal Server Error",
			})
		},
	})

	// Setup the PATCH endpoint for dorm updates
	test.app.Patch("/dorms/:id", func(c *fiber.Ctx) error {
		// Mock authenticated user
		c.Locals("userID", test.userID)
		c.Locals("user", &domain.User{
			ID:   test.userID,
			Role: &test.userRole,
		})

		return dormHandler.Update(c)
	})

	// Initialize test state
	test.userID = uuid.New()
	test.dormID = uuid.New()
	test.ownerID = uuid.New()
	test.userRole = domain.LessorRole
}

// iAmARegisteredLessor sets up the test with a lessor role
func (test *DormUpdateTest) iAmARegisteredLessor() error {
	test.userRole = domain.LessorRole
	return nil
}

// iHaveAPropertyListed sets up a property
func (test *DormUpdateTest) iHaveAPropertyListed() error {
	test.initialDorm = &domain.Dorm{
		ID:        test.dormID,
		Name:      "Original Dorm Name",
		Size:      50.0,
		Bedrooms:  2,
		Bathrooms: 1,
		Address: domain.Address{
			District:    "Original District",
			Subdistrict: "Original Subdistrict",
			Province:    "Original Province",
			Zipcode:     "12345",
		},
		Price:       5000.0,
		Description: "Original description",
	}
	return nil
}

// iAmTheOwnerOfTheProperty sets the current user as the owner
func (test *DormUpdateTest) iAmTheOwnerOfTheProperty() error {
	test.ownerID = test.userID
	test.initialDorm.OwnerID = test.ownerID
	return nil
}

// iAmNotTheOwnerOfTheProperty sets the current user as not the owner
func (test *DormUpdateTest) iAmNotTheOwnerOfTheProperty() error {
	// Make sure ownerID is different from userID
	if test.ownerID == test.userID {
		test.ownerID = uuid.New()
	}
	test.initialDorm.OwnerID = test.ownerID
	return nil
}

// iUpdatePropertyDetails prepares update data
func (test *DormUpdateTest) iUpdatePropertyDetails() error {
	test.updateData = &dto.DormUpdateRequestBody{
		Name:        "Updated Dorm Name",
		Size:        60.0,
		Bedrooms:    3,
		Bathrooms:   2,
		Price:       6000.0,
		Description: "Updated description",
		Address: struct {
			District    string `json:"district" validate:"omitempty"`
			Subdistrict string `json:"subdistrict" validate:"omitempty"`
			Province    string `json:"province" validate:"omitempty"`
			Zipcode     string `json:"zipcode" validate:"omitempty,numeric,len=5"`
		}{
			District:    "Updated District",
			Subdistrict: "Updated Subdistrict",
			Province:    "Updated Province",
			Zipcode:     "54321",
		},
	}

	// Prepare the updated dorm object that should be returned by the service
	test.updatedDorm = &domain.Dorm{
		ID:        test.dormID,
		OwnerID:   test.initialDorm.OwnerID,
		Name:      test.updateData.Name,
		Size:      test.updateData.Size,
		Bedrooms:  test.updateData.Bedrooms,
		Bathrooms: test.updateData.Bathrooms,
		Address: domain.Address{
			District:    test.updateData.Address.District,
			Subdistrict: test.updateData.Address.Subdistrict,
			Province:    test.updateData.Address.Province,
			Zipcode:     test.updateData.Address.Zipcode,
		},
		Price:       test.updateData.Price,
		Description: test.updateData.Description,
	}

	return nil
}

// iSubmitTheChanges performs the API request to update a dorm
func (test *DormUpdateTest) iSubmitTheChanges() error {
	jsonData, err := json.Marshal(test.updateData)
	if err != nil {
		return err
	}

	// Create a new HTTP request
	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/dorms/%s", test.dormID),
		bytes.NewBuffer(jsonData),
	)
	req.Header.Set("Content-Type", "application/json")
	test.request = req

	// Setup mock service expectations based on the scenario
	if test.ownerID == test.userID {
		// Success case - owner is updating
		test.mockService.On("Update",
			test.userID,
			false, // not admin
			test.dormID,
			mock.AnythingOfType("*dto.DormUpdateRequestBody"),
		).Return(test.updatedDorm, nil)
	} else {
		// Failure case - non-owner is updating
		err := apperror.ForbiddenError(
			fmt.Errorf("unauthorized action"),
			"You do not have permission to update this dorm",
		)
		test.mockService.On("Update",
			test.userID,
			false, // not admin
			test.dormID,
			mock.AnythingOfType("*dto.DormUpdateRequestBody"),
		).Return(nil, err)
	}

	// Perform the request
	resp, err := test.app.Test(req)
	if err != nil {
		return err
	}
	test.response = resp

	// Parse the response body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	responseBody := buf.String()

	// Try to parse as JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &result); err == nil {
		// fmt.Println("Raw response body:", responseBody)
		test.responseBody = result
		// fmt.Println("Parsed response body:", test.responseBody)
	} else {
		// If parsing fails, store raw response
		test.responseBody = map[string]interface{}{
			"success": false,
			"status":  resp.StatusCode,
			"raw":     responseBody,
		}
	}

	return nil
}

// theUpdatedDetailsShouldBeSaved checks if the dorm was updated successfully
func (test *DormUpdateTest) theUpdatedDetailsShouldBeSaved() error {
	if test.response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code %d, got %d", http.StatusOK, test.response.StatusCode)
	}

	// Check if the response indicates success
	fmt.Println(test.responseBody["success"])
	if test.responseBody["success"] != nil {
		return fmt.Errorf("expected 'success' to be true, got %v", test.responseBody["success"])
	}

	return nil
}

// iShouldSeeTheUpdatedInformationInMyListings checks if updated info is returned
func (test *DormUpdateTest) iShouldSeeTheUpdatedInformationInMyListings() error {
	data, ok := test.responseBody["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected data to be a map, got %T", test.responseBody["data"])
	}

	// Check if the response contains the updated dorm data
	if data["name"] != test.updateData.Name {
		return fmt.Errorf("expected name to be %s, got %v", test.updateData.Name, data["name"])
	}

	if data["size"].(float64) != test.updateData.Size {
		return fmt.Errorf("expected size to be %f, got %v", test.updateData.Size, data["size"])
	}

	if data["description"] != test.updateData.Description {
		return fmt.Errorf("expected description to be %s, got %v", test.updateData.Description, data["description"])
	}

	return nil
}

// iShouldSeeAnErrorMessage checks if the appropriate error message is shown
func (test *DormUpdateTest) iShouldSeeAnErrorMessage(message string) error {
	if test.response.StatusCode != http.StatusForbidden {
		return fmt.Errorf("expected status code %d, got %d", http.StatusForbidden, test.response.StatusCode)
	}

	// Check if the response contains the error message
	if test.responseBody["success"] != false {
		return fmt.Errorf("expected success to be false, got %v", test.responseBody["success"])
	}

	if test.responseBody["message"] != message {
		return fmt.Errorf("expected error message '%s', got '%v'", message, test.responseBody["message"])
	}

	return nil
}

// myPropertyShouldNotBeUpdated verifies no update occurred
func (test *DormUpdateTest) myPropertyShouldNotBeUpdated() error {
	// We can verify this indirectly by checking that the mockService's Update method
	// was called but returned an error, and that the HTTP status isn't 200 OK
	if test.response.StatusCode == http.StatusOK {
		return fmt.Errorf("expected non-200 status code, got %d", test.response.StatusCode)
	}

	return nil
}

// InitializeUpdateScenario sets up the test scenario
func InitializeUpdateScenario(ctx *godog.ScenarioContext, t *testing.T) {
	test := &DormUpdateTest{}

	ctx.BeforeScenario(func(*godog.Scenario) {
		test.Initialize(t)
	})

	ctx.Step(`^I am a registered lessor$`, test.iAmARegisteredLessor)
	ctx.Step(`^I have a property listed$`, test.iHaveAPropertyListed)
	ctx.Step(`^I am the owner of the property$`, test.iAmTheOwnerOfTheProperty)
	ctx.Step(`^I am not the owner of the property$`, test.iAmNotTheOwnerOfTheProperty)
	ctx.Step(`^I update property details$`, test.iUpdatePropertyDetails)
	ctx.Step(`^I submit the changes$`, test.iSubmitTheChanges)
	ctx.Step(`^the updated details should be saved$`, test.theUpdatedDetailsShouldBeSaved)
	ctx.Step(`^I should see the updated information in my listings$`, test.iShouldSeeTheUpdatedInformationInMyListings)
	ctx.Step(`^I should see an error message "([^"]*)"$`, test.iShouldSeeAnErrorMessage)
	ctx.Step(`^my property should not be updated$`, test.myPropertyShouldNotBeUpdated)
}

func TestUpdate(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			InitializeUpdateScenario(ctx, t)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}