package service_test

import (
	"errors"
	"testing"

	"project-uas/app/model"

	"github.com/google/uuid"
)

/* ============================================================
   MOCK HELPERS (PASSWORD & TOKEN)
============================================================ */

func mockCheckPassword(input, hashed string) bool {
	return input == "password123" && hashed == "hashed"
}

func mockGenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("invalid user")
	}
	return "mock-token", nil
}

/* ============================================================
   SERVICE LOGIN (TEST-ONLY)
============================================================ */

func LoginServiceTest(repo func(string) (*model.User, error), username, password string) (string, error) {
	user, err := repo(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !mockCheckPassword(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", errors.New("inactive user")
	}

	return mockGenerateToken(user.ID.String())
}

/* ========================== TEST LOGIN ====================== */

func TestLoginService_Success(t *testing.T) {
	mockRepo := func(username string) (*model.User, error) {
		return &model.User{
			ID:           uuid.New(),
			Username:     username,
			PasswordHash: "hashed",
			IsActive:     true,
		}, nil
	}

	token, err := LoginServiceTest(mockRepo, "admin", "password123")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if token == "" {
		t.Fatalf("token should not be empty")
	}
}

func TestLoginService_WrongPassword(t *testing.T) {
	mockRepo := func(username string) (*model.User, error) {
		return &model.User{
			ID:           uuid.New(),
			PasswordHash: "hashed",
			IsActive:     true,
		}, nil
	}

	_, err := LoginServiceTest(mockRepo, "admin", "wrong")
	if err == nil {
		t.Fatalf("expected error")
	}
}

/* ============================================================
   SERVICE ACHIEVEMENT (TEST-ONLY)
============================================================ */

func CreateAchievementService(studentID uuid.UUID) (*model.AchievementReference, error) {
	if studentID == uuid.Nil {
		return nil, errors.New("student required")
	}

	return &model.AchievementReference{
		ID:        uuid.New(),
		StudentID: studentID,
		Status:    model.StatusDraft,
	}, nil
}

/* ====================== TEST ACHIEVEMENT ==================== */

func TestCreateAchievementService(t *testing.T) {
	ref, err := CreateAchievementService(uuid.New())
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if ref.Status != model.StatusDraft {
		t.Fatalf("status should be draft")
	}
}

func TestCreateAchievementService_Invalid(t *testing.T) {
	_, err := CreateAchievementService(uuid.Nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}
