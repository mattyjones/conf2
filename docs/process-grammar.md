    *(function-stmt)

    function-stmt =
        function_name
        *(indent operation-stmt)

    indent = "  " ; 2 spaces

    function_name =
        token_ident EOL

    token_ident = <alphanumeric characters begining with letter>

    operation-stmt =
        *(
            select-stmt /
            if-stmt /
            let-stmt /
            set-stmt
        )

    if-stmt =
        "if" expression-stmt
        *(indent operation-stmt)

    select-stmt =
        "select" token_ident

    let-stmt =
        token_ident "=" expression-stmt

    expression-stmt =
         *(
            function-stmt
            arithmetic-stmt
         )

    arithmetic-stmt =
        *(
            number
            variable
            string
            arithmetic-operator
        )

    string = double-quote * double-quote

    arithmetic-operator =
      *(
        "+" "-" "*" "/"
      )

    function-stmt =
         token_ident "(" [ arguments ] ")"

    arguments =
        expression-stmt ["," expression-stmt]

