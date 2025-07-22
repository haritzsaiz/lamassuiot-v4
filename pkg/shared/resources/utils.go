package resources

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type StorageListRequest[E any] struct {
	ExhaustiveRun bool
	ApplyFunc     func(E)
	QueryParams   *QueryParameters
	ExtraOpts     map[string]interface{}
}

func FilterQuery(f *fiber.Ctx, filterFieldMap map[string]FilterFieldType) *QueryParameters {
	queryParams := QueryParameters{
		NextBookmark: "",
		Filters:      []FilterOption{},
		PageSize:     25,
	}

	f.Context().QueryArgs().VisitAll(func(k, value []byte) {
		v := []string{string(value)}
		switch string(k) {
		case "sort_by":
			value := v[len(v)-1] //only get last
			sortQueryParam := value
			sortField := strings.Trim(sortQueryParam, " ")

			_, exists := filterFieldMap[sortField]
			if exists {
				queryParams.Sort.SortField = sortField
			}

		case "sort_mode":
			value := v[len(v)-1] //only get last
			sortQueryParam := value
			sortMode := SortModeAsc

			if sortQueryParam == "desc" {
				sortMode = SortModeDesc
			}

			queryParams.Sort.SortMode = sortMode

		case "page_size":
			value := v[len(v)-1] //only get last
			pageS, err := strconv.Atoi(value)
			if err == nil {
				queryParams.PageSize = pageS
			}

		case "bookmark":
			value := v[len(v)-1] //only get last
			queryParams.NextBookmark = value

		case "filter":
			for _, value := range v {
				bs := strings.Index(value, "[")
				es := strings.Index(value, "]")
				if bs != -1 && es != -1 && bs < es {
					field, rest, _ := strings.Cut(value, "[")
					operand, arg, _ := strings.Cut(rest, "]")
					operand = strings.ToLower(operand)

					fieldOperandType, exists := filterFieldMap[field]
					if !exists {
						continue
					}

					var filterOperand FilterOperation
					switch fieldOperandType {
					case StringFilterFieldType:
						switch operand {
						case "eq", "equal":
							filterOperand = StringEqual
						case "eq_ic", "equal_ignorecase":
							filterOperand = StringEqualIgnoreCase
						case "ne", "notequal":
							filterOperand = StringNotEqual
						case "ne_ic", "notequal_ignorecase":
							filterOperand = StringNotEqualIgnoreCase
						case "ct", "contains":
							filterOperand = StringContains
						case "ct_ic", "contains_ignorecase":
							filterOperand = StringContainsIgnoreCase
						case "nc", "notcontains":
							filterOperand = StringNotContains
						case "nc_ic", "notcontains_ignorecase":
							filterOperand = StringNotContainsIgnoreCase
						}

					case StringArrayFilterFieldType:
						if strings.Contains(operand, "ignorecase") {
							filterOperand = StringArrayContainsIgnoreCase
						} else {
							filterOperand = StringArrayContains
						}

					case DateFilterFieldType:
						switch operand {
						case "bf", "before":
							filterOperand = DateBefore
						case "eq", "equal":
							filterOperand = DateEqual
						case "af", "after":
							filterOperand = DateAfter
						}
					case NumberFilterFieldType:
						switch operand {
						case "eq", "equal":
							filterOperand = NumberEqual
						case "ne", "notequal":
							filterOperand = NumberNotEqual
						case "lt", "lessthan":
							filterOperand = NumberLessThan
						case "le", "lessequal", "lessorequal":
							filterOperand = NumberLessOrEqualThan
						case "gt", "greaterthan":
							filterOperand = NumberGreaterThan
						case "ge", "greaterequal", "greaterorequal":
							filterOperand = NumberGreaterOrEqualThan
						}
					case EnumFilterFieldType:
						switch operand {
						case "eq", "equal":
							filterOperand = EnumEqual
						case "ne", "notequal":
							filterOperand = EnumNotEqual
						}
					}
					if exists {
						queryParams.Filters = append(queryParams.Filters, FilterOption{
							Field:           field,
							Value:           arg,
							FilterOperation: filterOperand,
						})
					}
				}

			}
		}
	})

	return &queryParams
}
