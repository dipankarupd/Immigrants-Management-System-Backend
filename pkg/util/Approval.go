package util

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ApprovalRequest struct {
	Approval string `json:"approval"`
}

func GetApprovalFromRequest(r *http.Request) (string, error) {
	reqBody := ApprovalRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		return "", err
	}

	if reqBody.Approval != "approved" && reqBody.Approval != "pending" && reqBody.Approval != "rejected" {
		return "", errors.New("invalid approval status")
	}

	return reqBody.Approval, nil
}
