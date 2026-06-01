package public

import (
	context "silent/src/internal/context/abstract"
	"silent/src/internal/dto"
	"silent/src/internal/service/abstract"
)

func RegisterMemberHandler(ctx context.HandlerContext, req *dto.RegisterMemberRequest, memberService abstract.IMemberService) (dto.RegisterMemberResponse, error) {
	err := memberService.Register(req.TelegramUserID, req.Name, req.TeamID, req.IsLead)
	if err != nil {
		return dto.RegisterMemberResponse{}, err
	}
	return dto.RegisterMemberResponse{Status: "registered"}, nil
}
