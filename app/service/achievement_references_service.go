package service

import (
	"fmt"
	"path/filepath"
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GET /api/v1/achievements (Filtered by Role)
func ListAchievements(c *fiber.Ctx) error {
	// Ambil UserID dan RoleID dari Middleware Locals
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	
	// Cek apakah user adalah Student, Dosen, atau Admin
	// Cara terbaik: Cek di tabel Users / Lecturers / Students
	// Disini kita asumsikan check sederhana via repository helpers
	
	var refs []model.AchievementReference
	var err error

	// 1. Cek apakah dia Student?
	student, errStudent := repository.GetStudentByID(database.DB, userID)
	if errStudent == nil {
		// Dia Student: Tampilkan miliknya saja
		refs, err = repository.GetAchievementReferencesByStudentID(database.DB, student.ID)
	} else {
		// 2. Cek apakah dia Lecturer?
		lecturer, errLecturer := repository.GetLecturerByID(database.DB, userID)
		if errLecturer == nil {
			// Dia Lecturer: Tampilkan milik mahasiswa bimbingannya
			refs, err = repository.GetAchievementReferencesByAdvisorID(database.DB, lecturer.ID)
		} else {
			// 3. Asumsi Admin: Tampilkan Semua
			refs, err = repository.GetAllAchievementReferences(database.DB)
		}
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data"})
	}
	return c.JSON(fiber.Map{"success": true, "data": refs})
}

// GET /api/v1/achievements/:id
func GetAchievementDetail(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	detail, _ := repository.GetAchievementMongoByID(database.MongoDB, ref.MongoAchievementID)
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"reference": ref, "detail": detail}})
}

// POST /api/v1/achievements (Create Draft)
func CreateAchievement(c *fiber.Ctx) error {
	var req model.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}

	// Logic Auth: Pastikan student_id di body sesuai dengan user yang login (Optional security check)
	
	// 1. Simpan ke Mongo
	mongoData := model.Achievement{
		StudentID: req.StudentID, AchievementType: req.AchievementType,
		Title: req.Title, Description: req.Description, Details: req.Details,
		Tags: req.Tags, Points: req.Points, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	mongoID, err := repository.InsertAchievementMongo(database.MongoDB, mongoData)
	if err != nil { return c.Status(500).JSON(fiber.Map{"success": false, "message": "Mongo Error"})}

	// 2. Simpan ke Postgres (Status Draft)
	ref := &model.AchievementReference{
		StudentID: req.StudentID, MongoAchievementID: mongoID, Status: model.StatusDraft,
	}
	if err := repository.CreateAchievementReference(database.DB, ref); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Postgres Error"})}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Draft created", "data": ref})
}

// PUT /api/v1/achievements/:id (Update Draft)
func UpdateAchievement(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	var req model.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil { return c.SendStatus(400) }

	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil { return c.SendStatus(404) }

	if ref.Status != model.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Hanya draft yang bisa diedit"})
	}

	// Update Mongo
	updateData := map[string]interface{}{
		"title": req.Title, "description": req.Description, 
		"details": req.Details, "tags": req.Tags, "points": req.Points,
	}
	if err := repository.UpdateAchievementMongo(database.MongoDB, ref.MongoAchievementID, updateData); err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Prestasi updated"})
}

// DELETE /api/v1/achievements/:id
func DeleteAchievement(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil { return c.SendStatus(404) }

	if ref.Status != model.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Hanya draft yang bisa dihapus"})
	}

	// Delete Mongo & Postgres
	repository.DeleteAchievementMongo(database.MongoDB, ref.MongoAchievementID)
	repository.DeleteAchievementReference(database.DB, id)

	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

// POST /api/v1/achievements/:id/submit
func SubmitAchievement(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil { return c.SendStatus(404) }

	if ref.Status != model.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Hanya draft yang bisa disubmit"})
	}

	now := time.Now()
	ref.Status = model.StatusSubmitted
	ref.SubmittedAt = &now
	repository.UpdateAchievementStatus(database.DB, ref)

	return c.JSON(fiber.Map{"success": true, "message": "Submitted for verification"})
}

// POST /api/v1/achievements/:id/verify
func VerifyAchievement(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	verifierIDStr := c.Locals("user_id").(string) // Ambil ID Dosen dari Token
	verifierID, _ := uuid.Parse(verifierIDStr)

	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil { return c.SendStatus(404) }

	if ref.Status != model.StatusSubmitted {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Status harus submitted"})
	}

	now := time.Now()
	ref.Status = model.StatusVerified
	ref.VerifiedAt = &now
	ref.VerifiedBy = &verifierID
	repository.UpdateAchievementStatus(database.DB, ref)

	return c.JSON(fiber.Map{"success": true, "message": "Verified"})
}

// POST /api/v1/achievements/:id/reject
func RejectAchievement(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	var req model.RejectAchievementRequest
	c.BodyParser(&req)

	verifierIDStr := c.Locals("user_id").(string)
	verifierID, _ := uuid.Parse(verifierIDStr)

	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil { return c.SendStatus(404) }

	if ref.Status != model.StatusSubmitted {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Status harus submitted"})
	}

	now := time.Now()
	ref.Status = model.StatusRejected
	ref.VerifiedAt = &now
	ref.VerifiedBy = &verifierID
	ref.RejectionNote = &req.RejectionNote
	repository.UpdateAchievementStatus(database.DB, ref)

	return c.JSON(fiber.Map{"success": true, "message": "Rejected"})
}

// GET /api/v1/achievements/:id/history
func GetAchievementHistory(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil { return c.SendStatus(404) }

	// Construct History Manual dari Data Reference
	var history []model.AchievementHistoryResponse

	// 1. Created
	history = append(history, model.AchievementHistoryResponse{
		Status: model.StatusDraft, Timestamp: ref.CreatedAt, Actor: "Mahasiswa",
	})

	// 2. Submitted
	if ref.SubmittedAt != nil {
		history = append(history, model.AchievementHistoryResponse{
			Status: model.StatusSubmitted, Timestamp: *ref.SubmittedAt, Actor: "Mahasiswa",
		})
	}

	// 3. Verified/Rejected
	if ref.VerifiedAt != nil {
		actor := "Dosen"
		if ref.Status == model.StatusRejected {
			history = append(history, model.AchievementHistoryResponse{
				Status: model.StatusRejected, Timestamp: *ref.VerifiedAt, Actor: actor, Note: ref.RejectionNote,
			})
		} else {
			history = append(history, model.AchievementHistoryResponse{
				Status: model.StatusVerified, Timestamp: *ref.VerifiedAt, Actor: actor,
			})
		}
	}

	return c.JSON(fiber.Map{"success": true, "data": history})
}

// POST /api/v1/achievements/:id/attachments
func UploadAttachment(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil { return c.SendStatus(404) }

	file, err := c.FormFile("file")
	if err != nil { return c.Status(400).JSON(fiber.Map{"message": "File required"})}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	c.SaveFile(file, filepath.Join("./uploads", filename))

	attachment := model.Attachment{
		FileName: file.Filename, FileURL: "/uploads/"+filename, FileType: file.Header.Get("Content-Type"), UploadedAt: time.Now(),
	}

	repository.AddAttachmentMongo(database.MongoDB, ref.MongoAchievementID, attachment)
	return c.JSON(fiber.Map{"success": true, "message": "File uploaded", "data": attachment})
}