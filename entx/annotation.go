package entx

var AnnotationName = "I12R_ENTX"

// Annotation provides a ent.Annotaion spec
type Annotation struct {
	IsNamespacedDataJSONField bool
	ExternalEdge              *ExternalEdgeConfig
}

// Name implements the ent Annotation interface.
func (a Annotation) Name() string {
	return AnnotationName
}
