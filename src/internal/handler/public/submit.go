package public

import (
	context "silent/src/internal/context/abstract"
	"silent/src/internal/dto"
	"silent/src/internal/service/abstract"
)

func SubmitProvider(ctx context.HandlerContext, params dto.SubmitUpdateRequest, updateService abstract.IUpdateService) (dto.SubmitUpdateResponse, error) {
	err := updateService.Submit(params.TelegramUserID, params.RawText)
	if err != nil {
		return dto.SubmitUpdateResponse{}, err
	}

	return buildResponse(), nil
}

func buildResponse() dto.SubmitUpdateResponse {
	return dto.SubmitUpdateResponse{
		Status: "queued",
	}
}
