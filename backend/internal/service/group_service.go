package service

import (
	"time"

	"github.com/google/uuid"

	"thinh/gin-app/internal/dto"
	"thinh/gin-app/internal/model"
	"thinh/gin-app/internal/repository"
	apperrors "thinh/gin-app/pkg/errors"
)

type GroupService struct {
	groupRepo *repository.GroupRepository
	userRepo  *repository.UserRepository
}

func NewGroupService(groupRepo *repository.GroupRepository, userRepo *repository.UserRepository) *GroupService {
	return &GroupService{groupRepo: groupRepo, userRepo: userRepo}
}

func (s *GroupService) Create(userID string, req *dto.CreateGroupRequest) (*dto.GroupResponse, error) {
	now := time.Now()
	group := &model.Group{
		ID:             uuid.New().String(),
		Name:           req.Name,
		Description:    req.Description,
		CenterLocation: req.CenterLocation,
		RadiusMeters:   req.RadiusMeters,
		CreatedBy:      userID,
		CreatedAt:      now,
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	// Creator automatically joins the group
	if err := s.groupRepo.AddMember(group.ID, userID); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	return s.groupToResponse(group, 1), nil
}

func (s *GroupService) GetByID(id string) (*dto.GroupResponse, error) {
	group, err := s.groupRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if group == nil {
		return nil, apperrors.NewNotFound("Group not found")
	}

	// Get member count
	members, err := s.groupRepo.GetMembers(id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}

	return s.groupToResponse(group, len(members)), nil
}

func (s *GroupService) GetAll(page, pageSize int) ([]dto.GroupResponse, int, error) {
	offset := (page - 1) * pageSize
	groups, total, err := s.groupRepo.FindAll(page, pageSize, offset)
	if err != nil {
		return nil, 0, apperrors.NewInternal(err)
	}

	resp := make([]dto.GroupResponse, len(groups))
	for i, g := range groups {
		resp[i] = *s.groupToResponse(&g, 0)
	}
	return resp, total, nil
}

func (s *GroupService) Update(id string, req *dto.UpdateGroupRequest) (*dto.GroupResponse, error) {
	group, err := s.groupRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if group == nil {
		return nil, apperrors.NewNotFound("Group not found")
	}

	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Description != nil {
		group.Description = req.Description
	}
	if req.CenterLocation != nil {
		group.CenterLocation = req.CenterLocation
	}
	if req.RadiusMeters != nil {
		group.RadiusMeters = *req.RadiusMeters
	}

	if err := s.groupRepo.Update(group); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	members, _ := s.groupRepo.GetMembers(id)
	return s.groupToResponse(group, len(members)), nil
}

func (s *GroupService) Delete(id string) error {
	if err := s.groupRepo.Delete(id); err != nil {
		return apperrors.NewInternal(err)
	}
	return nil
}

func (s *GroupService) JoinGroup(groupID, userID string) error {
	group, err := s.groupRepo.FindByID(groupID)
	if err != nil {
		return apperrors.NewInternal(err)
	}
	if group == nil {
		return apperrors.NewNotFound("Group not found")
	}

	isMember, err := s.groupRepo.IsMember(groupID, userID)
	if err != nil {
		return apperrors.NewInternal(err)
	}
	if isMember {
		return apperrors.NewConflict("Already a member of this group")
	}

	return s.groupRepo.AddMember(groupID, userID)
}

func (s *GroupService) LeaveGroup(groupID, userID string) error {
	isMember, err := s.groupRepo.IsMember(groupID, userID)
	if err != nil {
		return apperrors.NewInternal(err)
	}
	if !isMember {
		return apperrors.NewNotFound("Not a member of this group")
	}

	return s.groupRepo.RemoveMember(groupID, userID)
}

func (s *GroupService) GetMembers(groupID string) ([]dto.GroupMemberResponse, error) {
	members, err := s.groupRepo.GetMembers(groupID)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}

	resp := make([]dto.GroupMemberResponse, len(members))
	for i, m := range members {
		user, err := s.userRepo.FindByID(m.UserID)
		if err != nil || user == nil {
			resp[i] = dto.GroupMemberResponse{
				UserID:   m.UserID,
				JoinedAt: m.JoinedAt,
			}
			continue
		}
		resp[i] = dto.GroupMemberResponse{
			UserID:   m.UserID,
			FullName: user.FullName,
			Email:    user.Email,
			JoinedAt: m.JoinedAt,
		}
	}
	return resp, nil
}

func (s *GroupService) groupToResponse(group *model.Group, memberCount int) *dto.GroupResponse {
	return &dto.GroupResponse{
		ID:             group.ID,
		Name:           group.Name,
		Description:    group.Description,
		CenterLocation: group.CenterLocation,
		RadiusMeters:   group.RadiusMeters,
		CreatedBy:      group.CreatedBy,
		CreatedAt:      group.CreatedAt,
		MemberCount:    memberCount,
	}
}
