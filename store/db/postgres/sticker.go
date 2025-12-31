package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/usememos/memos/store"
)

func (d *DB) CreateStickerPack(ctx context.Context, create *store.StickerPack) (*store.StickerPack, error) {
	stmt := `
		INSERT INTO sticker_pack (display_name, description, enabled)
		VALUES ($1, $2, $3)
		RETURNING id, created_ts, updated_ts
	`
	if err := d.db.QueryRowContext(ctx, stmt, create.DisplayName, create.Description, create.Enabled).Scan(
		&create.ID,
		&create.CreatedTs,
		&create.UpdatedTs,
	); err != nil {
		return nil, err
	}

	return create, nil
}

func (d *DB) ListStickerPacks(ctx context.Context, find *store.FindStickerPack) ([]*store.StickerPack, error) {
	where, args := []string{"1 = 1"}, []any{}
	argIndex := 1

	if v := find.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.Enabled; v != nil {
		where = append(where, fmt.Sprintf("enabled = $%d", argIndex))
		args = append(args, *v)
		argIndex++
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
	set, args := []string{}, []any{}
	argIndex := 1

	if v := update.DisplayName; v != nil {
		set = append(set, fmt.Sprintf("display_name = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := update.Description; v != nil {
		set = append(set, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := update.Enabled; v != nil {
		set = append(set, fmt.Sprintf("enabled = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := update.UpdatedTs; v != nil {
		set = append(set, fmt.Sprintf("updated_ts = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}

	args = append(args, update.ID)

	stmt := `
		UPDATE sticker_pack
		SET ` + strings.Join(set, ", ") + `
		WHERE id = $` + fmt.Sprintf("%d", argIndex) + `
		RETURNING id, display_name, description, enabled, created_ts, updated_ts
	`

	pack := &store.StickerPack{}
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&pack.ID,
		&pack.DisplayName,
		&pack.Description,
		&pack.Enabled,
		&pack.CreatedTs,
		&pack.UpdatedTs,
	); err != nil {
		return nil, err
	}

	return pack, nil
}

func (d *DB) DeleteStickerPack(ctx context.Context, delete *store.DeleteStickerPack) error {
	stmt := `DELETE FROM sticker_pack WHERE id = $1`
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
		INSERT INTO sticker (pack_id, display_name, content, external_link, type, size, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_ts
	`
	if err := d.db.QueryRowContext(ctx, stmt,
		create.PackID,
		create.DisplayName,
		create.Content,
		create.ExternalLink,
		create.Type,
		create.Size,
		tagsStr,
	).Scan(
		&create.ID,
		&create.CreatedTs,
	); err != nil {
		return nil, err
	}

	return create, nil
}

func (d *DB) ListStickers(ctx context.Context, find *store.FindSticker) ([]*store.Sticker, error) {
	where, args := []string{"1 = 1"}, []any{}
	argIndex := 1

	if v := find.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.PackID; v != nil {
		where = append(where, fmt.Sprintf("pack_id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
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
	stmt := `DELETE FROM sticker WHERE id = $1`
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
