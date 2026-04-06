package grouprepo

import (
	"context"
	"database/sql"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

type dbQuerier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

//--------------------------------------------------------------------------------------|

type sqlGroupRepository struct {
	db dbQuerier
}

//--------------------------------------------------------------------------------------|

// NewGroupRepository creates a new SQL-backed group repository.
func NewGroupRepository(db *sql.DB) models.GroupRepo {
	return &sqlGroupRepository{db: db}
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) WithTx(tx any) models.GroupRepo {
	if tx == nil {
		return r
	}
	if t, ok := tx.(*sql.Tx); ok {
		return &sqlGroupRepository{db: t}
	}
	return r
}

//--------------------------------------------------------------------------------------|
// Group CRUD
//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) CreateGroup(ctx context.Context, group *models.Group) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO groups (creator_id, title, description) VALUES (?, ?, ?) RETURNING id, created_at`,
		group.CreatorID, group.Title, group.Description,
	).Scan(&group.ID, &group.CreatedAt)
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetGroupByID(ctx context.Context, id int) (*models.Group, error) {
	var g models.Group
	err := r.db.QueryRowContext(ctx,
		`SELECT id, creator_id, title, description, created_at FROM groups WHERE id = ?`, id,
	).Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) ListGroups(ctx context.Context, limit, offset int) ([]models.Group, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, creator_id, title, description, created_at FROM groups ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

//--------------------------------------------------------------------------------------|
// Membership
//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) AddMember(ctx context.Context, groupID, userID int, role string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO group_members (group_id, user_id, role) VALUES (?, ?, ?)`,
		groupID, userID, role)
	return err
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) RemoveMember(ctx context.Context, groupID, userID int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, userID)
	return err
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetMembers(ctx context.Context, groupID int) ([]models.GroupMember, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT gm.group_id, gm.user_id, gm.role, gm.joined_at, u.username, u.first_name, u.last_name
		 FROM group_members gm
		 JOIN users u ON gm.user_id = u.id
		 WHERE gm.group_id = ?
		 ORDER BY gm.joined_at`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.GroupMember
	for rows.Next() {
		var m models.GroupMember
		if err := rows.Scan(&m.GroupID, &m.UserID, &m.Role, &m.JoinedAt, &m.Username, &m.FirstName, &m.LastName); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) IsMember(ctx context.Context, groupID, userID int) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?`,
		groupID, userID).Scan(&count)
	return count > 0, err
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetMemberGroupIDs(ctx context.Context, userID int) ([]int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT group_id FROM group_members WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

//--------------------------------------------------------------------------------------|
// Invitations
//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) CreateInvitation(ctx context.Context, inv *models.GroupInvitation) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO group_invitations (group_id, inviter_id, invitee_id) VALUES (?, ?, ?) RETURNING id, created_at`,
		inv.GroupID, inv.InviterID, inv.InviteeID,
	).Scan(&inv.ID, &inv.CreatedAt)
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetInvitationByID(ctx context.Context, id int) (*models.GroupInvitation, error) {
	var inv models.GroupInvitation
	err := r.db.QueryRowContext(ctx,
		`SELECT gi.id, gi.group_id, gi.inviter_id, gi.invitee_id, gi.status, gi.created_at,
		        g.title, u.username
		 FROM group_invitations gi
		 JOIN groups g ON gi.group_id = g.id
		 JOIN users u ON gi.inviter_id = u.id
		 WHERE gi.id = ?`, id,
	).Scan(&inv.ID, &inv.GroupID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt,
		&inv.GroupTitle, &inv.InviterName)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetPendingInvitations(ctx context.Context, userID int) ([]models.GroupInvitation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT gi.id, gi.group_id, gi.inviter_id, gi.invitee_id, gi.status, gi.created_at,
		        g.title, u.username
		 FROM group_invitations gi
		 JOIN groups g ON gi.group_id = g.id
		 JOIN users u ON gi.inviter_id = u.id
		 WHERE gi.invitee_id = ? AND gi.status = 'pending'
		 ORDER BY gi.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []models.GroupInvitation
	for rows.Next() {
		var inv models.GroupInvitation
		if err := rows.Scan(&inv.ID, &inv.GroupID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt,
			&inv.GroupTitle, &inv.InviterName); err != nil {
			return nil, err
		}
		invitations = append(invitations, inv)
	}
	return invitations, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) UpdateInvitationStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE group_invitations SET status = ? WHERE id = ?`, status, id)
	return err
}

//--------------------------------------------------------------------------------------|
// Join Requests
//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) CreateJoinRequest(ctx context.Context, req *models.GroupJoinRequest) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO group_join_requests (group_id, user_id) VALUES (?, ?) RETURNING id, created_at`,
		req.GroupID, req.UserID,
	).Scan(&req.ID, &req.CreatedAt)
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetJoinRequestByID(ctx context.Context, id int) (*models.GroupJoinRequest, error) {
	var req models.GroupJoinRequest
	err := r.db.QueryRowContext(ctx,
		`SELECT gjr.id, gjr.group_id, gjr.user_id, gjr.status, gjr.created_at, u.username
		 FROM group_join_requests gjr
		 JOIN users u ON gjr.user_id = u.id
		 WHERE gjr.id = ?`, id,
	).Scan(&req.ID, &req.GroupID, &req.UserID, &req.Status, &req.CreatedAt, &req.Username)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetPendingJoinRequests(ctx context.Context, groupID int) ([]models.GroupJoinRequest, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT gjr.id, gjr.group_id, gjr.user_id, gjr.status, gjr.created_at, u.username
		 FROM group_join_requests gjr
		 JOIN users u ON gjr.user_id = u.id
		 WHERE gjr.group_id = ? AND gjr.status = 'pending'
		 ORDER BY gjr.created_at DESC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.GroupJoinRequest
	for rows.Next() {
		var req models.GroupJoinRequest
		if err := rows.Scan(&req.ID, &req.GroupID, &req.UserID, &req.Status, &req.CreatedAt, &req.Username); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) UpdateJoinRequestStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE group_join_requests SET status = ? WHERE id = ?`, status, id)
	return err
}

//--------------------------------------------------------------------------------------|
// Events
//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) CreateEvent(ctx context.Context, event *models.GroupEvent) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO group_events (group_id, creator_id, title, description, event_time) 
		 VALUES (?, ?, ?, ?, ?) RETURNING id, created_at`,
		event.GroupID, event.CreatorID, event.Title, event.Description, event.EventTime,
	).Scan(&event.ID, &event.CreatedAt)
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetEventByID(ctx context.Context, id int) (*models.GroupEvent, error) {
	var e models.GroupEvent
	err := r.db.QueryRowContext(ctx,
		`SELECT ge.id, ge.group_id, ge.creator_id, ge.title, ge.description, ge.event_time, ge.created_at,
		        COALESCE(SUM(CASE WHEN ger.response = 'going' THEN 1 ELSE 0 END), 0),
		        COALESCE(SUM(CASE WHEN ger.response = 'not_going' THEN 1 ELSE 0 END), 0)
		 FROM group_events ge
		 LEFT JOIN group_event_responses ger ON ge.id = ger.event_id
		 WHERE ge.id = ?
		 GROUP BY ge.id`, id,
	).Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventTime, &e.CreatedAt,
		&e.GoingCount, &e.NotGoingCount)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetGroupEvents(ctx context.Context, groupID int) ([]models.GroupEvent, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT ge.id, ge.group_id, ge.creator_id, ge.title, ge.description, ge.event_time, ge.created_at,
		        COALESCE(SUM(CASE WHEN ger.response = 'going' THEN 1 ELSE 0 END), 0),
		        COALESCE(SUM(CASE WHEN ger.response = 'not_going' THEN 1 ELSE 0 END), 0)
		 FROM group_events ge
		 LEFT JOIN group_event_responses ger ON ge.id = ger.event_id
		 WHERE ge.group_id = ?
		 GROUP BY ge.id
		 ORDER BY ge.event_time`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.GroupEvent
	for rows.Next() {
		var e models.GroupEvent
		if err := rows.Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventTime, &e.CreatedAt,
			&e.GoingCount, &e.NotGoingCount); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) RespondToEvent(ctx context.Context, resp *models.GroupEventResponse) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO group_event_responses (event_id, user_id, response) VALUES (?, ?, ?)
		 ON CONFLICT (event_id, user_id) DO UPDATE SET response = excluded.response`,
		resp.EventID, resp.UserID, resp.Response)
	return err
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetEventResponses(ctx context.Context, eventID int) ([]models.GroupEventResponse, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT ger.event_id, ger.user_id, ger.response, u.username
		 FROM group_event_responses ger
		 JOIN users u ON ger.user_id = u.id
		 WHERE ger.event_id = ?`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []models.GroupEventResponse
	for rows.Next() {
		var r models.GroupEventResponse
		if err := rows.Scan(&r.EventID, &r.UserID, &r.Response, &r.Username); err != nil {
			return nil, err
		}
		responses = append(responses, r)
	}
	return responses, nil
}

//--------------------------------------------------------------------------------------|
// Group Messages
//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) SaveGroupMessage(ctx context.Context, msg *models.GroupMessage) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO group_messages (group_id, sender_id, body, image_url) VALUES (?, ?, ?, ?) 
		 RETURNING id, created_at`,
		msg.GroupID, msg.SenderID, msg.Body, msg.ImageURL,
	).Scan(&msg.ID, &msg.CreatedAt)
}

//--------------------------------------------------------------------------------------|

func (r *sqlGroupRepository) GetGroupMessages(ctx context.Context, groupID, limit, offset int) ([]models.GroupMessage, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT gm.id, gm.group_id, gm.sender_id, gm.body, gm.image_url, gm.created_at, u.username
		 FROM group_messages gm
		 JOIN users u ON gm.sender_id = u.id
		 WHERE gm.group_id = ?
		 ORDER BY gm.created_at DESC, gm.id DESC
		 LIMIT ? OFFSET ?`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.GroupMessage
	for rows.Next() {
		var m models.GroupMessage
		if err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.Body, &m.ImageURL, &m.CreatedAt, &m.Username); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
