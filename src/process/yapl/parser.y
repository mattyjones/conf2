%{
package yapl

import (
    "process"
)

func (l *lexer) Lex(lval *yySymType) int {
    t, _ := l.nextToken()
    if t.typ == ParseEof {
        return 0
    }
    lval.token = t.val
    lval.stack = &l.stack
    return int(t.typ)
}

%}

%union {
    ident string
    token string
    stack *operatorStack
}

%token <token> token_ident
%token <token> token_string
%token <token> token_number
%token token_space_indent
%token token_eol
%token token_open_paren
%token token_close_paren
%token token_comma
%token token_equal

%token kywd_select
%token kywd_into
%token kywd_let
%token kywd_goto
%token kywd_if

%%

script :
    script_def
    operations;

operations :
     operation | operations operation;

operation :
    indentation operation_stmt;

operation_stmt :
    set
    | let

set :
    token_ident token_equal expression token_eol;

let :
    kywd_let token_ident token_equal expression token_eol;

expression :
    token_ident;

indentation :
    token_space_indent | indentation token_space_indent;

script_def : token_ident token_eol {
   s := &process.Script{Name:$1}
   yylval.stack.PushScript(s)
}

%%