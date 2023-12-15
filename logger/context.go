package logger

import (
	"context"
)

type ctxKeyType string

var LoggerCtxFieldsKey ctxKeyType = "logger.fields"

// AppendFieldsToContext appends the given fields to the context
func AppendFieldsToContext(ctx context.Context, fields ...Field) context.Context {

	var existingLogFields []Field
	if f, ok := ctx.Value(LoggerCtxFieldsKey).([]Field); ok {
		existingLogFields = f
	}

	// Merge the existing fields with the new fields
	existingFieldsMap := make(map[string]Field)
	for _, existingField := range existingLogFields {
		existingFieldsMap[existingField.Key] = existingField
	}
	for _, incomingField := range fields {
		existingFieldsMap[incomingField.Key] = incomingField
	}

	var newLogFields []Field

	for _, newField := range existingFieldsMap {
		newLogFields = append(newLogFields, newField)
	}
	return context.WithValue(ctx, LoggerCtxFieldsKey, newLogFields)
}

// GetFieldValueFromContext returns the string value from the context
func GetFieldValueFromContext(ctx context.Context, FieldName string) string {
	if f, ok := ctx.Value(LoggerCtxFieldsKey).([]Field); ok {
		for _, field := range f {
			if field.Key == FieldName {
				return field.String
			}
		}
	}
	return ""
}
