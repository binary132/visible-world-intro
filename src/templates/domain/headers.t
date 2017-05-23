package {{pkgName .DomainPkg}}

import (
  "inventory_service/spec"
  "inventory_service/spec/validator"
  "inventory_service/models"
	"inventory_service/errors"
	ctx "inventory_service/context"
	"inventory_service/parser"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/vattle/sqlboiler/boil"
  "encoding/json"
	"github.com/vattle/sqlboiler/queries/qm"
  "fmt"
)

{{template "domain.base" .}}
{{template "domain.show" .}}
{{template "domain.create" .}}
{{template "domain.update" .}}
{{template "domain.upsert" .}}
{{template "domain.search" .}}
{{template "domain.relation" .}}
