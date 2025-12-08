package service

import (
	"log"
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllAchievementReferences(c *fiber.Ctx) error {
	refs, err := repository.GetAllAchievementReferences(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data referensi", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": refs})
}

func GetAchievementReferenceByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Referensi tidak ditemukan"})
	}

	detail, err := repository.GetAchievementMongoByID(database.MongoDB, ref.MongoAchievementID)

	response := fiber.Map{
		"reference": ref,
		"detail":    detail,
	}

	if err != nil {
		response["warning"] = "Detail prestasi tidak ditemukan di MongoDB"
	}

	return c.JSON(fiber.Map{"success": true, "data": response})
}

// CreateAchievementReference
func CreateAchievementReference(c *fiber.Ctx) error {
	var req model.SubmitAchievementRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Request body tidak valid",
		})
	}

	if req.StudentID == uuid.Nil || req.Title == "" || req.AchievementType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Field student_id, title, dan achievement_type wajib diisi",
		})
	}

	mongoData := model.Achievement{
		StudentID:       req.StudentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
	}

	mongoID, err := repository.InsertAchievementMongo(database.MongoDB, mongoData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menyimpan detail ke MongoDB", "error": err.Error(),
		})
	}

	now := time.Now()
	ref := &model.AchievementReference{
		StudentID:          req.StudentID,
		MongoAchievementID: mongoID,
		Status:             model.StatusSubmitted, // FIX UTAMA
		SubmittedAt:        &now,
	}

	if err := repository.CreateAchievementReference(database.DB, ref); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menyimpan referensi ke PostgreSQL", "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Prestasi berhasil disubmit",
		"data": fiber.Map{
			"reference_id": ref.ID,
			"mongo_id":     mongoID,
		},
	})
}

// Update Status
func UpdateAchievementReference(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID referensi tidak valid"})
	}

	var req model.UpdateAchievementStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	if req.Status != model.StatusVerified && req.Status != model.StatusRejected {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Status harus 'verified' atau 'rejected'",
		})
	}

	if req.Status == model.StatusRejected && (req.RejectionNote == nil || *req.RejectionNote == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Catatan penolakan wajib diisi jika rejected",
		})
	}

	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Referensi tidak ditemukan"})
	}

	now := time.Now()
	ref.Status = req.Status
	ref.VerifiedBy = &req.VerifiedBy
	ref.VerifiedAt = &now

	if req.Status == model.StatusRejected {
		ref.RejectionNote = req.RejectionNote
	}

	if err := repository.UpdateAchievementReference(database.DB, ref); err != nil {
		log.Println("Error update achievement ref:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengupdate referensi", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Status prestasi berhasil diupdate", "data": ref})
}

func DeleteAchievementReference(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	if err := repository.DeleteAchievementReference(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghapus referensi", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Referensi berhasil dihapus"})
}
