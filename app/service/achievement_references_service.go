package service

import (
	"log"
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"strings" 
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
	return c.JSON(fiber.Map{"success": true, "data": ref})
}

// CreateAchievementReference (Logika 'Submit' oleh Mahasiswa)
func CreateAchievementReference(c *fiber.Ctx) error {
	// 1. Menggunakan struct request 'Submit' yang baru
	var req model.CreateAchievementReferenceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	// 2. Validasi field wajib
	if req.StudentID == uuid.Nil || req.MongoAchievementID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Field 'student_id' dan 'mongo_achievement_id' wajib diisi",
		})
	}

	now := time.Now()

	// 3. Siapkan model.AchievementReference untuk disimpan
	ref := &model.AchievementReference{
		StudentID:          req.StudentID,
		MongoAchievementID: req.MongoAchievementID,
		Status:             model.StatusPending, // Status default saat submit
		SubmittedAt:        &now,                // Waktu submit
		// verified_at, verified_by, rejection_note otomatis NULL
	}

	// 4. Panggil repository (yang sudah diperbarui)
	if err := repository.CreateAchievementReference(database.DB, ref); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false, "message": "Prestasi ini sudah pernah disubmit.", "error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menambah referensi", "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Prestasi berhasil disubmit untuk verifikasi", "data": ref})
}

// UpdateAchievementReference (Logika 'Verify/Reject' oleh Dosen/Admin)
func UpdateAchievementReference(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID referensi tidak valid"})
	}

	// 1. Menggunakan struct request 'Update Status' yang baru
	var req model.UpdateAchievementStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	// 2. Validasi
	if req.Status != model.StatusApproved && req.Status != model.StatusRejected {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Status harus 'approved' atau 'rejected'",
		})
	}
	if req.Status == model.StatusRejected && (req.RejectionNote == nil || *req.RejectionNote == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Catatan penolakan (rejection_note) wajib diisi jika status 'rejected'",
		})
	}
	if req.VerifiedBy == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Field 'verified_by' (ID verifikator) wajib diisi",
		})
	}

	// 3. Ambil data lama
	ref, err := repository.GetAchievementReferenceByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Referensi tidak ditemukan"})
	}
	
	if ref.Status != model.StatusPending {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false, "message": "Gagal: Prestasi ini sudah diverifikasi/ditolak sebelumnya.",
		})
	}

	now := time.Now()

	// 4. Update data
	ref.Status = req.Status
	ref.VerifiedBy = &req.VerifiedBy
	ref.VerifiedAt = &now
	if req.Status == model.StatusRejected {
		ref.RejectionNote = req.RejectionNote
	}

	// 5. Panggil repository (yang sudah diperbarui)
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