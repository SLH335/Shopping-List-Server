package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	. "github.com/slh335/shoppinglistserver"
)

func (server *Server) GetInvitations(c echo.Context) error {
	user, success, err := verifySession(c, server)
	if !success {
		return err
	}

	invitations, err := server.InvitationService.GetInvitations(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load invitations",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    invitations,
	})
}

func (server *Server) Invite(c echo.Context) error {
	inviter, success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "username", "list_id")
	if !success {
		return err
	}
	username, listIdStr := values[0], values[1]
	listId, err := strconv.Atoi(listIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: path parameter '%s' is not a valid integer", listIdStr),
		})
	}

	members, err := server.ListService.Members(listId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load list members",
		})
	}
	foundInviter := false
	for _, member := range members {
		if member.Id == inviter.Id {
			foundInviter = true
		}
	}
	if !foundInviter {
		return c.JSON(http.StatusForbidden, Response{
			Success: false,
			Message: "error: user is not authorized to invite to this list",
		})
	}
	for _, member := range members {
		if member.Username == username {
			return c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "error: that user is already a member of the list",
			})
		}
	}

	invitee, err := server.UserService.GetUser(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load user",
		})
	}

	invitation, err := server.InvitationService.AddInvitation(inviter.Id, invitee.Id, listId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to create invitation",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully invited %s to '%s'", invitation.Invitee.Username, invitation.List.Name),
		Data:    invitation,
	})
}

func (server *Server) AcceptInvitation(c echo.Context) error {
	invitee, success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "invitation_token")
	if !success {
		return err
	}
	token := values[0]

	invitation, err := server.InvitationService.GetInvitation(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load invitation",
		})
	}

	if invitation.Invitee.Id != invitee.Id {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: only the invitee can accept an invitation",
		})
	}

	err = server.ListService.Join(invitation.List.Id, invitee.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to join list",
		})
	}

	err = server.InvitationService.DeleteInvitation(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to delete invitation",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully accepted invitation to \"%s\"", invitation.List.Name),
		Data:    invitation,
	})
}

func (server *Server) DeclineInvitation(c echo.Context) error {
	invitee, success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "invitation_token")
	if !success {
		return err
	}
	token := values[0]

	invitation, err := server.InvitationService.GetInvitation(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load invitation",
		})
	}

	if invitation.Invitee.Id != invitee.Id {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: only the invitee can decline an invitation",
		})
	}

	err = server.InvitationService.DeleteInvitation(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to delete invitation",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "successfully declined invitation",
	})
}

func (server *Server) RevokeInvitation(c echo.Context) error {
	inviter, success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "invitation_token")
	if !success {
		return err
	}
	token := values[0]

	invitation, err := server.InvitationService.GetInvitation(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load invitation",
		})
	}

	if invitation.Invitee.Id != inviter.Id {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: only the inviter can revoke an invitation",
		})
	}

	err = server.InvitationService.DeleteInvitation(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to delete invitation",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "successfully revoked invitation",
	})
}
