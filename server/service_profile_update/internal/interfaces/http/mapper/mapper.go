package mapper

import (
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/interfaces/http/dto"
)

// ToCreateFaceProfileUpdateRequestInput maps DTO to application model
func ToCreateFaceProfileUpdateRequestInput(d *dto.CreateFaceProfileUpdateRequestDTO, session *model.SessionInfo) *model.CreateFaceProfileUpdateRequestInput {
	return &model.CreateFaceProfileUpdateRequestInput{
		Session: session,
		Reason:  d.Reason,
	}
}

// ToGetMyUpdateRequestsInput maps DTO to application model
func ToGetMyUpdateRequestsInput(d *dto.GetMyRequestsQueryDTO, session *model.SessionInfo) *model.GetMyUpdateRequestsInput {
	return &model.GetMyUpdateRequestsInput{
		Session: session,
		Month:   d.Month,
		Limit:   d.Limit,
		Offset:  d.Offset,
	}
}

// ToGetPendingRequestsInput maps DTO to application model
func ToGetPendingRequestsInput(d *dto.GetPendingRequestsQueryDTO, session *model.SessionInfo) *model.GetPendingRequestsInput {
	return &model.GetPendingRequestsInput{
		Session: session,
		Limit:   d.Limit,
		Offset:  d.Offset,
	}
}

// ToApproveRequestInput maps DTO to application model
func ToApproveRequestInput(d *dto.ApproveRequestParamDTO, session *model.SessionInfo) *model.ApproveRequestInput {
	return &model.ApproveRequestInput{
		Session:   session,
		RequestID: d.RequestID,
	}
}

// ToRejectRequestInput maps DTO to application model
func ToRejectRequestInput(paramDTO *dto.RejectRequestParamDTO, bodyDTO *dto.RejectRequestBodyDTO, session *model.SessionInfo) *model.RejectRequestInput {
	return &model.RejectRequestInput{
		Session:   session,
		RequestID: paramDTO.RequestID,
		Reason:    bodyDTO.Reason,
	}
}

// ToValidateUpdateTokenInput maps DTO to application model
func ToValidateUpdateTokenInput(d *dto.ValidateTokenQueryDTO) *model.ValidateUpdateTokenInput {
	return &model.ValidateUpdateTokenInput{
		Token: d.Token,
	}
}

// ToUpdateFaceProfileInput maps DTO to application model
func ToUpdateFaceProfileInput(d *dto.UpdateFaceProfileFormDTO, imageData []byte, filename string) *model.UpdateFaceProfileInput {
	return &model.UpdateFaceProfileInput{
		Token:     d.Token,
		ImageData: imageData,
		Filename:  filename,
	}
}

// ToResetEmployeePasswordInput maps DTO to application model
func ToResetEmployeePasswordInput(d *dto.ResetPasswordRequestDTO, session *model.SessionInfo) *model.ResetEmployeePasswordInput {
	return &model.ResetEmployeePasswordInput{
		Session:    session,
		EmployeeID: d.EmployeeID,
	}
}
