package store

import (
	"context"
)

// StickerPack represents a collection of stickers.
type StickerPack struct {
	ID int32

	// Standard fields
	CreatedTs int64
	UpdatedTs int64

	// Domain specific fields
	DisplayName string
	Description string
	Enabled     bool
}

// Sticker represents an individual sticker image.
type Sticker struct {
	ID int32

	// Standard fields
	CreatedTs int64

	// Domain specific fields
	PackID       int32
	DisplayName  string
	Content      []byte
	ExternalLink string
	Type         string
	Size         int64
	Tags         []string
}

type FindStickerPack struct {
	ID      *int32
	Enabled *bool
}

type UpdateStickerPack struct {
	ID          int32
	DisplayName *string
	Description *string
	Enabled     *bool
	UpdatedTs   *int64
}

type DeleteStickerPack struct {
	ID int32
}

type FindSticker struct {
	ID     *int32
	PackID *int32
	Limit  *int
	Offset *int
}

type DeleteSticker struct {
	ID int32
}

func (s *Store) CreateStickerPack(ctx context.Context, create *StickerPack) (*StickerPack, error) {
	return s.driver.CreateStickerPack(ctx, create)
}

func (s *Store) ListStickerPacks(ctx context.Context, find *FindStickerPack) ([]*StickerPack, error) {
	return s.driver.ListStickerPacks(ctx, find)
}

func (s *Store) GetStickerPack(ctx context.Context, find *FindStickerPack) (*StickerPack, error) {
	packs, err := s.ListStickerPacks(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(packs) == 0 {
		return nil, nil
	}
	return packs[0], nil
}

func (s *Store) UpdateStickerPack(ctx context.Context, update *UpdateStickerPack) (*StickerPack, error) {
	return s.driver.UpdateStickerPack(ctx, update)
}

func (s *Store) DeleteStickerPack(ctx context.Context, delete *DeleteStickerPack) error {
	return s.driver.DeleteStickerPack(ctx, delete)
}

func (s *Store) CreateSticker(ctx context.Context, create *Sticker) (*Sticker, error) {
	return s.driver.CreateSticker(ctx, create)
}

func (s *Store) ListStickers(ctx context.Context, find *FindSticker) ([]*Sticker, error) {
	return s.driver.ListStickers(ctx, find)
}

func (s *Store) DeleteSticker(ctx context.Context, delete *DeleteSticker) error {
	return s.driver.DeleteSticker(ctx, delete)
}
