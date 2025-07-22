package storage

import (
	"context"
	"encoding"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func TableQuery[E any](log *logger.Logger, db *gorm.DB, tableName string, primaryKeyColumn string, model E) (*PostgresDBQuerier[E], error) {
	schema.RegisterSerializer("text", TextSerializer{})
	querier := newPostgresDBQuerier[E](db, tableName, primaryKeyColumn)
	return &querier, nil
}

// TextSerializer string serializer
type TextSerializer struct{}

func (TextSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	// Create a new instance of the field type
	fieldValue := reflect.New(field.FieldType).Interface()

	// Check if the fieldValue implements encoding.TextUnmarshaler
	unmarshaler, ok := fieldValue.(encoding.TextUnmarshaler)
	if !ok {
		return fmt.Errorf("field type does not implement encoding.TextUnmarshaler")
	}

	// Convert dbValue to a string or []byte
	var textData []byte
	switch v := dbValue.(type) {
	case string:
		textData = []byte(v)
	case []byte:
		textData = v
	default:
		return fmt.Errorf("unsupported dbValue type: %T", dbValue)
	}

	// Use the UnmarshalText method to populate the field
	if err := unmarshaler.UnmarshalText(textData); err != nil {
		return fmt.Errorf("failed to unmarshal text: %w", err)
	}

	// Set the value back to the destination
	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(fieldValue).Elem())
	return nil
}

// Value implements serializer interface
func (TextSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	// Check if fieldValue implements encoding.TextMarshaler
	if marshaler, ok := fieldValue.(encoding.TextMarshaler); ok {
		text, err := marshaler.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal text: %w", err)
		}
		return string(text), nil // Return the text representation as a string
	}

	return nil, fmt.Errorf("fieldValue does not implement encoding.TextMarshaler")
}

type PostgresDBQuerier[E any] struct {
	*gorm.DB
	tableName        string
	primaryKeyColumn string
}

func newPostgresDBQuerier[E any](db *gorm.DB, tableName string, primaryKeyColumn string) PostgresDBQuerier[E] {
	return PostgresDBQuerier[E]{
		DB:               db,
		tableName:        tableName,
		primaryKeyColumn: primaryKeyColumn,
	}
}

type GormExtraOps struct {
	query           interface{}
	additionalWhere []interface{}
	joins           []string
}

func applyExtraOpts(tx *gorm.DB, extraOpts []GormExtraOps) *gorm.DB {
	for _, join := range extraOpts {
		for _, j := range join.joins {
			tx = tx.Joins(j)
		}
	}

	for _, whereQuery := range extraOpts {
		tx = tx.Where(whereQuery.query, whereQuery.additionalWhere...)
	}

	return tx
}

func (db *PostgresDBQuerier[E]) Count(ctx context.Context, extraOpts []GormExtraOps) (int, error) {
	var count int64
	tx := db.Table(db.tableName).WithContext(ctx)

	tx = applyExtraOpts(tx, extraOpts)

	tx.Count(&count)
	if err := tx.Error; err != nil {
		return -1, err
	}

	return int(count), nil
}

func (db *PostgresDBQuerier[E]) SelectAll(ctx context.Context, queryParams *resources.QueryParameters, extraOpts []GormExtraOps, exhaustiveRun bool, applyFunc func(elem E)) (string, error) {
	var elems []E
	tx := db.Table(db.tableName)

	offset := 0
	limit := 15

	var sortMode string
	var sortBy string

	nextBookmark := ""

	if queryParams != nil {
		if queryParams.NextBookmark == "" {
			if queryParams.PageSize > 0 {
				limit = queryParams.PageSize
			}

			if queryParams.Sort.SortMode == "" {
				sortMode = string(resources.SortModeAsc)
			} else {
				sortMode = string(queryParams.Sort.SortMode)
			}

			nextBookmark = fmt.Sprintf("off:%d;lim:%d;", limit+offset, limit)

			if queryParams.Sort.SortField != "" {
				sortBy = strings.ReplaceAll(queryParams.Sort.SortField, ".", "_")
				nextBookmark = nextBookmark + fmt.Sprintf("sortM:%s;sortB:%s;", sortMode, sortBy)
				tx = tx.Order(sortBy + " " + sortMode)
			}

			for _, filter := range queryParams.Filters {
				tx = FilterOperandToWhereClause(filter, tx)
				nextBookmark = nextBookmark + fmt.Sprintf("filter:%s-%d-%s;", base64.StdEncoding.EncodeToString([]byte(filter.Field)), filter.FilterOperation, base64.StdEncoding.EncodeToString([]byte(filter.Value)))
			}

		} else {
			nextBookmark = ""
			decodedBookmark, err := base64.RawURLEncoding.DecodeString(queryParams.NextBookmark)
			if err != nil {
				return "", fmt.Errorf("not a valid bookmark")
			}

			splits := strings.SplitSeq(string(decodedBookmark), ";")

			for splitPart := range splits {
				queryPart := strings.Split(splitPart, ":")
				switch queryPart[0] {
				case "off":
					offset, err = strconv.Atoi(queryPart[1])
					if err != nil {
						return "", fmt.Errorf("not a valid bookmark")
					}
				case "lim":
					limit, err = strconv.Atoi(queryPart[1])
					if err != nil {
						return "", fmt.Errorf("not a valid bookmark")
					}
				case "sortM":
					sortMode = queryPart[1]
					if err != nil {
						return "", fmt.Errorf("not a valid bookmark")
					}
				case "sortB":
					sortBy = strings.ReplaceAll(queryPart[1], ".", "_")
					if err != nil {
						return "", fmt.Errorf("not a valid bookmark")
					}
				case "filter":
					filter := queryPart[1]
					if err != nil {
						return "", fmt.Errorf("not a valid bookmark")
					}
					filterSplit := strings.Split(filter, "-")
					if len(filterSplit) == 3 {
						field, err := base64.StdEncoding.DecodeString(filterSplit[0])
						if err != nil {
							continue
						}
						value, err := base64.StdEncoding.DecodeString(filterSplit[2])
						if err != nil {
							continue
						}

						operand, err := strconv.Atoi(filterSplit[1])
						if err != nil {
							continue
						}

						tx = FilterOperandToWhereClause(resources.FilterOption{
							Field:           string(field),
							FilterOperation: resources.FilterOperation(operand),
							Value:           string(value),
						}, tx)

						nextBookmark = nextBookmark + fmt.Sprintf("filter:%s-%d-%s;", base64.StdEncoding.EncodeToString([]byte(field)), operand, base64.StdEncoding.EncodeToString([]byte(value)))
					}
				}
				if sortMode != "" && sortBy != "" {
					tx = tx.Order(sortBy + " " + sortMode)
				}
			}
			nextBookmark = nextBookmark + fmt.Sprintf("off:%d;lim:%d;", offset+limit, limit)
			if queryParams.Sort.SortField != "" {
				sortBy = queryParams.Sort.SortField
				nextBookmark = nextBookmark + fmt.Sprintf("sortM:%s;sortB:%s;", sortMode, sortBy)
			}
		}
	}

	tx = applyExtraOpts(tx, extraOpts)

	if offset > 0 {
		tx.Offset(offset)
	}

	if exhaustiveRun {
		res := tx.WithContext(ctx).Preload(clause.Associations).FindInBatches(&elems, limit, func(tx *gorm.DB, batch int) error {
			for _, elem := range elems {
				applyFunc(elem)
			}

			return nil
		})
		if res.Error != nil {
			return "", res.Error
		}

		return "", nil
	} else {
		tx.Offset(offset)
		tx.Limit(limit + 1)
		rs := tx.WithContext(ctx).Preload(clause.Associations).Find(&elems)

		if rs.Error != nil {
			return "", rs.Error
		}

		// Check if we got more than the requested limit
		hasMore := len(elems) > limit

		// Trim elems to the requested limit
		if hasMore {
			elems = elems[:limit] // Keep only the requested limit
		}

		for _, elem := range elems {
			// batch processing found records
			applyFunc(elem)
		}

		if !hasMore {
			// no more records to fetch. Reset nextBookmark to empty string
			return "", nil
		}

		return base64.RawURLEncoding.EncodeToString([]byte(nextBookmark)), nil
	}
}

// Selects first element from DB. if queryCol is empty or nil, the primary key column
// defined in the creation process, is used.
func (db *PostgresDBQuerier[E]) SelectExists(ctx context.Context, queryID string, queryCol *string) (bool, *E, error) {
	searchCol := db.primaryKeyColumn
	if queryCol != nil && *queryCol != "" {
		searchCol = *queryCol
	}

	var elem E
	tx := db.Table(db.tableName).WithContext(ctx).Preload(clause.Associations).Limit(1).Find(&elem, fmt.Sprintf("%s = ?", searchCol), queryID)
	if tx.Error != nil {
		return false, nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return false, nil, nil // No record found, but no error
	}

	return true, &elem, nil
}

func (db *PostgresDBQuerier[E]) Insert(ctx context.Context, elem *E) (*E, error) {
	tx := db.Table(db.tableName).WithContext(ctx).Create(elem)
	if err := tx.Error; err != nil {
		return nil, err
	}

	return elem, nil
}

func (db *PostgresDBQuerier[E]) Update(ctx context.Context, elem *E, elemID string) (*E, error) {
	tx := db.Session(&gorm.Session{FullSaveAssociations: true}).Table(db.tableName).WithContext(ctx).Where(fmt.Sprintf("%s = ?", db.primaryKeyColumn), elemID).Save(elem)
	if err := tx.Error; err != nil {
		return nil, err
	}

	if tx.RowsAffected != 1 {
		return nil, gorm.ErrRecordNotFound
	}

	return elem, nil
}

func (db *PostgresDBQuerier[E]) Delete(ctx context.Context, elemID string) error {
	tx := db.Table(db.tableName).WithContext(ctx).Delete(nil, db.Where(fmt.Sprintf("%s = ?", db.primaryKeyColumn), elemID))
	if err := tx.Error; err != nil {
		return err
	}

	if tx.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func FilterOperandToWhereClause(filter resources.FilterOption, tx *gorm.DB) *gorm.DB {
	if strings.Contains(filter.Field, ".") {
		filter.Field = strings.ReplaceAll(filter.Field, ".", "_")
	}

	switch filter.FilterOperation {
	case resources.StringEqual:
		return tx.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case resources.StringEqualIgnoreCase:
		return tx.Where(fmt.Sprintf("%s ILIKE ?", filter.Field), filter.Value)
	case resources.StringNotEqual:
		return tx.Where(fmt.Sprintf("%s <> ?", filter.Field), filter.Value)
	case resources.StringNotEqualIgnoreCase:
		return tx.Where(fmt.Sprintf("%s NOT ILIKE ?", filter.Field), filter.Value)
	case resources.StringContains:
		return tx.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%s%%", filter.Value))
	case resources.StringContainsIgnoreCase:
		return tx.Where(fmt.Sprintf("%s ILIKE ?", filter.Field), fmt.Sprintf("%%%s%%", filter.Value))
	case resources.StringArrayContains:
		// return tx.Where(fmt.Sprintf("? = ANY(%s)", filter.Field), filter.Value)
		return tx.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%s%%", filter.Value))
	case resources.StringArrayContainsIgnoreCase:
		// return tx.Where(fmt.Sprintf("? = ANY(%s)", filter.Field), filter.Value)
		return tx.Where(fmt.Sprintf("%s ILIKE ?", filter.Field), fmt.Sprintf("%%%s%%", filter.Value))
	case resources.StringNotContains:
		return tx.Where(fmt.Sprintf("%s NOT LIKE ?", filter.Field), fmt.Sprintf("%%%s%%", filter.Value))
	case resources.StringNotContainsIgnoreCase:
		return tx.Where(fmt.Sprintf("%s NOT ILIKE ?", filter.Field), fmt.Sprintf("%%%s%%", filter.Value))
	case resources.DateEqual:
		return tx.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case resources.DateBefore:
		return tx.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case resources.DateAfter:
		return tx.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case resources.NumberEqual:
		return tx.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case resources.NumberNotEqual:
		return tx.Where(fmt.Sprintf("%s <> ?", filter.Field), filter.Value)
	case resources.NumberLessThan:
		return tx.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case resources.NumberLessOrEqualThan:
		return tx.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case resources.NumberGreaterThan:
		return tx.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case resources.NumberGreaterOrEqualThan:
		return tx.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case resources.EnumEqual:
		return tx.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case resources.EnumNotEqual:
		return tx.Where(fmt.Sprintf("%s <> ?", filter.Field), filter.Value)
	default:
		return tx
	}
}
