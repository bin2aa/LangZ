package repository

import (
	"database/sql"
	"fmt"
	"time"

	"thinh/gin-app/internal/model"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(group *model.Group) error {
	query := `
		INSERT INTO groups (id, name, description, center_location, radius_meters, created_by, created_at)
		VALUES ($1, $2, $3, ST_GeogFromText($4), $5, $6, $7)
		RETURNING id, created_at`
	return r.db.QueryRow(query,
		group.ID, group.Name, group.Description, group.CenterLocation,
		group.RadiusMeters, group.CreatedBy, group.CreatedAt,
	).Scan(&group.ID, &group.CreatedAt)
}

func (r *GroupRepository) FindByID(id string) (*model.Group, error) {
	query := `SELECT id, name, description, ST_AsText(center_location) as center_location, 
	           radius_meters, created_by, created_at FROM groups WHERE id = $1`
	group := &model.Group{}
	var description, centerLocation sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&group.ID, &group.Name, &description, &centerLocation,
		&group.RadiusMeters, &group.CreatedBy, &group.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find group by id: %w", err)
	}
	if description.Valid {
		group.Description = &description.String
	}
	if centerLocation.Valid {
		group.CenterLocation = &centerLocation.String
	}
	return group, nil
}

func (r *GroupRepository) FindAll(page, pageSize, offset int) ([]model.Group, int, error) {
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM groups`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count groups: %w", err)
	}

	query := `SELECT g.id, g.name, g.description, ST_AsText(g.center_location) as center_location, 
	           g.radius_meters, g.created_by, g.created_at,
	           COUNT(gm.user_id) as member_count
	           FROM groups g
	           LEFT JOIN group_members gm ON g.id = gm.group_id
	           GROUP BY g.id
	           ORDER BY g.created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list groups: %w", err)
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		var description, centerLocation sql.NullString
		var memberCount int
		if err := rows.Scan(
			&g.ID, &g.Name, &description, &centerLocation,
			&g.RadiusMeters, &g.CreatedBy, &g.CreatedAt, &memberCount,
		); err != nil {
			return nil, 0, fmt.Errorf("scan group: %w", err)
		}
		if description.Valid {
			g.Description = &description.String
		}
		if centerLocation.Valid {
			g.CenterLocation = &centerLocation.String
		}
		groups = append(groups, g)
	}
	return groups, total, nil
}

func (r *GroupRepository) Update(group *model.Group) error {
	query := `UPDATE groups SET name = $1, description = $2, center_location = ST_GeogFromText($3), 
	           radius_meters = $4 WHERE id = $5`
	_, err := r.db.Exec(query, group.Name, group.Description, group.CenterLocation,
		group.RadiusMeters, group.ID)
	if err != nil {
		return fmt.Errorf("update group: %w", err)
	}
	return nil
}

func (r *GroupRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM groups WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete group: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Join group
func (r *GroupRepository) AddMember(groupID, userID string) error {
	query := `INSERT INTO group_members (group_id, user_id, joined_at) VALUES ($1, $2, $3) 
	           ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, groupID, userID, time.Now())
	return err
}

// Leave group
func (r *GroupRepository) RemoveMember(groupID, userID string) error {
	_, err := r.db.Exec(`DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`, groupID, userID)
	return err
}

// Check if user is member
func (r *GroupRepository) IsMember(groupID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM group_members WHERE group_id = $1 AND user_id = $2`, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// List members of a group
func (r *GroupRepository) GetMembers(groupID string) ([]model.GroupMember, error) {
	rows, err := r.db.Query(`SELECT group_id, user_id, joined_at FROM group_members WHERE group_id = $1 ORDER BY joined_at`, groupID)
	if err != nil {
		return nil, fmt.Errorf("get group members: %w", err)
	}
	defer rows.Close()

	var members []model.GroupMember
	for rows.Next() {
		var m model.GroupMember
		if err := rows.Scan(&m.GroupID, &m.UserID, &m.JoinedAt); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		members = append(members, m)
	}
	return members, nil
}
