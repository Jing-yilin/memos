package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/usememos/memos/store"
)

func (d *DB) CreateStickerPack(ctx context.Context, create *store.StickerPack) (*store.StickerPack, error) {
	stmt := `
		INSERT INTO sticker_pack (display_name, description, enabled, created_ts, updated_ts)
		VALUES (?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
	`
	result, err := d.db.ExecContext(ctx, stmt, create.DisplayName, create.Description, create.Enabled)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	create.ID = int32(id)
	pack, err := d.GetStickerPack(ctx, &store.FindStickerPack{ID: &create.ID})
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (d *DB) GetStickerPack(ctx context.Context, find *store.FindStickerPack) (*store.StickerPack, error) {
	packs, err := d.ListStickerPacks(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(packs) == 0 {
		return nil, nil
	}
	return packs[0], nil
}

func (d *DB) ListStickerPacks(ctx context.Context, find *store.FindStickerPack) ([]*store.StickerPack, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.Enabled; v != nil {
		where, args = append(where, "enabled = ?"), append(args, *v)
	}

	query := `
		SELECT id, display_name, description, enabled, created_ts, updated_ts
		FROM sticker_pack
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY created_ts DESC
	`

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packs := []*store.StickerPack{}
	for rows.Next() {
		pack := &store.StickerPack{}
		if err := rows.Scan(
			&pack.ID,
			&pack.DisplayName,
			&pack.Description,
			&pack.Enabled,
			&pack.CreatedTs,
			&pack.UpdatedTs,
		); err != nil {
			return nil, err
		}
		packs = append(packs, pack)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return packs, nil
}

func (d *DB) UpdateStickerPack(ctx context.Context, update *store.UpdateStickerPack) (*store.StickerPack, error) {
	set, args := []string{"updated_ts = UNIX_TIMESTAMP()"}, []any{}

	if v := update.DisplayName; v != nil {
		set, args = append(set, "display_name = ?"), append(args, *v)
	}
	if v := update.Description; v != nil {
		set, args = append(set, "description = ?"), append(args, *v)
	}
	if v := update.Enabled; v != nil {
		set, args = append(set, "enabled = ?"), append(args, *v)
	}

	args = append(args, update.ID)

	stmt := `UPDATE sticker_pack SET ` + strings.Join(set, ", ") + ` WHERE id = ?`
	if _, err := d.db.ExecContext(ctx, stmt, args...); err != nil {
		return nil, err
	}

	pack, err := d.GetStickerPack(ctx, &store.FindStickerPack{ID: &update.ID})
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (d *DB) DeleteStickerPack(ctx context.Context, delete *store.DeleteStickerPack) error {
	stmt := `DELETE FROM sticker_pack WHERE id = ?`
	result, err := d.db.ExecContext(ctx, stmt, delete.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *DB) CreateSticker(ctx context.Context, create *store.Sticker) (*store.Sticker, error) {
	tagsStr := strings.Join(create.Tags, ",")

	stmt := `
		INSERT INTO sticker (pack_id, display_name, content, external_link, type, size, tags, created_ts)
		VALUES (?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP())
	`
	result, err := d.db.ExecContext(ctx, stmt,
		create.PackID,
		create.DisplayName,
		create.Content,
		create.ExternalLink,
		create.Type,
		create.Size,
		tagsStr,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	create.ID = int32(id)
	stickers, err := d.ListStickers(ctx, &store.FindSticker{ID: &create.ID})
	if err != nil {
		return nil, err
	}
	if len(stickers) == 0 {
		return nil, sql.ErrNoRows
	}

	return stickers[0], nil
}

func (d *DB) ListStickers(ctx context.Context, find *store.FindSticker) ([]*store.Sticker, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.PackID; v != nil {
		where, args = append(where, "pack_id = ?"), append(args, *v)
	}

	query := `
		SELECT id, pack_id, display_name, content, external_link, type, size, tags, created_ts
		FROM sticker
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY created_ts ASC
	`

	if find.Limit != nil {
		query = fmt.Sprintf("%s LIMIT %d", query, *find.Limit)
		if find.Offset != nil {
			query = fmt.Sprintf("%s OFFSET %d", query, *find.Offset)
		}
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stickers := []*store.Sticker{}
	for rows.Next() {
		sticker := &store.Sticker{}
		var tagsStr string

		if err := rows.Scan(
			&sticker.ID,
			&sticker.PackID,
			&sticker.DisplayName,
			&sticker.Content,
			&sticker.ExternalLink,
			&sticker.Type,
			&sticker.Size,
			&tagsStr,
			&sticker.CreatedTs,
		); err != nil {
			return nil, err
		}

		if tagsStr != "" {
			sticker.Tags = strings.Split(tagsStr, ",")
		} else {
			sticker.Tags = []string{}
		}

		stickers = append(stickers, sticker)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stickers, nil
}

func (d *DB) DeleteSticker(ctx context.Context, delete *store.DeleteSticker) error {
	stmt := `DELETE FROM sticker WHERE id = ?`
	result, err := d.db.ExecContext(ctx, stmt, delete.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
