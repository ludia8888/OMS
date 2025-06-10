package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/openfoundry/oms/internal/domain/entity"
	"github.com/openfoundry/oms/internal/domain/repository"
)

// PostgresObjectTypeRepository implements ObjectTypeRepository using PostgreSQL
type PostgresObjectTypeRepository struct {
	db *sql.DB
}

// NewPostgresObjectTypeRepository creates a new PostgreSQL repository
func NewPostgresObjectTypeRepository(db *sql.DB) repository.ObjectTypeRepository {
	return &PostgresObjectTypeRepository{db: db}
}

// Create creates a new object type
func (r *PostgresObjectTypeRepository) Create(ctx context.Context, objectType *entity.ObjectType) error {
	// Serialize properties and metadata to JSON
	propertiesJSON, err := json.Marshal(objectType.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	metadataJSON, err := json.Marshal(objectType.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	baseDatasetsJSON, err := json.Marshal(objectType.BaseDatasets)
	if err != nil {
		return fmt.Errorf("failed to marshal base datasets: %w", err)
	}

	// Insert object type
	query := `
		INSERT INTO object_types (
			id, name, display_name, description, category, tags,
			properties, base_datasets, metadata, version, is_deleted,
			created_at, created_by, updated_at, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	_, err = r.db.ExecContext(ctx, query,
		objectType.ID,
		objectType.Name,
		objectType.DisplayName,
		objectType.Description,
		objectType.Category,
		pq.Array(objectType.Tags),
		propertiesJSON,
		baseDatasetsJSON,
		metadataJSON,
		objectType.Version,
		objectType.IsDeleted,
		objectType.CreatedAt,
		objectType.CreatedBy,
		objectType.UpdatedAt,
		objectType.UpdatedBy,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return entity.ErrObjectTypeNameExists
			}
		}
		return fmt.Errorf("failed to create object type: %w", err)
	}

	// Create initial version record
	if err := r.createVersion(ctx, objectType); err != nil {
		return fmt.Errorf("failed to create version record: %w", err)
	}

	return nil
}

// GetByID retrieves an object type by ID
func (r *PostgresObjectTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ObjectType, error) {
	query := `
		SELECT id, name, display_name, description, category, tags,
			   properties, base_datasets, metadata, version,
			   created_at, created_by, updated_at, updated_by
		FROM object_types
		WHERE id = $1 AND is_deleted = FALSE`

	return r.scanObjectType(r.db.QueryRowContext(ctx, query, id))
}

// GetByName retrieves an object type by name
func (r *PostgresObjectTypeRepository) GetByName(ctx context.Context, name string) (*entity.ObjectType, error) {
	query := `
		SELECT id, name, display_name, description, category, tags,
			   properties, base_datasets, metadata, version,
			   created_at, created_by, updated_at, updated_by
		FROM object_types
		WHERE name = $1 AND is_deleted = FALSE`

	return r.scanObjectType(r.db.QueryRowContext(ctx, query, name))
}

// Update updates an existing object type
func (r *PostgresObjectTypeRepository) Update(ctx context.Context, objectType *entity.ObjectType) error {
	// Serialize properties and metadata to JSON
	propertiesJSON, err := json.Marshal(objectType.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	metadataJSON, err := json.Marshal(objectType.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	baseDatasetsJSON, err := json.Marshal(objectType.BaseDatasets)
	if err != nil {
		return fmt.Errorf("failed to marshal base datasets: %w", err)
	}

	// Update object type
	query := `
		UPDATE object_types SET
			display_name = $2,
			description = $3,
			category = $4,
			tags = $5,
			properties = $6,
			base_datasets = $7,
			metadata = $8,
			version = $9,
			updated_at = $10,
			updated_by = $11
		WHERE id = $1 AND is_deleted = FALSE`

	result, err := r.db.ExecContext(ctx, query,
		objectType.ID,
		objectType.DisplayName,
		objectType.Description,
		objectType.Category,
		pq.Array(objectType.Tags),
		propertiesJSON,
		baseDatasetsJSON,
		metadataJSON,
		objectType.Version,
		objectType.UpdatedAt,
		objectType.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update object type: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrObjectTypeNotFound
	}

	// Create version record
	if err := r.createVersion(ctx, objectType); err != nil {
		return fmt.Errorf("failed to create version record: %w", err)
	}

	return nil
}

// Delete soft deletes an object type
func (r *PostgresObjectTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE object_types 
		SET is_deleted = TRUE, updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete object type: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrObjectTypeNotFound
	}

	return nil
}

// List retrieves a list of object types based on filter
func (r *PostgresObjectTypeRepository) List(ctx context.Context, filter repository.ObjectTypeFilter) ([]*entity.ObjectType, error) {
	query := `
		SELECT id, name, display_name, description, category, tags,
			   properties, base_datasets, metadata, version,
			   created_at, created_by, updated_at, updated_by
		FROM object_types
		WHERE is_deleted = FALSE`

	var args []interface{}
	argCount := 0

	// Handle cursor-based pagination
	if filter.PageCursor != "" {
		cursor, err := r.decodeCursor(filter.PageCursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		argCount++
		query += fmt.Sprintf(" AND (created_at, id) < ($%d, $%d)", argCount, argCount+1)
		args = append(args, cursor.Timestamp, cursor.ID)
		argCount++
	}

	// Apply filters
	if filter.Category != nil {
		argCount++
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, *filter.Category)
	}

	if len(filter.Tags) > 0 {
		argCount++
		query += fmt.Sprintf(" AND tags && $%d", argCount)
		args = append(args, pq.Array(filter.Tags))
	}

	// Order and limit
	query += " ORDER BY created_at DESC, id DESC"
	if filter.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.PageSize)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list object types: %w", err)
	}
	defer rows.Close()

	var objectTypes []*entity.ObjectType
	for rows.Next() {
		ot, err := r.scanObjectTypeFromRows(rows)
		if err != nil {
			return nil, err
		}
		objectTypes = append(objectTypes, ot)
	}

	return objectTypes, rows.Err()
}

// Count counts object types based on filter
func (r *PostgresObjectTypeRepository) Count(ctx context.Context, filter repository.ObjectTypeFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM object_types WHERE is_deleted = FALSE`

	var args []interface{}
	argCount := 0

	// Apply filters
	if filter.Category != nil {
		argCount++
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, *filter.Category)
	}

	if len(filter.Tags) > 0 {
		argCount++
		query += fmt.Sprintf(" AND tags && $%d", argCount)
		args = append(args, pq.Array(filter.Tags))
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count object types: %w", err)
	}

	return count, nil
}

// Search implements full-text search using PostgreSQL's tsvector
func (r *PostgresObjectTypeRepository) Search(ctx context.Context, query string, limit int) ([]*entity.ObjectType, error) {
	sql := `
		SELECT id, name, display_name, description, category, tags,
			   properties, base_datasets, metadata, version,
			   created_at, created_by, updated_at, updated_by
		FROM object_types 
		WHERE to_tsvector('english', name || ' ' || display_name || ' ' || COALESCE(description, '')) 
		@@ plainto_tsquery('english', $1)
		AND is_deleted = FALSE
		ORDER BY ts_rank(to_tsvector('english', name || ' ' || display_name || ' ' || COALESCE(description, '')), 
						plainto_tsquery('english', $1)) DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, sql, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search object types: %w", err)
	}
	defer rows.Close()

	var results []*entity.ObjectType
	for rows.Next() {
		ot, err := r.scanObjectTypeFromRows(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, ot)
	}

	return results, rows.Err()
}

// GetVersion retrieves a specific version of an object type
func (r *PostgresObjectTypeRepository) GetVersion(ctx context.Context, id uuid.UUID, version int) (*entity.ObjectType, error) {
	query := `
		SELECT snapshot
		FROM object_type_versions
		WHERE object_type_id = $1 AND version = $2`

	var snapshotJSON []byte
	err := r.db.QueryRowContext(ctx, query, id, version).Scan(&snapshotJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrObjectTypeNotFound
		}
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	var objectType entity.ObjectType
	if err := json.Unmarshal(snapshotJSON, &objectType); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &objectType, nil
}

// ListVersions lists all versions of an object type
func (r *PostgresObjectTypeRepository) ListVersions(ctx context.Context, id uuid.UUID) ([]*repository.ObjectTypeVersion, error) {
	query := `
		SELECT id, object_type_id, version, snapshot, change_description, created_at, created_by
		FROM object_type_versions
		WHERE object_type_id = $1
		ORDER BY version DESC`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}
	defer rows.Close()

	var versions []*repository.ObjectTypeVersion
	for rows.Next() {
		var v repository.ObjectTypeVersion
		var snapshotJSON []byte

		err := rows.Scan(
			&v.ID,
			&v.ObjectTypeID,
			&v.Version,
			&snapshotJSON,
			&v.ChangeDescription,
			&v.CreatedAt,
			&v.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}

		if err := json.Unmarshal(snapshotJSON, &v.Snapshot); err != nil {
			return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
		}

		versions = append(versions, &v)
	}

	return versions, rows.Err()
}

// CompareVersions compares two versions of an object type
func (r *PostgresObjectTypeRepository) CompareVersions(ctx context.Context, id uuid.UUID, v1, v2 int) (*repository.VersionDiff, error) {
	// Get both versions
	version1, err := r.GetVersion(ctx, id, v1)
	if err != nil {
		return nil, fmt.Errorf("failed to get version %d: %w", v1, err)
	}

	version2, err := r.GetVersion(ctx, id, v2)
	if err != nil {
		return nil, fmt.Errorf("failed to get version %d: %w", v2, err)
	}

	// Compare versions
	diff := &repository.VersionDiff{
		ObjectTypeID: id,
		Version1:     v1,
		Version2:     v2,
		Changes:      []repository.FieldChange{},
	}

	// Compare basic fields
	if version1.Name != version2.Name {
		diff.Changes = append(diff.Changes, repository.FieldChange{
			Field:    "name",
			OldValue: version1.Name,
			NewValue: version2.Name,
			Type:     repository.ChangeTypeModified,
		})
	}

	if version1.DisplayName != version2.DisplayName {
		diff.Changes = append(diff.Changes, repository.FieldChange{
			Field:    "displayName",
			OldValue: version1.DisplayName,
			NewValue: version2.DisplayName,
			Type:     repository.ChangeTypeModified,
		})
	}

	// Compare properties
	propChanges := r.compareProperties(version1.Properties, version2.Properties)
	diff.Changes = append(diff.Changes, propChanges...)

	return diff, nil
}

// BatchCreate creates multiple object types
func (r *PostgresObjectTypeRepository) BatchCreate(ctx context.Context, objectTypes []*entity.ObjectType) error {
	// Use transaction for batch operation
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO object_types (
			id, name, display_name, description, category, tags,
			properties, base_datasets, metadata, version, is_deleted,
			created_at, created_by, updated_at, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, ot := range objectTypes {
		propertiesJSON, _ := json.Marshal(ot.Properties)
		metadataJSON, _ := json.Marshal(ot.Metadata)
		baseDatasetsJSON, _ := json.Marshal(ot.BaseDatasets)

		_, err := stmt.ExecContext(ctx,
			ot.ID, ot.Name, ot.DisplayName, ot.Description, ot.Category,
			pq.Array(ot.Tags), propertiesJSON, baseDatasetsJSON, metadataJSON,
			ot.Version, ot.IsDeleted, ot.CreatedAt, ot.CreatedBy,
			ot.UpdatedAt, ot.UpdatedBy,
		)
		if err != nil {
			return fmt.Errorf("failed to insert object type %s: %w", ot.Name, err)
		}

		// Create version record
		if err := r.createVersionTx(ctx, tx, ot); err != nil {
			return fmt.Errorf("failed to create version for %s: %w", ot.Name, err)
		}
	}

	return tx.Commit()
}

// BatchUpdate updates multiple object types
func (r *PostgresObjectTypeRepository) BatchUpdate(ctx context.Context, objectTypes []*entity.ObjectType) error {
	// Use transaction for batch operation
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE object_types SET
			display_name = $2,
			description = $3,
			category = $4,
			tags = $5,
			properties = $6,
			base_datasets = $7,
			metadata = $8,
			version = $9,
			updated_at = $10,
			updated_by = $11
		WHERE id = $1 AND is_deleted = FALSE`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, ot := range objectTypes {
		propertiesJSON, _ := json.Marshal(ot.Properties)
		metadataJSON, _ := json.Marshal(ot.Metadata)
		baseDatasetsJSON, _ := json.Marshal(ot.BaseDatasets)

		_, err := stmt.ExecContext(ctx,
			ot.ID, ot.DisplayName, ot.Description, ot.Category,
			pq.Array(ot.Tags), propertiesJSON, baseDatasetsJSON, metadataJSON,
			ot.Version, ot.UpdatedAt, ot.UpdatedBy,
		)
		if err != nil {
			return fmt.Errorf("failed to update object type %s: %w", ot.Name, err)
		}

		// Create version record
		if err := r.createVersionTx(ctx, tx, ot); err != nil {
			return fmt.Errorf("failed to create version for %s: %w", ot.Name, err)
		}
	}

	return tx.Commit()
}

// Helper methods

func (r *PostgresObjectTypeRepository) scanObjectType(row *sql.Row) (*entity.ObjectType, error) {
	var ot entity.ObjectType
	var propertiesJSON, baseDatasetsJSON, metadataJSON []byte

	err := row.Scan(
		&ot.ID,
		&ot.Name,
		&ot.DisplayName,
		&ot.Description,
		&ot.Category,
		pq.Array(&ot.Tags),
		&propertiesJSON,
		&baseDatasetsJSON,
		&metadataJSON,
		&ot.Version,
		&ot.CreatedAt,
		&ot.CreatedBy,
		&ot.UpdatedAt,
		&ot.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrObjectTypeNotFound
		}
		return nil, fmt.Errorf("failed to scan object type: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(propertiesJSON, &ot.Properties); err != nil {
		return nil, fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	if err := json.Unmarshal(baseDatasetsJSON, &ot.BaseDatasets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base datasets: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &ot.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &ot, nil
}

func (r *PostgresObjectTypeRepository) scanObjectTypeFromRows(rows *sql.Rows) (*entity.ObjectType, error) {
	var ot entity.ObjectType
	var propertiesJSON, baseDatasetsJSON, metadataJSON []byte

	err := rows.Scan(
		&ot.ID,
		&ot.Name,
		&ot.DisplayName,
		&ot.Description,
		&ot.Category,
		pq.Array(&ot.Tags),
		&propertiesJSON,
		&baseDatasetsJSON,
		&metadataJSON,
		&ot.Version,
		&ot.CreatedAt,
		&ot.CreatedBy,
		&ot.UpdatedAt,
		&ot.UpdatedBy,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan object type: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(propertiesJSON, &ot.Properties); err != nil {
		return nil, fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	if err := json.Unmarshal(baseDatasetsJSON, &ot.BaseDatasets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base datasets: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &ot.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &ot, nil
}

func (r *PostgresObjectTypeRepository) createVersion(ctx context.Context, objectType *entity.ObjectType) error {
	return r.createVersionTx(ctx, r.db, objectType)
}

func (r *PostgresObjectTypeRepository) createVersionTx(ctx context.Context, tx interface{ ExecContext(context.Context, string, ...interface{}) (sql.Result, error) }, objectType *entity.ObjectType) error {
	snapshotJSON, err := json.Marshal(objectType)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	query := `
		INSERT INTO object_type_versions (
			object_type_id, version, snapshot, created_at, created_by
		) VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.ExecContext(ctx, query,
		objectType.ID,
		objectType.Version,
		snapshotJSON,
		objectType.UpdatedAt,
		objectType.UpdatedBy,
	)

	return err
}

func (r *PostgresObjectTypeRepository) compareProperties(props1, props2 []entity.Property) []repository.FieldChange {
	var changes []repository.FieldChange

	// Create maps for easier comparison
	props1Map := make(map[string]entity.Property)
	props2Map := make(map[string]entity.Property)

	for _, p := range props1 {
		props1Map[p.Name] = p
	}
	for _, p := range props2 {
		props2Map[p.Name] = p
	}

	// Check for removed and modified properties
	for name, p1 := range props1Map {
		if p2, exists := props2Map[name]; exists {
			// Check if property was modified
			if p1.DataType != p2.DataType || p1.Required != p2.Required {
				changes = append(changes, repository.FieldChange{
					Field:    fmt.Sprintf("properties.%s", name),
					OldValue: p1,
					NewValue: p2,
					Type:     repository.ChangeTypeModified,
				})
			}
		} else {
			// Property was removed
			changes = append(changes, repository.FieldChange{
				Field:    fmt.Sprintf("properties.%s", name),
				OldValue: p1,
				NewValue: nil,
				Type:     repository.ChangeTypeRemoved,
			})
		}
	}

	// Check for added properties
	for name, p2 := range props2Map {
		if _, exists := props1Map[name]; !exists {
			changes = append(changes, repository.FieldChange{
				Field:    fmt.Sprintf("properties.%s", name),
				OldValue: nil,
				NewValue: p2,
				Type:     repository.ChangeTypeAdded,
			})
		}
	}

	return changes
}

func (r *PostgresObjectTypeRepository) encodeCursor(timestamp time.Time, id uuid.UUID) string {
	data := fmt.Sprintf("%d:%s", timestamp.Unix(), id.String())
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (r *PostgresObjectTypeRepository) decodeCursor(cursor string) (*repository.PageCursor, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(data), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cursor format")
	}

	timestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, err
	}

	return &repository.PageCursor{
		Timestamp: time.Unix(timestamp, 0),
		ID:        id,
	}, nil
}