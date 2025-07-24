package database

// Re-export types and functions from internal/db package
import (
	internaldb "sample-micro-service-api/package-go/database/internal/db"
)

// Re-export main entity types
type (
	GcasUser              = internaldb.GcasUser
	GcasGroup             = internaldb.GcasGroup
	GcasGroupUserRelation = internaldb.GcasGroupUserRelation
	GcasGroupSystemRelation = internaldb.GcasGroupSystemRelation
	Project               = internaldb.Project
	ProjectCost           = internaldb.ProjectCost
	ProjectSystemRelation = internaldb.ProjectSystemRelation
	System                = internaldb.System
	SystemBasicInformation = internaldb.SystemBasicInformation
	Queries               = internaldb.Queries
	DBTX                  = internaldb.DBTX
)

// Re-export master table types
type (
	MLocalGovernment      = internaldb.MLocalGovernment
	MOrganizationCategory = internaldb.MOrganizationCategory
	MUserRole             = internaldb.MUserRole
)


// Re-export parameter types for System
type (
	CreateSystemParams        = internaldb.CreateSystemParams
	UpdateSystemParams        = internaldb.UpdateSystemParams
	UpdateSystemContactParams = internaldb.UpdateSystemContactParams
)

// Re-export constructor
func New(db DBTX) *Queries {
	return internaldb.New(db)
} 