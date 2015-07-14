%{
package yang

import (
    "fmt"
)

func (l *lexer) Lex(lval *yySymType) int {
    t, _ := l.nextToken()
    if t.typ == ParseEof {
        return 0
    }
    lval.token = t.val
    lval.stack = l.stack
    return int(t.typ)
}

func (l *lexer) Error(e string) {
    line, col := l.Position()
    fmt.Println(fmt.Sprintf("%s at line %d, col %d", e, line, col))
}

func popAndAddDef(yylval *yySymType) bool {
    i := yylval.stack.Pop()
    def, ok := i.(Def)
    if ok {
        parent := yylval.stack.Peek()
        parentList, ok := parent.(DefList)
        if ok {
            err := parentList.AddDef(def)
            if err == nil {
                return true
            }
            if yyDebug > 1 {
                __yyfmt__.Printf(err.Error())
            }
        } else if yyDebug > 1 {
            __yyfmt__.Printf("Internal Error: %s doesn't implement DefList.", parent.GetIdent())
        }
    } else {
        __yyfmt__.Printf("Internal Error: %s doesn't implement Def.", i.GetIdent())
    }

    return false
}

%}

%union {
    ident string
    token string
    stack *yangDefStack
}

%token <token> token_ident
%token <token> token_string
%token <token> token_int
%token token_curly_open
%token token_curly_close
%token token_semi
%token <token> token_rev_ident

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

%%

module :
    module_def
    module_stmts
    revision_stmt
    module_body_stmts
    token_curly_close;

module_def :
    kywd_module token_ident token_curly_open {
      m:= &Module{Ident:$2}
      yylval.stack.Push(m)
    }

revision_def :
    kywd_revision token_rev_ident {
        d := yylval.stack.Peek()
        r := &Revision{Ident:$2}
        d.(*Module).Revision = r
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
        yylval.stack.Peek().(Describable).SetDescription($2)
    }

module_stmts :
    module_stmt token_semi
    | module_stmts module_stmt token_semi;

/* TODO: are these optional? */
module_stmt :
    kywd_namespace token_string {
         d := yylval.stack.Peek()
         d.(*Module).Namespace = $2
    }
    | description
    | kywd_prefix token_string {
         m := yylval.stack.Peek().(*Module)
         m.Prefix = $2
    }

module_body_stmt :
    typedef_stmt
    | grouping_stmt
    | list_stmt
    | container_stmt
    | rpc_stmt
    | notification_stmt

module_body_stmts :
    module_body_stmt {
      if ! popAndAddDef(&yylval) {
        goto ret1
      }
    }
    | module_body_stmts module_body_stmt {
      if ! popAndAddDef(&yylval) {
        goto ret1
      }
    };

optional_body_stmts :
    /*empty*/
    | body_stmts;

body_stmt :
    list_stmt
    | container_stmt
    | leaf_stmt
    | leaf_list_stmt

body_stmts :
    body_stmt  {
        if ! popAndAddDef(&yylval) {
          goto ret1
        }
      }
    | body_stmts body_stmt  {
       if ! popAndAddDef(&yylval) {
         goto ret1
       }
     };

typedef_stmt :
        typedef_def
        token_curly_open
        typedef_stmt_body
        token_curly_close;

typedef_def :
    kywd_typedef token_ident {
        yylval.stack.Push(&Typedef{Ident:$2})
    };

typedef_stmt_body :
        typedef_stmt_body_stmt
        | typedef_stmt_body typedef_stmt_body_stmt
        ;

typedef_stmt_body_stmt:
        type_stmt
        | description token_semi
        | kywd_default token_string token_semi;

type_stmt :
        kywd_type token_ident type_stmt_body
        ;

type_stmt_body :
        token_curly_open type_stmt_types token_curly_close
        | token_semi
        ;

type_stmt_types :
        kywd_length token_string token_semi
        | enum_stmts;

container_stmt :
    container_def
    token_curly_open
    container_body_stmts
    token_curly_close;

container_def :
    kywd_container token_ident {
        yylval.stack.Push(&Container{Ident:$2})
    };

container_body_stmts :
    container_body_stmt
    | container_body_stmts container_body_stmt;

container_body_stmt :
    description token_semi
    | kywd_uses token_ident token_semi
    | kywd_config token_string token_semi
    | body_stmt {
         if ! popAndAddDef(&yylval) {
           goto ret1
         }
       }

rpc_stmt :
    rpc_def
    token_curly_open
    rpc_body_stmts
    token_curly_close;

rpc_def :
    kywd_rpc token_ident {
        yylval.stack.Push(&Rpc{Ident:$2})
    };

rpc_body_stmts :
    rpc_body_stmt | rpc_body_stmts rpc_body_stmt;

/* TODO: add if, status, typedef, grouping, output  */
rpc_body_stmt:
    description token_semi
    | reference_stmt
    | rpc_input optional_body_stmts token_curly_close {
         input := yylval.stack.Pop().(*RpcInput)
         rpc := yylval.stack.Peek().(*Rpc)
         rpc.Input = input
    }
    | rpc_output optional_body_stmts token_curly_close {
         output := yylval.stack.Pop().(*RpcOutput)
         rpc := yylval.stack.Peek().(*Rpc)
         rpc.Output = output
    };

rpc_input :
    kywd_input token_curly_open {
        yylval.stack.Push(&RpcInput{})
    };

rpc_output :
    kywd_output token_curly_open {
        yylval.stack.Push(&RpcOutput{})
    };

notification_stmt :
    notification_def
    token_curly_open
    notification_body_stmts
    token_curly_close;

notification_def :
    kywd_notification token_ident {
        yylval.stack.Push(&Notification{Ident:$2})
    };

notification_body_stmts :
    notification_body_stmt
    | notification_body_stmts notification_body_stmt;

/* TODO: if, stats, reference, typedef*/
notification_body_stmt :
    description token_semi
    | kywd_uses token_ident token_semi
    | kywd_config token_string token_semi
    | body_stmt {
        if ! popAndAddDef(&yylval) {
            goto ret1
        }
    };

grouping_stmt :
    grouping_def
    grouping_body_defined;

grouping_body_defined:
    token_curly_open
    grouping_body_stmts
    token_curly_close;

grouping_def :
    kywd_grouping token_ident {
        yylval.stack.Push(&Grouping{Ident:$2})
    };

grouping_body_stmts :
    grouping_body_stmt |
    grouping_body_stmts grouping_body_stmt;

grouping_body_stmt :
    description token_semi
    | reference_stmt
    | body_stmt {
       if ! popAndAddDef(&yylval) {
           goto ret1
       }
    };

list_stmt :
    list_def token_curly_open
    list_body_stmts
    token_curly_close;

list_def :
    kywd_list token_ident {
        yylval.stack.Push(&List{Ident:$2})
    };

list_body_stmts :
    list_body_stmt
    | list_body_stmts list_body_stmt;

list_body_stmt :
    description token_semi
    | kywd_max_elements token_int token_semi
    | kywd_uses token_string token_semi
    | kywd_config token_string token_semi
    | kywd_key token_string token_semi
    | kywd_unique token_string token_semi
    | body_stmt  {
        if ! popAndAddDef(&yylval) {
            goto ret1
        }
    }

leaf_stmt:
    leaf_def
    token_curly_open
    leaf_body_stmts
    token_curly_close;

leaf_def :
    kywd_leaf token_ident {
        yylval.stack.Push(&Leaf{Ident:$2})
    };

leaf_body_stmts :
    leaf_body_stmt
    | leaf_body_stmts leaf_body_stmt;

/* TODO: when, if, units, must, status, reference */
leaf_body_stmt :
    type_stmt
    | description token_semi
    | kywd_uses token_ident token_semi
    | kywd_config token_string token_semi
    | kywd_default token_string token_semi
    | kywd_mandatory token_string token_semi

/* TODO: when, if, units, must, status, reference, min, max */
leaf_list_stmt :
    leaf_list_def
    token_curly_open
    leaf_body_stmts
    token_curly_close;

leaf_list_def :
    kywd_leaf_list token_ident {
        yylval.stack.Push(&LeafList{Ident:$2})
    };

enum_stmts :
    enum_stmt | enum_stmts enum_stmt;

enum_stmt :
    kywd_enum token_ident token_semi;

reference_stmt :
    kywd_reference token_string token_semi;

%%

func parse(yang string) int {
    l := lex(yang)
    err := yyParse(l)
    return err
}
