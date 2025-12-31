package v1

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/usememos/memos/internal/util"
	v1pb "github.com/usememos/memos/proto/gen/api/v1"
	"github.com/usememos/memos/store"
)

// Helper function to extract sticker pack ID from resource name.
// Format: sticker-packs/{sticker_pack}.
func extractStickerPackIDFromName(name string) (int32, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 2 || parts[0] != "sticker-packs" {
		return 0, errors.Errorf("invalid sticker pack name format: %s", name)
	}

	id, err := util.ConvertStringToInt32(parts[1])
	if err != nil {
		return 0, errors.Errorf("invalid sticker pack ID %q", parts[1])
	}

	return id, nil
}

// Helper function to extract sticker pack ID and sticker ID from resource name.
// Format: sticker-packs/{sticker_pack}/stickers/{sticker}.
func extractStickerPackAndStickerIDFromName(name string) (int32, int32, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 4 || parts[0] != "sticker-packs" || parts[2] != "stickers" {
		return 0, 0, errors.Errorf("invalid sticker name format: %s", name)
	}

	packID, err := util.ConvertStringToInt32(parts[1])
	if err != nil {
		return 0, 0, errors.Errorf("invalid sticker pack ID %q", parts[1])
	}

	stickerID, err := util.ConvertStringToInt32(parts[3])
	if err != nil {
		return 0, 0, errors.Errorf("invalid sticker ID %q", parts[3])
	}

	return packID, stickerID, nil
}

func (s *APIV1Service) ListStickerPacks(ctx context.Context, _ *v1pb.ListStickerPacksRequest) (*v1pb.ListStickerPacksResponse, error) {
	stickerPacks, err := s.Store.ListStickerPacks(ctx, &store.FindStickerPack{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list sticker packs: %v", err)
	}

	response := &v1pb.ListStickerPacksResponse{}
	for _, pack := range stickerPacks {
		packPb, err := s.convertStickerPackFromStore(ctx, pack)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert sticker pack: %v", err)
		}
		response.StickerPacks = append(response.StickerPacks, packPb)
	}

	return response, nil
}

func (s *APIV1Service) GetStickerPack(ctx context.Context, request *v1pb.GetStickerPackRequest) (*v1pb.StickerPack, error) {
	packID, err := extractStickerPackIDFromName(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sticker pack name: %v", err)
	}

	pack, err := s.Store.GetStickerPack(ctx, &store.FindStickerPack{
		ID: &packID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get sticker pack: %v", err)
	}
	if pack == nil {
		return nil, status.Errorf(codes.NotFound, "sticker pack not found")
	}

	return s.convertStickerPackFromStore(ctx, pack)
}

func (s *APIV1Service) CreateStickerPack(ctx context.Context, request *v1pb.CreateStickerPackRequest) (*v1pb.StickerPack, error) {
	currentUser, err := s.fetchCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	// Only admins can create sticker packs
	if currentUser.Role != store.RoleAdmin && currentUser.Role != store.RoleHost {
		return nil, status.Errorf(codes.PermissionDenied, "only admins can create sticker packs")
	}

	pack := &store.StickerPack{
		DisplayName: request.StickerPack.DisplayName,
		Description: request.StickerPack.Description,
		Enabled:     request.StickerPack.Enabled,
	}

	createdPack, err := s.Store.CreateStickerPack(ctx, pack)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create sticker pack: %v", err)
	}

	return s.convertStickerPackFromStore(ctx, createdPack)
}

func (s *APIV1Service) UpdateStickerPack(ctx context.Context, request *v1pb.UpdateStickerPackRequest) (*v1pb.StickerPack, error) {
	currentUser, err := s.fetchCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	// Only admins can update sticker packs
	if currentUser.Role != store.RoleAdmin && currentUser.Role != store.RoleHost {
		return nil, status.Errorf(codes.PermissionDenied, "only admins can update sticker packs")
	}

	packID, err := extractStickerPackIDFromName(request.StickerPack.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sticker pack name: %v", err)
	}

	update := &store.UpdateStickerPack{
		ID: packID,
	}

	if request.UpdateMask != nil {
		for _, field := range request.UpdateMask.Paths {
			switch field {
			case "display_name":
				update.DisplayName = &request.StickerPack.DisplayName
			case "description":
				update.Description = &request.StickerPack.Description
			case "enabled":
				update.Enabled = &request.StickerPack.Enabled
			}
		}
	}

	updatedPack, err := s.Store.UpdateStickerPack(ctx, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update sticker pack: %v", err)
	}

	return s.convertStickerPackFromStore(ctx, updatedPack)
}

func (s *APIV1Service) DeleteStickerPack(ctx context.Context, request *v1pb.DeleteStickerPackRequest) (*emptypb.Empty, error) {
	currentUser, err := s.fetchCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	// Only admins can delete sticker packs
	if currentUser.Role != store.RoleAdmin && currentUser.Role != store.RoleHost {
		return nil, status.Errorf(codes.PermissionDenied, "only admins can delete sticker packs")
	}

	packID, err := extractStickerPackIDFromName(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sticker pack name: %v", err)
	}

	if err := s.Store.DeleteStickerPack(ctx, &store.DeleteStickerPack{
		ID: packID,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete sticker pack: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) CreateSticker(ctx context.Context, request *v1pb.CreateStickerRequest) (*v1pb.Sticker, error) {
	currentUser, err := s.fetchCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	// Only admins can create stickers
	if currentUser.Role != store.RoleAdmin && currentUser.Role != store.RoleHost {
		return nil, status.Errorf(codes.PermissionDenied, "only admins can create stickers")
	}

	packID, err := extractStickerPackIDFromName(request.Parent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent name: %v", err)
	}

	sticker := &store.Sticker{
		PackID:       packID,
		DisplayName:  request.Sticker.DisplayName,
		Content:      request.Sticker.Content,
		ExternalLink: request.Sticker.ExternalLink,
		Type:         request.Sticker.Type,
		Tags:         request.Sticker.Tags,
	}

	createdSticker, err := s.Store.CreateSticker(ctx, sticker)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create sticker: %v", err)
	}

	return s.convertStickerFromStore(createdSticker), nil
}

func (s *APIV1Service) DeleteSticker(ctx context.Context, request *v1pb.DeleteStickerRequest) (*emptypb.Empty, error) {
	currentUser, err := s.fetchCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	// Only admins can delete stickers
	if currentUser.Role != store.RoleAdmin && currentUser.Role != store.RoleHost {
		return nil, status.Errorf(codes.PermissionDenied, "only admins can delete stickers")
	}

	_, stickerID, err := extractStickerPackAndStickerIDFromName(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sticker name: %v", err)
	}

	if err := s.Store.DeleteSticker(ctx, &store.DeleteSticker{
		ID: stickerID,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete sticker: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) convertStickerPackFromStore(ctx context.Context, pack *store.StickerPack) (*v1pb.StickerPack, error) {
	packPb := &v1pb.StickerPack{
		Name:        fmt.Sprintf("sticker-packs/%d", pack.ID),
		DisplayName: pack.DisplayName,
		Description: pack.Description,
		Enabled:     pack.Enabled,
		CreateTime:  timestamppb.New(time.Unix(pack.CreatedTs, 0)),
		UpdateTime:  timestamppb.New(time.Unix(pack.UpdatedTs, 0)),
	}

	// Load stickers for this pack
	stickers, err := s.Store.ListStickers(ctx, &store.FindSticker{
		PackID: &pack.ID,
	})
	if err != nil {
		return nil, err
	}

	for _, sticker := range stickers {
		packPb.Stickers = append(packPb.Stickers, s.convertStickerFromStore(sticker))
	}

	return packPb, nil
}

func (s *APIV1Service) convertStickerFromStore(sticker *store.Sticker) *v1pb.Sticker {
	return &v1pb.Sticker{
		Name:         fmt.Sprintf("sticker-packs/%d/stickers/%d", sticker.PackID, sticker.ID),
		DisplayName:  sticker.DisplayName,
		Content:      sticker.Content,
		ExternalLink: sticker.ExternalLink,
		Type:         sticker.Type,
		Size:         sticker.Size,
		Tags:         sticker.Tags,
		CreateTime:   timestamppb.New(time.Unix(sticker.CreatedTs, 0)),
	}
}
