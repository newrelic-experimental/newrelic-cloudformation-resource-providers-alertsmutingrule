package resource

import (
   "fmt"
   "github.com/newrelic/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

//
// Generic, should be able to leave these as-is
//

type Payload struct {
   model  *Model
   models []interface{}
}

func (p *Payload) SetIdentifier(g *string) {
   p.model.Id = g
}

func (p *Payload) GetIdentifier() *string {
   return p.model.Id
}

func (p *Payload) GetIdentifierKey(a model.Action) string {
   return "id"
}

func (p *Payload) HasTags() bool {
   return false
}

func (p *Payload) GetTags() map[string]string {
   return nil
}

var blank = ""

func (p *Payload) GetTagIdentifier() *string {
   return &blank
}

func NewPayload(m *Model) *Payload {
   return &Payload{
      model:  m,
      models: make([]interface{}, 0),
   }
}

func (p *Payload) GetResourceModel() interface{} {
   return p.model
}

func (p *Payload) GetResourceModels() []interface{} {
   log.Debugf("GetResourceModels: returning %+v", p.models)
   return p.models
}

func (p *Payload) AppendToResourceModels(m model.Model) {
   p.models = append(p.models, m.GetResourceModel())
}

//
// These are API specific, must be configured per API
//

var typeName = "NewRelic::Observability::AlertsNrqlCondition"

func (p *Payload) NewModelFromGuid(g interface{}) (m model.Model) {
   s := fmt.Sprintf("%s", g)
   return NewPayload(&Model{Id: &s})
}

var emptyString = ""

func (p *Payload) GetGraphQLFragment() *string {
   return &emptyString
}

func (p *Payload) GetListQueryNextCursor() string {
   return p.GetListQuery()
}

func (p *Payload) GetVariables() map[string]string {
   vars := make(map[string]string)
   if p.model.Variables != nil {
      for k, v := range p.model.Variables {
         vars[k] = v
      }
   }

   if p.model.Id != nil {
      vars["ID"] = *p.model.Id
   }

   if p.model.Rule != nil {
      vars["RULE"] = *p.model.Rule
   }

   if p.model.ListQueryFilter != nil {
      vars["LISTQUERYFILTER"] = *p.model.ListQueryFilter
   }

   return vars
}

func (p *Payload) GetErrorKey() string {
   return ""
}

func (p *Payload) GetCreateMutation() string {
   return `
mutation {
      alertsMutingRuleCreate(accountId: {{{ACCOUNTID}}}, {{{RULE}}}
      ) {
        id
      }
    }`
}

func (p *Payload) GetDeleteMutation() string {
   return `
mutation {
  alertsMutingRuleDelete(accountId: {{{ACCOUNTID}}}, id: "{{{ID}}}") {
    id
  }
}
`
}

func (p *Payload) GetUpdateMutation() string {
   return `
mutation {
      alertsMutingRuleUpdate(accountId: {{{ACCOUNTID}}}, id: "{{{ID}}}" {{{RULE}}}
      ) {
        id
      }
    }
`
}

func (p *Payload) GetReadQuery() string {
   return `
{
  actor {
    account(id: {{{ACCOUNTID}}}) {
      alerts {
        mutingRule(id: "{{{ID}}}") {
          id
        }
      }
    }
  }
}
`
}

func (p *Payload) GetListQuery() string {
   return `
{
  actor {
    account(id: {{{ACCOUNTID}}}) {
      alerts {
        mutingRules {
          accountId
          id
        }
      }
    }
  }
}
`
}
