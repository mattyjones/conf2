%{
package yang

import (
    "fmt"
    "strings"
    "schema"
)

type yangError struct {
	s string
}

func (err *yangError) Error() string {
	return err.s
}

func tokenString(s string) string {
    return s[1:len(s) - 1]
}

func (l *lexer) Lex(lval *yySymType) int {
    t, _ := l.nextToken()
    if t.typ == ParseEof {
        return 0
    }
    lval.token = t.val
    lval.stack = l.stack
    lval.importer = l.importer
    return int(t.typ)
}

func (l *lexer) Error(e string) {
    line, col := l.Position()
    msg := fmt.Sprintf("%s - line %d, col %d", e, line, col)
    l.lastError = &yangError{msg}
}

func HasError(l yyLexer, e error) bool {
    if e == nil {
        return false
    }
    l.Error(e.Error())
    return true
}

func popAndAddMeta(yylval *yySymType) error {
    i := yylval.stack.Pop()
    if def, ok := i.(schema.Meta); ok {
        parent := yylval.stack.Peek()
        if parentList, ok := parent.(schema.MetaList); ok {
            return parentList.AddMeta(def)
        } else {
            return &yangError{fmt.Sprintf("Cannot add \"%s\" to \"%s\"; not collection type.", i.GetIdent(), parent.GetIdent())}
        }
    } else {
        return &yangError{fmt.Sprintf("\"%s\" cannot be stored in a collection type.", i.GetIdent())}
    }
}

%}

%union {
    ident string
    token string
    dataType *schema.DataType
    stack *yangMetaStack
    importer ImportModule
}

%token <token> token_ident
%token <token> token_string
%token <token> token_int
%token token_curly_open
%token token_curly_close
%token token_semi
%token <token> token_rev_ident

/* KEEP LIST IN SYNC WITH lexer.go */
%token kywd_namespace
%token kywd_description
%token kywd_revision
%token kywd_type
%token kywd_prefix
%token kywd_default
%token kywd_length
%token kywd_enum
%token kywd_key
%token kywd_config
%token kywd_uses
%token kywd_unique
%token kywd_input
%token kywd_output
%token kywd_module
%token kywd_container
%token kywd_list
%token kywd_rpc
%token kywd_notification
%token kywd_typedef
%token kywd_grouping
%token kywd_leaf
%token kywd_mandatory
%token kywd_reference
%token kywd_leaf_list
%token kywd_max_elements
%token kywd_choice
%token kywd_case
%token kywd_import
%token kywd_include
%token kywd_action
%token kywd_anyxml
%token kywd_path

%%

module :
    module_def
    module_stmts
    revision_stmt
    module_body_stmts
    token_curly_close;

module_def :
    kywd_module token_ident token_curly_open {
      m:= &schema.Module{Ident:$2}
      yylval.stack.Push(m)
    }

revision_def :
    kywd_revision token_rev_ident {
        d := yylval.stack.Peek()
        r := &schema.Revision{Ident:$2}
        d.(*schema.Module).Revision = r
        yylval.stack.Push(r)
    }

revision_stmt :
    revision_def token_semi {
      yylval.stack.Pop()
    }
    | revision_def token_curly_open description token_semi token_curly_close {
      yylval.stack.Pop()
    };

description : kywd_description token_string {
        yylval.stack.Peek().(schema.Describable).SetDescription(tokenString($2))
    }

module_stmts :
    module_stmt token_semi
    | module_stmts module_stmt token_semi;

/* TODO: are these optional? */
module_stmt :
    kywd_namespace token_string {
         d := yylval.stack.Peek()
         d.(*schema.Module).Namespace = tokenString($2)
    }
    | description
    | import_stmt
    | kywd_prefix token_string {
         m := yylval.stack.Peek().(*schema.Module)
         m.Prefix = tokenString($2)
    }

import_stmt : kywd_import token_ident {
    var err error
    if yylval.importer == nil {
        yylex.Error("No importer defined")
        goto ret1
    } else {
        m := yylval.stack.Peek().(*schema.Module)
        if err = yylval.importer(m, $2); err != nil {
            yylex.Error(err.Error())
            goto ret1
        }
    }
}

module_body_stmt :
    rpc_stmt
    | notification_stmt
    | body_stmt

module_body_stmts :
    module_body_stmt
    | module_body_stmts module_body_stmt;

optional_body_stmts :
    /*empty*/
    | body_stmts;

body_stmt :
    typedef_stmt
    | grouping_stmt
    | list_stmt
    | container_stmt
    | leaf_stmt
    | leaf_list_stmt
    | anyxml_stmt
    | uses_stmt
    | choice_stmt
    | action_stmt

body_stmts :
    body_stmt | body_stmts body_stmt;

choice_stmt :
    choice_def
    token_curly_open
    choice_stmt_body
    token_curly_close {
      if HasError(yylex, popAndAddMeta(&yylval)) {
        goto ret1
      }
    }

choice_stmt_body :
    /* empty */
    | description token_semi
    | case_stmts

choice_def :
    kywd_choice token_ident {
        yylval.stack.Push(&schema.Choice{Ident:$2})
    };

case_stmts :
    case_stmt | case_stmts case_stmt;

case_stmt :
    case_def token_curly_open
    container_body_stmts
    token_curly_close {
      if HasError(yylex, popAndAddMeta(&yylval)) {
        goto ret1
      }
    }

case_def :
    kywd_case token_ident {
        yylval.stack.Push(&schema.ChoiceCase{Ident:$2})
    };

typedef_stmt :
    typedef_def
    token_curly_open
    typedef_stmt_body
    token_curly_close {
      if HasError(yylex, popAndAddMeta(&yylval)) {
        goto ret1
      }
    };

typedef_def :
    kywd_typedef token_ident {
        yylval.stack.Push(&schema.Typedef{Ident:$2})
    };

typedef_stmt_body :
        typedef_stmt_body_stmt
        | typedef_stmt_body typedef_stmt_body_stmt
        ;

typedef_stmt_body_stmt:
        type_stmt
        | description token_semi
        | default_stmt
        ;
default_stmt : kywd_default token_string token_semi {
     if hasType, valid := yylval.stack.Peek().(schema.HasDataType); valid {
        hasType.GetDataType().SetDefault(tokenString($2))
     } else {
        yylex.Error("expected default statement on schema supporting details")
        goto ret1
     }
}

type_stmt : type_stmt_def type_stmt_body {
         y := yylval.stack.Peek().(schema.HasDataType)
         y.SetDataType(yylval.dataType)
        };

type_stmt_def : kywd_type token_ident {
            y := yylval.stack.Peek().(schema.HasDataType)
            yylval.dataType = schema.NewDataType(y, $2)
        };

type_stmt_body :
        token_semi
        | token_curly_open type_stmt_types token_curly_close;

type_stmt_types :
        kywd_length token_string token_semi {
            var err error
            if err = yylval.dataType.DecodeLength(tokenString($2)); err != nil {
                yylex.Error(err.Error())
                goto ret1
            }
        }
        | enum_stmts
        | kywd_path token_string  token_semi {
              yylval.dataType.SetPath(tokenString($2))
        };

container_stmt :
    container_def
    token_curly_open
    container_body_stmts
    token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    };

container_def :
    kywd_container token_ident {
        yylval.stack.Push(&schema.Container{Ident:$2})
    };

container_body_stmts :
    container_body_stmt
    | container_body_stmts container_body_stmt;

container_body_stmt :
    description token_semi
    | config_stmt
    | mandatory_stmt
    | body_stmt

uses_stmt :
     kywd_uses token_ident token_semi {
         yylval.stack.Push(&schema.Uses{Ident:$2})
         if HasError(yylex, popAndAddMeta(&yylval)) {
           goto ret1
         }
     }
/* TODO: support alternate uses form */

rpc_stmt :
    rpc_def
    token_curly_open
    optional_rpc_body_stmts
    token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    };

rpc_def :
    kywd_rpc token_ident {
        yylval.stack.Push(&schema.Rpc{Ident:$2})
    };

optional_rpc_body_stmts :
    /* empty */
    | rpc_body_stmts;

rpc_body_stmts :
    rpc_body_stmt | rpc_body_stmts rpc_body_stmt;

rpc_body_stmt:
    description token_semi
    | reference_stmt
    | rpc_input optional_body_stmts token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    }
    | rpc_output optional_body_stmts token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    };

rpc_input :
    kywd_input token_curly_open {
        yylval.stack.Push(&schema.RpcInput{})
    };

rpc_output :
    kywd_output token_curly_open {
        yylval.stack.Push(&schema.RpcOutput{})
    };

action_stmt :
    action_def
    token_curly_open
    optional_action_body_stmts
    token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    };

action_def :
    kywd_action token_ident {
        yylval.stack.Push(&schema.Rpc{Ident:$2})
    };

optional_action_body_stmts :
    /* empty */
    | action_body_stmts;

action_body_stmts :
    action_body_stmt | action_body_stmts action_body_stmt;

action_body_stmt:
    description token_semi
    | reference_stmt
    | action_input optional_body_stmts token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    }
    | action_output optional_body_stmts token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    };

action_input :
    kywd_input token_curly_open {
        yylval.stack.Push(&schema.RpcInput{})
    };

action_output :
    kywd_output token_curly_open {
        yylval.stack.Push(&schema.RpcOutput{})
    };

notification_stmt :
    notification_def
    token_curly_open
    notification_body_stmts
    token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    };

notification_def :
    kywd_notification token_ident {
        yylval.stack.Push(&schema.Notification{Ident:$2})
    };

notification_body_stmts :
    notification_body_stmt
    | notification_body_stmts notification_body_stmt;

/* TODO: if, stats, reference, typedef*/
notification_body_stmt :
    description token_semi
    | config_stmt
    | mandatory_stmt
    | body_stmt;

grouping_stmt :
    grouping_def
    grouping_body_defined {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    };

grouping_body_defined:
    token_curly_open
    grouping_body_stmts
    token_curly_close;

grouping_def :
    kywd_grouping token_ident {
        yylval.stack.Push(&schema.Grouping{Ident:$2})
    };

grouping_body_stmts :
    grouping_body_stmt |
    grouping_body_stmts grouping_body_stmt;

grouping_body_stmt :
    description token_semi
    | reference_stmt
    | body_stmt;

list_stmt :
    list_def token_curly_open
    list_body_stmts
    token_curly_close{
         if HasError(yylex, popAndAddMeta(&yylval)) {
             goto ret1
         }
     };

list_def :
    kywd_list token_ident {
        yylval.stack.Push(&schema.List{Ident:$2})
    };

list_body_stmts :
    list_body_stmt
    | list_body_stmts list_body_stmt;

list_body_stmt :
    description token_semi
    | kywd_max_elements token_int token_semi
    | config_stmt
    | mandatory_stmt
    | key_stmt
    | kywd_unique token_string token_semi
    | body_stmt

key_stmt: kywd_key token_string token_semi {
     if list, valid := yylval.stack.Peek().(*schema.List); valid {
       list.Key = strings.Split(tokenString($2), " ")
     } else {
        yylex.Error("expected a list for key statement")
        goto ret1
     }
}

anyxml_stmt:
    anyxml_def token_semi {
         if HasError(yylex, popAndAddMeta(&yylval)) {
             goto ret1
         }
    }

anyxml_def :
    kywd_anyxml token_ident {
        yylval.stack.Push(schema.NewAny($2))
    }

leaf_stmt:
    leaf_def
    token_curly_open
    leaf_body_stmts
    token_curly_close {
         if HasError(yylex, popAndAddMeta(&yylval)) {
             goto ret1
         }
     };

leaf_def :
    kywd_leaf token_ident {
        yylval.stack.Push(&schema.Leaf{Ident:$2})
    };

leaf_body_stmts :
    leaf_body_stmt
    | leaf_body_stmts leaf_body_stmt;

/* TODO: when, if, units, must, status, reference */
leaf_body_stmt :
    type_stmt
    | description token_semi
    | config_stmt
    | mandatory_stmt
    | default_stmt

mandatory_stmt : kywd_mandatory token_string token_semi {
      if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
         hasDetails.Details().SetMandatory("true" == $2)
      } else {
         yylex.Error("expected mandatory statement on schema supporting details")
         goto ret1
      }
 };

config_stmt : kywd_config token_string token_semi {
     if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
        hasDetails.Details().SetConfig("true" == $2)
     } else {
        yylex.Error("expected config statement on schema supporting details")
        goto ret1
     }
};

/* TODO: when, if, units, must, status, reference, min, max */
leaf_list_stmt :
    leaf_list_def
    token_curly_open
    leaf_body_stmts
    token_curly_close {
        if HasError(yylex, popAndAddMeta(&yylval)) {
            goto ret1
        }
    };

leaf_list_def :
    kywd_leaf_list token_ident {
        yylval.stack.Push(&schema.LeafList{Ident:$2})
    };

enum_stmts :
    enum_stmt
    | enum_stmts enum_stmt;

enum_stmt :
    kywd_enum token_ident token_semi {
        yylval.dataType.AddEnumeration($2)
    };

reference_stmt :
    kywd_reference token_string token_semi;

%%

