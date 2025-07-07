package mapper

import (
	"reflect"
	"strings"
)

// FieldMapping represents a field mapping configuration
type FieldMapping struct {
	Source      string
	Target      string
	Transform   func(interface{}) interface{}
	Ignore      bool
	Default     interface{}
	Conditional func(interface{}) bool
}

// MapperConfig holds the configuration for a mapper
type MapperConfig struct {
	FieldMappings map[string]FieldMapping
	IgnoreFields  []string
	PostProcessor func(interface{}, interface{}) interface{}
}

// Engine handles the actual mapping logic
type Engine struct{}

// NewEngine creates a new mapping engine
func NewEngine() *Engine {
	return &Engine{}
}

// Map performs the mapping between source and target using the provided config
func (e *Engine) Map(source interface{}, target interface{}, config MapperConfig) interface{} {
	sourceVal := reflect.ValueOf(source)
	targetVal := reflect.ValueOf(target)

	// Handle pointers
	if sourceVal.Kind() == reflect.Ptr {
		sourceVal = sourceVal.Elem()
	}
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}

	// sourceType := sourceVal.Type()
	targetType := targetVal.Type()

	// Create new instance of target type
	newTarget := reflect.New(targetType).Elem()

	// Map fields
	for i := 0; i < targetType.NumField(); i++ {
		targetField := targetType.Field(i)
		targetFieldVal := newTarget.Field(i)

		// Skip unexported fields
		if !targetFieldVal.CanSet() {
			continue
		}

		// Check if field should be ignored
		if e.shouldIgnoreField(targetField.Name, config.IgnoreFields) {
			continue
		}

		// Check for custom field mapping
		if mapping, exists := config.FieldMappings[targetField.Name]; exists {
			if mapping.Ignore {
				continue
			}
			
			// Apply custom mapping
			e.applyFieldMapping(sourceVal, targetFieldVal, mapping)
			continue
		}

		// Default field mapping (same name)
		if sourceField := sourceVal.FieldByName(targetField.Name); sourceField.IsValid() {
			e.mapField(sourceField, targetFieldVal)
		}
	}

	result := newTarget.Interface()

	// Apply post-processing if configured
	if config.PostProcessor != nil {
		result = config.PostProcessor(result, source)
	}

	return result
}

// shouldIgnoreField checks if a field should be ignored
func (e *Engine) shouldIgnoreField(fieldName string, ignoreFields []string) bool {
	for _, ignoredField := range ignoreFields {
		if ignoredField == fieldName {
			return true
		}
	}
	return false
}

// applyFieldMapping applies a custom field mapping
func (e *Engine) applyFieldMapping(sourceVal reflect.Value, targetFieldVal reflect.Value, mapping FieldMapping) {
	var sourceFieldVal reflect.Value
	
	// Handle nested field access (e.g., "User.Name")
	if strings.Contains(mapping.Source, ".") {
		sourceFieldVal = e.getNestedField(sourceVal, mapping.Source)
	} else {
		sourceFieldVal = sourceVal.FieldByName(mapping.Source)
	}

	// Apply conditional logic
	if mapping.Conditional != nil && !mapping.Conditional(sourceFieldVal.Interface()) {
		return
	}

	// Use default value if source field is not valid
	if !sourceFieldVal.IsValid() {
		if mapping.Default != nil {
			defaultVal := reflect.ValueOf(mapping.Default)
			if defaultVal.Type().AssignableTo(targetFieldVal.Type()) {
				targetFieldVal.Set(defaultVal)
			}
		}
		return
	}

	// Apply transformation if provided
	if mapping.Transform != nil {
		transformed := mapping.Transform(sourceFieldVal.Interface())
		transformedVal := reflect.ValueOf(transformed)
		if transformedVal.Type().AssignableTo(targetFieldVal.Type()) {
			targetFieldVal.Set(transformedVal)
		}
		return
	}

	// Direct mapping
	e.mapField(sourceFieldVal, targetFieldVal)
}

// getNestedField retrieves a nested field value
func (e *Engine) getNestedField(val reflect.Value, fieldPath string) reflect.Value {
	parts := strings.Split(fieldPath, ".")
	current := val
	
	for _, part := range parts {
		if !current.IsValid() {
			return reflect.Value{}
		}
		current = current.FieldByName(part)
	}
	
	return current
}

// mapField maps a single field
func (e *Engine) mapField(sourceField, targetField reflect.Value) {
	if !sourceField.IsValid() || !targetField.CanSet() {
		return
	}

	// Handle different type conversions
	if sourceField.Type().AssignableTo(targetField.Type()) {
		targetField.Set(sourceField)
	} else if sourceField.Type().ConvertibleTo(targetField.Type()) {
		targetField.Set(sourceField.Convert(targetField.Type()))
	}
}

// Builder provides a fluent API for building mapper configurations
type Builder struct {
	config MapperConfig
}

// NewBuilder creates a new mapper configuration builder
func NewBuilder() *Builder {
	return &Builder{
		config: MapperConfig{
			FieldMappings: make(map[string]FieldMapping),
			IgnoreFields:  []string{},
		},
	}
}

// MapField adds a field mapping
func (b *Builder) MapField(target, source string) *Builder {
	b.config.FieldMappings[target] = FieldMapping{
		Source: source,
		Target: target,
	}
	return b
}

// MapFieldWithTransform adds a field mapping with transformation
func (b *Builder) MapFieldWithTransform(target, source string, transform func(interface{}) interface{}) *Builder {
	b.config.FieldMappings[target] = FieldMapping{
		Source:    source,
		Target:    target,
		Transform: transform,
	}
	return b
}

// MapFieldWithDefault adds a field mapping with default value
func (b *Builder) MapFieldWithDefault(target, source string, defaultValue interface{}) *Builder {
	b.config.FieldMappings[target] = FieldMapping{
		Source:  source,
		Target:  target,
		Default: defaultValue,
	}
	return b
}

// MapFieldWithCondition adds a field mapping with condition
func (b *Builder) MapFieldWithCondition(target, source string, condition func(interface{}) bool) *Builder {
	b.config.FieldMappings[target] = FieldMapping{
		Source:      source,
		Target:      target,
		Conditional: condition,
	}
	return b
}

// IgnoreField adds a field to ignore
func (b *Builder) IgnoreField(fieldName string) *Builder {
	b.config.IgnoreFields = append(b.config.IgnoreFields, fieldName)
	return b
}

// IgnoreFields adds multiple fields to ignore
func (b *Builder) IgnoreFields(fieldNames ...string) *Builder {
	b.config.IgnoreFields = append(b.config.IgnoreFields, fieldNames...)
	return b
}

// PostProcess adds a post-processing function
func (b *Builder) PostProcess(processor func(interface{}, interface{}) interface{}) *Builder {
	b.config.PostProcessor = processor
	return b
}

// Build creates the mapper configuration
func (b *Builder) Build() MapperConfig {
	return b.config
}