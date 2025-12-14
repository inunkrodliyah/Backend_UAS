package repository_test

import (
	"errors"
	"testing"
	"time"

	"project-uas/app/model"

	"github.com/google/uuid"
)

/* ============================================================
   MOCK ACHIEVEMENT REPOSITORY (POSTGRES)
============================================================ */

type MockAchievementRepo struct {
	data map[uuid.UUID]*model.AchievementReference
}

func NewMockAchievementRepo() *MockAchievementRepo {
	return &MockAchievementRepo{
		data: make(map[uuid.UUID]*model.AchievementReference),
	}
}

func (m *MockAchievementRepo) Create(r *model.AchievementReference) error {
	if r.StudentID == uuid.Nil {
		return errors.New("student_id required")
	}
	r.ID = uuid.New()
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	m.data[r.ID] = r
	return nil
}

func (m *MockAchievementRepo) GetByID(id uuid.UUID) (*model.AchievementReference, error) {
	if r, ok := m.data[id]; ok && r.Status != model.StatusDeleted {
		return r, nil
	}
	return nil, errors.New("not found")
}

func (m *MockAchievementRepo) SoftDelete(id uuid.UUID) error {
	r, ok := m.data[id]
	if !ok {
		return errors.New("not found")
	}
	now := time.Now()
	r.Status = model.StatusDeleted
	r.DeletedAt = &now
	return nil
}

/* ===================== TEST ACHIEVEMENT ===================== */

func TestAchievementRepository_CreateAndGet(t *testing.T) {
	repo := NewMockAchievementRepo()

	ref := &model.AchievementReference{
		StudentID: uuid.New(),
		Status:    model.StatusDraft,
	}

	err := repo.Create(ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.GetByID(ref.ID)
	if err != nil {
		t.Fatalf("should be found")
	}

	if found.Status != model.StatusDraft {
		t.Fatalf("status mismatch")
	}
}

func TestAchievementRepository_SoftDelete(t *testing.T) {
	repo := NewMockAchievementRepo()

	ref := &model.AchievementReference{
		StudentID: uuid.New(),
		Status:    model.StatusDraft,
	}
	repo.Create(ref)

	err := repo.SoftDelete(ref.ID)
	if err != nil {
		t.Fatalf("delete failed")
	}

	_, err = repo.GetByID(ref.ID)
	if err == nil {
		t.Fatalf("expected not found after delete")
	}
}

/* ============================================================
   MOCK USER / LOGIN REPOSITORY
============================================================ */

type MockUserRepo struct {
	users map[string]*model.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users: make(map[string]*model.User),
	}
}

func (m *MockUserRepo) GetByUsername(username string) (*model.User, error) {
	u, ok := m.users[username]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

/* ========================= TEST LOGIN ======================= */

func TestUserRepository_GetByUsername(t *testing.T) {
	repo := NewMockUserRepo()

	user := &model.User{
		ID:           uuid.New(),
		Username:     "admin",
		PasswordHash: "hashed",
		IsActive:     true,
	}
	repo.users["admin"] = user

	found, err := repo.GetByUsername("admin")
	if err != nil {
		t.Fatalf("user should exist")
	}

	if found.Username != "admin" {
		t.Fatalf("username mismatch")
	}
}
