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
    sel *process.Select
    args []process.Expression
    expr process.Expression
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
%token <token> token_function

%token kywd_select
%token kywd_into
%token kywd_let
%token kywd_goto
%token kywd_if

%type <sel> select_def
%type <args> arguments
%type <args> optional_arguments
%type <expr> expression
%type <expr> function

%%

scripts :
    script | scripts script

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
    | select
    | if
    | goto

eol : token_eol {
    yylval.stack.nextDepth = 0;
}

set :
    token_ident token_equal expression eol {
       op := &process.Set{Name:$1}
       yylval.stack.Push(op)
    }

let :
    kywd_let token_ident token_equal expression eol {
        op := &process.Let{Name:$2}
        yylval.stack.Push(op)
    }

if :
    kywd_if expression eol {
        op := &process.If{}
        yylval.stack.Push(op)
    }

goto :
    kywd_goto token_ident eol {
        op := &process.Goto{Script:$2}
        yylval.stack.Push(op)
    };

select_def :
    kywd_select token_ident {
        op := &process.Select{On:$2}
        yylval.stack.Push(op)
        $$ = op
    }

select :
    select_def eol
    | select_def kywd_into token_ident eol {
        $1.Into = $3
    }

expression :
    token_ident {
        $$ = &process.Primative{Var:$1}
    }
    | token_string {
        $$ = &process.Primative{Str:$1}
    }
    | token_number {
        $$ = &process.Primative{Num:0}
    }
    | function;

function :
    token_function token_open_paren optional_arguments token_close_paren {
        $$ = &process.Function{Name:$1, Arguments:$3}
    }

optional_arguments :
    /* empty */ {
        $$ = []process.Expression{}
    }
    | arguments;

arguments :
    expression {
        $$ = []process.Expression{$1}
    }
    | arguments expression {
        $$ = append($1, $2)
    }

indent : token_space_indent {
    yylval.stack.depth = yylval.stack.nextDepth
    yylval.stack.nextDepth++
}

indentation :
    indent | indentation indent;

script_def : token_ident eol {
   s := &process.Script{Name:$1}
   yylval.stack.PushScript(s)
}
%%
