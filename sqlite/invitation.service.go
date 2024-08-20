package sqlite

import (
	"database/sql"
	"fmt"

	. "github.com/slh335/shoppinglistserver"
	"github.com/slh335/shoppinglistserver/crypto"
)

type InvitationService struct {
	DB *sql.DB
}

func (m *InvitationService) GetInvitation(token string) (invitation Invitation, err error) {
	stmt := `
		SELECT invitations.token, inviter.id, inviter.username, invitee.id, invitee.username, lists.id, lists.name
		FROM invitations
		INNER JOIN users inviter ON invitations.inviter_id=inviter.id
		INNER JOIN users invitee ON invitations.invitee_id=invitee.id
		INNER JOIN lists ON invitations.list_id=lists.id
		WHERE invitations.token=?`
	row := m.DB.QueryRow(stmt, token)

	err = row.Scan(
		&invitation.Token,
		&invitation.Inviter.Id, &invitation.Inviter.Username,
		&invitation.Invitee.Id, &invitation.Invitee.Username,
		&invitation.List.Id, &invitation.List.Name,
	)
	if err != nil {
		return invitation, err
	}
	return invitation, nil
}

func (m *InvitationService) GetInvitations(userId int) (invitations []Invitation, err error) {
	stmt := `
		SELECT invitations.token, inviter.id, inviter.username, invitee.id, invitee.username, lists.id, lists.name
		FROM invitations
		INNER JOIN users inviter ON invitations.inviter_id=inviter.id
		INNER JOIN users invitee ON invitations.invitee_id=invitee.id
		INNER JOIN lists ON invitations.list_id=lists.id
		WHERE inviter.id=? OR invitee.id=?`
	rows, err := m.DB.Query(stmt, userId, userId)
	if err != nil {
		return []Invitation{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var invitation Invitation
		err = rows.Scan(
			&invitation.Token,
			&invitation.Inviter.Id, &invitation.Inviter.Username,
			&invitation.Invitee.Id, &invitation.Invitee.Username,
			&invitation.List.Id, &invitation.List.Name,
		)
		invitations = append(invitations, invitation)
	}

	err = rows.Err()
	if err != nil {
		return invitations, err
	}
	return invitations, nil
}

func (m *InvitationService) AddInvitation(inviterId, inviteeId, listId int) (Invitation, error) {
	token := crypto.GenerateToken(64)

	stmt := "INSERT INTO invitations (token, inviter_id, invitee_id, list_id) VALUES (?, ?, ?, ?)"
	_, err := m.DB.Exec(stmt, token, inviterId, inviteeId, listId)
	if err != nil {
		return Invitation{}, err
	}

	return m.GetInvitation(token)
}

func (m *InvitationService) DeleteInvitation(token string) error {
	stmt := "DELETE FROM invitations WHERE token=?"
	res, err := m.DB.Exec(stmt, token)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("error: failed to delete invitation")
	}

	return nil
}
