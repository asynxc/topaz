package ds

import (
	"bytes"

	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	"github.com/aserto-dev/topaz/resolvers"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"

	"github.com/rs/zerolog"
)

// RegisterRelation - ds.relation
//
// ds.relation({
// 	"object": {
// 	  "id": "",
// 	  "key": "",
// 	  "type": ""
// 	},
// 	"relation": {
// 	  "name": "",
// 	  "object_type": ""
// 	},
// 	"subject": {
// 	  "id": "",
// 	  "key": "",
// 	  "type": ""
// 	}
// })
//
func RegisterRelation(logger *zerolog.Logger, fnName string, dr resolvers.DirectoryResolver) (*rego.Function, rego.Builtin1) {
	return &rego.Function{
			Name:    fnName,
			Decl:    types.NewFunction(types.Args(types.A), types.A),
			Memoize: false,
		},
		func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
			var a *dsc.RelationIdentifier

			if err := ast.As(op1.Value, &a); err != nil {
				return nil, err
			}

			if a.Object == nil && a.Subject == nil && a.Relation == nil {
				a = &dsc.RelationIdentifier{
					Subject: &dsc.ObjectIdentifier{
						Id:   proto.String(""),
						Type: proto.String(""),
						Key:  proto.String(""),
					},
					Relation: &dsc.RelationTypeIdentifier{
						ObjectType: proto.String(""),
						Name:       proto.String(""),
					},
					Object: &dsc.ObjectIdentifier{
						Id:   proto.String(""),
						Type: proto.String(""),
						Key:  proto.String(""),
					},
				}
				return help(fnName, a)
			}

			client, err := dr.GetDS(bctx.Context)
			if err != nil {
				return nil, errors.Wrapf(err, "get directory client")
			}

			resp, err := client.GetRelation(bctx.Context, &dsr.GetRelationRequest{
				Param: a,
			})
			if err != nil {
				return nil, err
			}

			buf := new(bytes.Buffer)
			if resp != nil {
				if err := ProtoToBuf(buf, resp); err != nil {
					return nil, err
				}
			}

			v, err := ast.ValueFromReader(buf)
			if err != nil {
				return nil, err
			}

			return ast.NewTerm(v), nil
		}
}
