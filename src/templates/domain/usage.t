{{define "domain.usage.create.func" -}}
/*
Create {{ titleCase .Name }}
* Validate request body through generated validators
* Unmarshal request body to Spec Model
* Enter a transaction by `WithExecutor` and load Domain Model from Spec Model
* Perform biz related validations
* Reload the created entity by `WithRequest`

Make sure all TODOs are resolved properly.
*/
{{- end}}

{{define "domain.usage.create.func.1" -}}
// Validate the request body with the generated validator
{{- end}}
{{define "domain.usage.create.func.2" -}}
// Load Domain Model from Spec Model
{{- end}}
{{define "domain.usage.create.func.3" -}}
// Perform creating related operations within transaction by `WithExecutor`
{{- end}}
{{define "domain.usage.create.func.4" -}}
// TODO: Load Domain Model from the Spec Model,
// unmarshalled from the raw data of request body
{{- end}}
{{define "domain.usage.create.func.5" -}}
// TODO: Perform the creating operations through Domain Model.
// Here the generated skeleton is to perform insertion on the primary Database
// Model, such as creating associations and relations.
// Go to [SQLBoiler](https://github.com/vattle/sqlboiler) for more details.
// You can perform other database operations through other database models
// Checkout available Database Models under `inventory_service/models`.
{{- end}}
{{define "domain.usage.create.func.6" -}}
// TODO: Reload the Domain Model, as the response data,
// outside the transaction with Executor passed by `WithRequest`.
{{- end}}



{{define "domain.usage.update.func" -}}
/*
Update {{ titleCase .Name }}
* Validate request body through generated validators
* Unmarshal request body to Spec Model
* Enter a transaction by `WithExecutor` and load Domain Model from Spec Model
* Perform biz related validations
* Reload the created entity by `WithRequest`

Make sure all TODOs are resolved properly.
*/
{{- end}}

{{define "domain.usage.update.func.1" -}}
// Validate the request body with the generated validator
{{- end}}
{{define "domain.usage.update.func.2" -}}
// Get entity ID from params
{{- end}}
{{define "domain.usage.update.func.3" -}}
// Load exist entity from database
{{- end}}
{{define "domain.usage.update.func.4" -}}
// TODO: Load Domain Model from the Spec Model,
// unmarshalled from the raw data of request body
{{- end}}
{{define "domain.usage.update.func.5" -}}
// TODO: Perform the updating operations through Domain Model.
// Here the generated skeleton is to perform insertion on the primary Database
// Model, such as creating associations and relations.
// Go to [SQLBoiler](https://github.com/vattle/sqlboiler) for more details.
// You can perform other database operations through other database models
// Checkout available Database Models under `inventory_service/models`.
{{- end}}
{{define "domain.usage.update.func.6" -}}
// TODO: Reload the Domain Model, as the response data,
// outside the transaction with Executor passed by `WithRequest`.
{{- end}}



{{define "domain.usage.relations.func.1" -}}
// TODO: Handle creating/editing relation
// Get Relation Name through the parsed URL
// ```
// // Path: /sites/1/parent_site_groups
// // return: parent_site_groups
// relName := ctx.URL().Resource().Sub().Name()
// ```
// Create Relation:
// ```
// d.setRelation(exec, relName, s)
// ```
// Update/Append Relation:
// ```
// d.setRelation(exec, relName, s)
// ```
{{- end}}
{{define "domain.usage.relations.func.2" -}}
// TODO: List relations, such as Site.ParentSiteGroups or Site.ChildSiteSections
// Get Relation Name through the parsed URL
// ```
// // Path: /sites/1/parent_site_groups
// // return: parent_site_groups
// relName := ctx.URL().Resource().Sub().Name()
// ```
// Get Relation:
// ```
// d.getRelation(exec, relName, ctx.Params())
// ```
{{- end}}
{{define "domain.usage.one_relation.func.1" -}}
// TODO: Add one relation
// Get Relation Name through the parsed URL
// ```
// // Path: /sites/1/parent_site_groups
// // return: parent_site_groups
// relName := ctx.URL().Resource().Sub().Name()
// ```
// Replace one relation by id:
// ```
// d.setRelationByIDs(exec, relName, relID)
// ```
// Append one relation by id:
// ```
// d.addRelationByIDs(exec, relName, relID)
// ```
{{- end}}
{{define "domain.usage.one_relation.func.2" -}}
// TODO: Delete one relation
// Get Relation Name through the parsed URL
// ```
// // Path: /sites/1/parent_site_groups
// // return: parent_site_groups
// relName := ctx.URL().Resource().Sub().Name()
// ```
//
// Delete one relation by id:
// ```
// d.deleteRelationByIDs(exec, relName, relID)
// ```
{{- end}}



{{define "domain.usage.load_spec.func" -}}
/*
TODO: Load Spec Model to Domain Model, for instance,

Site.loadUpdateSpec

```
  // adaptJsonValue converts nullable value to field of Domain Model
	adaptJsonValue(s.Description, &d.Description)
	adaptJsonValue(s.Metadata, &d.Metadata)
	adaptJsonValue(s.URL, &d.URL)
	adaptJsonValue(s.ExternalID, &d.InternalID)
	adaptJsonValue(s.Name, &d.Name)
	adaptJsonValue(s.Tag, &d.GroupTag)
	adaptJsonValue(s.SessionDuration, &d.SessionDuration)

	adaptJsonValue(s.Status, &d.Status)
	if strings.ToUpper(d.Status) == "INACTIVE" {
		d.Status = models.SiteSectionGroupStatusIN_ACTIVE
	}

	var ratingStr string
	adaptJsonValue(s.Rating, &ratingStr)
	if ratingStr != "" {
		d.Rating = null.Int64From(constants.Ratings.Lookup(networkID, ratingStr))
	}
```

Site.loadCreateSpec

```
	d.GroupType = models.SiteSectionGroupGroupTypeSITE
	adaptJsonValue(s.Description, &d.Description)
	adaptJsonValue(s.Metadata, &d.Metadata)
	adaptJsonValue(s.URL, &d.URL)
	adaptJsonValue(s.ExternalID, &d.InternalID)
	adaptJsonValue(s.SessionDuration, &d.SessionDuration)
	d.NetworkID = networkID

	if s.Name != nil {
		d.Name = *s.Name
	}

	if s.Tag != nil {
		d.GroupTag = *s.Tag
	}

	adaptJsonValue(s.Status, &d.Status)
	if strings.ToUpper(d.Status) == "INACTIVE" {
		d.Status = models.SiteSectionGroupStatusIN_ACTIVE
	}

	var ratingStr string
	adaptJsonValue(s.Rating, &ratingStr)
	if len(ratingStr) > 0 {
		d.Rating = null.Int64From(constants.Ratings.Lookup(networkID, ratingStr))
	}
```
*/
{{- end}}



{{define "domain.usage.set_relation.func" -}}
// TODO: Set named relation by IDs, exist ones should be removed
{{- end}}

{{define "domain.usage.add_relation.func" -}}
// TODO: Add named relation by IDs, appending to existing ones
{{- end}}

{{define "domain.usage.delete_relation.func" -}}
// TODO: Delete named relation by IDs
{{- end}}

{{define "domain.usage.marshaljson.func" -}}
// TODO: marshal domain model to json
// For instance,
//
// ```
// o := &spec.SiteShow{
// 	ID:          &d.ID,
// 	Name:        &d.Name,
// 	Tag:         &d.GroupTag,
// 	NetworkID:   &d.NetworkID,
// 	Metadata:    d.Metadata.Ptr(),
// 	Description: d.Description.Ptr(),
// 	ExternalID:  d.InternalID.Ptr(),
// 	URL:         d.URL.Ptr(),
// 	CreatedAt:   null.TimeFrom(d.CreatedAt),
// 	UpdatedAt:   null.TimeFrom(d.UpdatedAt),
// 	Links: []*spec.LinkItem{
// 		{Ref: "self", Href: fmt.Sprintf("/sites/%d", d.ID)},
// 		{Ref: "parent_site_groups", Href: fmt.Sprintf("/sites/%d/parent_site_groups", d.ID)},
// 		{Ref: "child_site_sections", Href: fmt.Sprintf("/sites/%d/child_site_sections", d.ID)},
// 	},
// }
// ```
{{- end}}

{{define "domain.usage.search.func" -}}
// TODO: Search {{.}} from Search Server and Database
{{- end}}
