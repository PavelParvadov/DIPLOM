package usecase

import (
	"context"
	"testing"
	"time"

	"happyhouse/backend/internal/domain"
)

type fakeHouseRepo struct {
	memberships map[[2]int64]domain.Membership
	houses      map[int64]domain.House
	nextHouseID int64
}

func newFakeHouseRepo() *fakeHouseRepo {
	return &fakeHouseRepo{
		memberships: make(map[[2]int64]domain.Membership),
		houses:      make(map[int64]domain.House),
		nextHouseID: 1,
	}
}

func (r *fakeHouseRepo) Create(_ context.Context, createdBy int64, input domain.CreateHouseInput, defaultCategories []string) (*domain.House, error) {
	house := &domain.House{
		ID:        r.nextHouseID,
		Name:      input.Name,
		Address:   input.Address,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
	r.nextHouseID++
	r.houses[house.ID] = *house
	r.memberships[[2]int64{createdBy, house.ID}] = domain.Membership{
		UserID:   createdBy,
		HouseID:  house.ID,
		Role:     domain.RoleAdmin,
		JoinedAt: time.Now(),
	}
	return house, nil
}

func (r *fakeHouseRepo) ListByUser(_ context.Context, userID int64) ([]domain.UserHouse, error) {
	items := make([]domain.UserHouse, 0)
	for key, membership := range r.memberships {
		if key[0] == userID {
			items = append(items, domain.UserHouse{
				House: domain.House{ID: membership.HouseID, Name: "Дом"},
				Role:  membership.Role,
			})
		}
	}
	return items, nil
}

func (r *fakeHouseRepo) GetMembership(_ context.Context, userID, houseID int64) (*domain.Membership, error) {
	membership, ok := r.memberships[[2]int64{userID, houseID}]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &membership, nil
}

func (r *fakeHouseRepo) AddMembership(_ context.Context, userID, houseID int64, role string) error {
	if _, exists := r.memberships[[2]int64{userID, houseID}]; exists {
		return domain.ErrAlreadyMember
	}
	r.memberships[[2]int64{userID, houseID}] = domain.Membership{
		UserID:   userID,
		HouseID:  houseID,
		Role:     role,
		JoinedAt: time.Now(),
	}
	return nil
}

type fakeInviteRepo struct {
	invite *domain.InviteCode
}

func (r *fakeInviteRepo) ListByHouse(_ context.Context, houseID int64) ([]domain.InviteCode, error) {
	return []domain.InviteCode{}, nil
}

func (r *fakeInviteRepo) Create(_ context.Context, houseID, createdBy int64, input domain.CreateInviteCodeInput, code string) (*domain.InviteCode, error) {
	return nil, nil
}

func (r *fakeInviteRepo) GetActiveByCode(_ context.Context, code string) (*domain.InviteCode, error) {
	if r.invite == nil || r.invite.Code != code {
		return nil, domain.ErrNotFound
	}
	return r.invite, nil
}

func (r *fakeInviteRepo) Deactivate(_ context.Context, houseID, inviteCodeID int64) error {
	return nil
}

type fakeCategoryRepo struct{}

func (r *fakeCategoryRepo) ListByHouse(_ context.Context, houseID int64) ([]domain.Category, error) {
	return nil, nil
}

func (r *fakeCategoryRepo) Create(_ context.Context, houseID int64, input domain.CreateCategoryInput) (*domain.Category, error) {
	return &domain.Category{ID: 1, HouseID: houseID, Name: input.Name}, nil
}

func (r *fakeCategoryRepo) Update(_ context.Context, houseID, categoryID int64, input domain.UpdateCategoryInput) (*domain.Category, error) {
	return nil, nil
}

func (r *fakeCategoryRepo) Delete(_ context.Context, houseID, categoryID int64) error {
	return nil
}

func (r *fakeCategoryRepo) GetByID(_ context.Context, houseID, categoryID int64) (*domain.Category, error) {
	return &domain.Category{ID: categoryID, HouseID: houseID, Name: "Новости"}, nil
}

func TestHouseUseCaseJoinByCodeAddsMembership(t *testing.T) {
	houseRepo := newFakeHouseRepo()
	inviteRepo := &fakeInviteRepo{
		invite: &domain.InviteCode{
			ID:        1,
			HouseID:   7,
			Code:      "WELCOME7",
			IsActive:  true,
			CreatedAt: time.Now(),
		},
	}

	uc := NewHouseUseCase(houseRepo, inviteRepo)
	if err := uc.JoinByCode(context.Background(), 10, domain.JoinHouseInput{Code: "WELCOME7"}); err != nil {
		t.Fatalf("join by code returned error: %v", err)
	}

	membership, err := houseRepo.GetMembership(context.Background(), 10, 7)
	if err != nil {
		t.Fatalf("expected membership to be created: %v", err)
	}
	if membership.Role != domain.RoleResident {
		t.Fatalf("expected resident role, got %s", membership.Role)
	}
}

func TestHouseUseCaseCreateMakesCreatorAdmin(t *testing.T) {
	houseRepo := newFakeHouseRepo()
	uc := NewHouseUseCase(houseRepo, &fakeInviteRepo{})

	house, err := uc.Create(context.Background(), 22, domain.CreateHouseInput{
		Name:    "Дом на Невском",
		Address: "Невский проспект, 18",
	})
	if err != nil {
		t.Fatalf("create house returned error: %v", err)
	}
	if house.CreatedBy != 22 {
		t.Fatalf("expected creator id 22, got %d", house.CreatedBy)
	}

	membership, err := houseRepo.GetMembership(context.Background(), 22, house.ID)
	if err != nil {
		t.Fatalf("expected creator membership: %v", err)
	}
	if membership.Role != domain.RoleAdmin {
		t.Fatalf("expected admin role, got %s", membership.Role)
	}
}

func TestHouseUseCaseJoinByCodeReturnsFriendlyConflict(t *testing.T) {
	houseRepo := newFakeHouseRepo()
	houseRepo.memberships[[2]int64{10, 7}] = domain.Membership{
		UserID:   10,
		HouseID:  7,
		Role:     domain.RoleResident,
		JoinedAt: time.Now(),
	}
	inviteRepo := &fakeInviteRepo{
		invite: &domain.InviteCode{
			ID:        1,
			HouseID:   7,
			Code:      "WELCOME7",
			IsActive:  true,
			CreatedAt: time.Now(),
		},
	}

	uc := NewHouseUseCase(houseRepo, inviteRepo)
	err := uc.JoinByCode(context.Background(), 10, domain.JoinHouseInput{Code: "WELCOME7"})
	if err != domain.ErrAlreadyMember {
		t.Fatalf("expected already member error, got %v", err)
	}
}

func TestCategoryUseCaseCreateRequiresAdmin(t *testing.T) {
	houseRepo := newFakeHouseRepo()
	houseRepo.memberships[[2]int64{5, 1}] = domain.Membership{
		UserID:   5,
		HouseID:  1,
		Role:     domain.RoleResident,
		JoinedAt: time.Now(),
	}

	uc := NewCategoryUseCase(houseRepo, &fakeCategoryRepo{})
	_, err := uc.Create(context.Background(), 5, 1, domain.CreateCategoryInput{Name: "Объявления"})
	if err != domain.ErrForbidden {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}
