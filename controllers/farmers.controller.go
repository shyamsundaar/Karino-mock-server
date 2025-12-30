package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/shyamsundaar/karino-mock-server/initializers"
	"github.com/shyamsundaar/karino-mock-server/models"
	// "github.com/shyamsundaar/karino-mock-server/query"
	// "gorm.io/gorm"
)

// CreateCustomerDetailHandler handles POST /spic_to_erp/customers/:coopId/farmers
// @Summary      Create a new farmer detail
// @Description  Create a new record in the details table
// @Tags         Details
// @Accept       json
// @Produce      json
// @Param        coopId  path      string                            true  "Cooperative ID"
// @Param        detail  body      models.CreateDetailSchema          true  "Create Detail Payload"
// @Success      201     {object}  models.CreateSuccessFarmerResponse
// @Router       /spic_to_erp/customers/{coopId}/farmers [post]
func CreateCustomerDetailHandler(c *fiber.Ctx) error {
	// 1. Get CoopID from URL Parameter
	coopId := c.Params("coopId")
	if coopId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "CoopID is required in URL"})
	}

	var payload *models.CreateDetailSchema

	// 2. Parse the JSON Body
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// 3. Validate the Body
	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// 4. Map everything to the DB Model
	newDetail := models.FarmerDetails{
		CoopID:                      coopId, // Set from URL Param
		FarmerID:                    payload.FarmerID,
		FirstName:                   payload.FirstName,
		LastName:                    payload.LastName,
		MobileNumber:                payload.MobileNumber,
		RegionID:                    payload.RegionID,
		RegionPartID:                payload.RegionPartID,
		SettlementID:                payload.SettlementID,
		SettlementPartID:            payload.SettlementPartID,
		CustomGeographyStructure1ID: payload.CustomGeo1ID,
		CustomGeographyStructure2ID: payload.CustomGeo2ID,
		ZipCode:                     payload.ZipCode,
		FarmerKycTypeID:             payload.FarmerKycTypeID,
		FarmerKycType:               payload.FarmerKycType,
		FarmerKycID:                 payload.FarmerKycID,
		ClubID:                      payload.ClubID,
		ClubName:                    payload.ClubName,
		ClubLeaderFarmerID:          payload.ClubLeaderFarmerID,
	}

	// 5. Save to Database (GORM fills in CreatedAt/UpdatedAt here)
	result := initializers.DB.Create(&newDetail)
	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	// Inside CreateCustomerDetailHandler...

	response := models.CreateSuccessFarmerResponse{
		Success: true,
		Data: models.FarmerResponse{
			TempERPCustomerID: newDetail.TempID,
			ErpCustomerId:     newDetail.CustomerID,
			ErpVendorId:       newDetail.VendorID,
			FarmerId:          newDetail.FarmerID,
			CreatedAt:         newDetail.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:         newDetail.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Message:           "Farmer detail created successfully",
		},
	}

	return c.Status(fiber.StatusCreated).JSON(response)

}

// FindDetails handles GET /spic_to_erp/customers/:coopId/farmers
// @Summary      List farmer details
// @Description  Get a paginated list of farmer details for a specific cooperative
// @Tags         Details
// @Accept       json
// @Produce      json
// @Param        coopId path      string  true   "Cooperative ID"
// @Param        updatedFrom   query     string  false  " "
// @Param        updatedTo     query     string  false  " "
// @Param        page          query     int     false  "Page number"    default(1)
// @Param        limit         query     int     false  "Items per page" default(10)
// @Success      200    {object}  models.SuccessListResponse
// @Router       /spic_to_erp/customers/{coopId}/farmers [get]
func FindDetailsHandler(c *fiber.Ctx) error {
	coopId := c.Params("coopId")

	// 1. Get query params
	updatedFrom := c.Query("updatedFrom")
	updatedTo := c.Query("updatedTo")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	offset := (page - 1) * limit
	var farmers []models.FarmerDetails
	var totalRecords int64

	// 2. Start building the query
	query := initializers.DB.Model(&models.FarmerDetails{}).Where("coop_id = ?", coopId)

	// 3. Apply Date Filters ONLY if they are provided
	if updatedFrom != "" {
		// Use >= for the start date
		query = query.Where("updated_at >= ?", updatedFrom)
	}
	if updatedTo != "" {
		// Use <= for the end date
		query = query.Where("updated_at <= ?", updatedTo)
	}

	// 4. Execute Count and Find
	query.Count(&totalRecords)
	result := query.Limit(limit).Offset(offset).Find(&farmers)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})

	}
	// 4. Calculate total pages
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	// 5. Return nested response
	return c.Status(fiber.StatusOK).JSON(models.SuccessListResponse{
		Success: true,
		Data: models.ListFarmerResponse{
			Items: farmers,
			Metadata: models.PaginationMetadata{
				TotalRecord: int(totalRecords),
				TotalPage:   totalPages,
				CurrentPage: page,
				Limit:       limit,
			},
		},
	})
}

// // FindDetailById handles GET /api/details/:detailId
// // @Summary      Get detail by ID
// // @Description  Retrieve a single farmer detail record
// // @Tags         Details
// // @Produce      json
// // @Param        detailId  path      string  true  "Detail ID"
// // @Success      200       {object}  models.Detail
// // @Failure      502      {object}  map[string]interface{}
// // @Router       /details/{detailId} [get]
// func FindDetailById(c *fiber.Ctx) error {
// 	detailId := c.Params("detailId")

// 	var detail models.FarmerDetails
// 	result := initializers.DB.First(&detail, "id = ?", detailId)
// 	if err := result.Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No detail with that ID exists"})
// 		}
// 		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"detail": detail}})
// }

// // DeleteDetail handles DELETE /api/details/:detailId
// // @Summary      Delete a detail
// // @Description  Remove a farmer detail record from the database
// // @Tags         Details
// // @Param        detailId  path      string  true  "Detail ID"
// // @Success      204       "No Content"
// // @Failure      502      {object}  map[string]interface{}
// // @Router       /details/{detailId} [delete]
// func DeleteDetail(c *fiber.Ctx) error {
// 	detailId := c.Params("detailId")

// 	result := initializers.DB.Delete(&models.FarmerDetails{}, "id = ?", detailId)

// 	if result.RowsAffected == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No detail with that ID exists"})
// 	} else if result.Error != nil {
// 		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
// 	}

// 	return c.SendStatus(fiber.StatusNoContent)
// }
