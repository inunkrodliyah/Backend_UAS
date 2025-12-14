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
// ListAchievements godoc
// @Summary      Lihat Daftar Prestasi
// @Description  Menampilkan daftar prestasi berdasarkan role (Mahasiswa: milik sendiri, Dosen: milik bimbingan, Admin: semua)
// @Tags         Achievements
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  fiber.Map{data=[]model.AchievementReference}
// @Failure      500  {object}  fiber.Map
// @Router       /achievements [get]
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
// GetAchievementDetail godoc
// @Summary      Lihat Detail Prestasi
// @Description  Mendapatkan detail lengkap prestasi (gabungan data Postgres & MongoDB) berdasarkan ID
// @Tags         Achievements
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Achievement ID (UUID)"
// @Success      200  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /achievements/{id} [get]
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
// CreateAchievement godoc
// @Summary      Buat Prestasi Baru (Draft)
// @Description  Mahasiswa membuat draft prestasi baru
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body model.CreateAchievementRequest true "Data Prestasi"
// @Success      201  {object}  fiber.Map
// @Router       /achievements [post]
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
// UpdateAchievement godoc
// @Summary      Edit Prestasi (Hanya Draft)
// @Description  Mahasiswa mengedit data prestasi yang masih berstatus draft
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string true "Achievement ID"
// @Param        request body model.UpdateAchievementRequest true "Data Update"
// @Success      200  {object}  fiber.Map
// @Router       /achievements/{id} [put]
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
// DeleteAchievement godoc
// @Summary      Hapus Prestasi (Soft Delete)
// @Description  Menghapus prestasi (hanya draft) secara soft delete
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id   path      string true "Achievement ID"
// @Success      200  {object}  fiber.Map
// @Router       /achievements/{id} [delete]
func DeleteAchievement(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil { return c.SendStatus(404) }

	// Validasi: Hanya draft yang bisa dihapus (tetap dipertahankan)
	if ref.Status != model.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Hanya draft yang bisa dihapus"})
	}

	// 1. JANGAN Hapus data di Mongo jika ingin bisa di-restore nantinya.
	// Jika tetap ingin menghapus mongo secara permanen, biarkan baris ini.
	// repository.DeleteAchievementMongo(database.MongoDB, ref.MongoAchievementID) 
    
    // ATAU: Implementasikan Soft Delete juga di Mongo (menambah field deletedAt di Mongo).
    // Untuk sekarang, kita anggap Mongo dibiarkan saja (data sampah) atau dihapus permanen.
    // Jika Anda menghapus permanen Mongo:
    repository.DeleteAchievementMongo(database.MongoDB, ref.MongoAchievementID)

	// 2. Soft Delete di Postgres (Fungsi repository sudah diubah jadi UPDATE)
	err = repository.DeleteAchievementReference(database.DB, id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus"})
    }

	return c.JSON(fiber.Map{"success": true, "message": "Deleted (Soft)"})
}

// POST /api/v1/achievements/:id/submit
// SubmitAchievement godoc
// @Summary      Submit ke Dosen Wali
// @Description  Mengubah status Draft menjadi Submitted agar bisa diverifikasi
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id   path      string true "Achievement ID"
// @Success      200  {object}  fiber.Map
// @Router       /achievements/{id}/submit [post]
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
// VerifyAchievement godoc
// @Summary      Verifikasi Prestasi (Dosen)
// @Description  Dosen menyetujui prestasi mahasiswa
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id   path      string true "Achievement ID"
// @Success      200  {object}  fiber.Map
// @Router       /achievements/{id}/verify [post]
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
// RejectAchievement godoc
// @Summary      Tolak Prestasi (Dosen)
// @Description  Dosen menolak prestasi dengan catatan revisi
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id   path      string true "Achievement ID"
// @Param        request body model.RejectAchievementRequest true "Alasan Penolakan"
// @Success      200  {object}  fiber.Map
// @Router       /achievements/{id}/reject [post]
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
// GetHistory godoc
// @Summary      Lihat History Status
// @Description  Melihat jejak perubahan status (Draft -> Submitted -> Rejected -> dst)
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id   path      string true "Achievement ID"
// @Success      200  {array}   model.AchievementHistoryResponse
// @Router       /achievements/{id}/history [get]
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
// UploadAttachment godoc
// @Summary      Upload Bukti Dokumen
// @Description  Upload file PDF/Gambar bukti prestasi
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Param        id    path      string  true  "Achievement ID"
// @Param        file  formData  file    true  "Bukti File (PDF/JPG)"
// @Success      200   {object}  fiber.Map
// @Router       /achievements/{id}/attachments [post]
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